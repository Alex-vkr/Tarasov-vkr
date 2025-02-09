package parser

import (
	"diplom-vuls-server/internal/model"
	custom_errors "diplom-vuls-server/pkg/custom-errors"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func ParseSoftwareProgram() ([]*model.SoftwareProgramItem, []*model.VersionItem, []*model.PlatformItem, *custom_errors.ErrHttp) {
	resp, err := http.Get("https://bdu.fstec.ru/files/documents/vullist.xlsx")
	if err != nil {
		return nil, nil, nil, custom_errors.New(http.StatusInternalServerError, "send request: "+err.Error())
	}

	if resp.StatusCode != 200 {
		return nil, nil, nil, custom_errors.New(http.StatusInternalServerError, "wrong response status: "+resp.Status)
	}

	f, err := os.Create("files/programs.xlsx")
	if err != nil {
		return nil, nil, nil, custom_errors.New(http.StatusInternalServerError, "create file: "+err.Error())
	}
	if _, err = io.Copy(f, resp.Body); err != nil {
		return nil, nil, nil, custom_errors.New(http.StatusInternalServerError, "copy to file: "+err.Error())
	}
	f.Close()
	resp.Body.Close()

	file, err := excelize.OpenFile("files/programs.xlsx")
	if err != nil {
		return nil, nil, nil, custom_errors.New(http.StatusInternalServerError, "open file: "+err.Error())
	}

	rows, err := file.GetRows("Sheet")
	if err != nil {
		return nil, nil, nil, custom_errors.New(http.StatusInternalServerError, "get rows: "+err.Error())
	}

	var skipped int

	programsAndVersions := make(map[string][]string)
	programsAndPlatforms := make(map[string][]string)
	for idxRow, row := range rows {
		if idxRow == 0 || idxRow == 1 || idxRow == 2 {
			continue
		}

		versions := strings.Split(row[5], ", ")
		for _, version := range versions {
			versionAndName := strings.Split(version, " (")
			if len(versionAndName) != 2 {
				skipped++
				continue
			}
			var minVersion, maxVersion, justVersion string
			programName := strings.Trim(versionAndName[1], ")")
			if strings.Contains(versionAndName[0], "от") && strings.Contains(versionAndName[0], "до") {
				if len(strings.Split(versionAndName[0], " до ")) != 1 {
					minVersion, maxVersion = strings.Split(versionAndName[0], " до ")[0], strings.Split(versionAndName[0], " до ")[1]
					minVersion = strings.ReplaceAll(minVersion, "от ", "")
					minVersion = strings.ReplaceAll(minVersion, "от", "")
				}
			} else if strings.Contains(versionAndName[0], "до") {
				maxVersion = strings.ReplaceAll(versionAndName[0], "до ", "")
				maxVersion = strings.ReplaceAll(maxVersion, "до", "")
			} else {
				justVersion = versionAndName[0]
			}

			maxVersion = strings.ReplaceAll(maxVersion, " включительно ", "")
			maxVersion = strings.ReplaceAll(maxVersion, " включительно", "")
			maxVersion = strings.ReplaceAll(maxVersion, "включительно ", "")
			maxVersion = strings.ReplaceAll(maxVersion, "включительно", "")

			if minVersion == "" && maxVersion == "" && versionAndName[0] != "-" {
				justVersion = versionAndName[0]
			}

			var versionsStr []string
			if maxVersion == "" && minVersion != "" {
				versionsStr = append(versionsStr, minVersion)
			} else if maxVersion != "" && minVersion == "" {
				versionsStr = append(versionsStr, maxVersion)
			} else if justVersion != "" {
				versionsStr = append(versionsStr, justVersion)
			} else {
				minVersionArrStr := strings.Split(minVersion, ".")
				maxVersionArrStr := strings.Split(maxVersion, ".")
				var minVersionArr, maxVersionArr []int

				notBig := true
				onlyNumbers := true
				for _, elem := range minVersionArrStr {
					elemNumber, err := strconv.Atoi(elem)
					if err != nil {
						onlyNumbers = false
					} else {
						minVersionArr = append(minVersionArr, elemNumber)
					}
				}
				for _, elem := range maxVersionArrStr {
					elemNumber, err := strconv.Atoi(elem)
					if err != nil {
						onlyNumbers = false
					} else {
						maxVersionArr = append(maxVersionArr, elemNumber)
					}
				}

				if len(minVersionArr) != len(maxVersionArr) {
					if len(minVersionArr) < len(maxVersionArr) {
						for len(minVersionArr) != len(maxVersionArr) {
							minVersionArr = append(minVersionArr, 0)
						}
					} else {
						for len(minVersionArr) != len(maxVersionArr) {
							maxVersionArr = append(maxVersionArr, 0)
						}
					}
				}

				if onlyNumbers && maxVersionArr[len(maxVersionArr)-1]-minVersionArr[len(minVersionArr)-1] > 20 {
					notBig = false
				}

				if !onlyNumbers || !notBig {
					versionsStr = append(versionsStr, minVersion, maxVersion)
				} else {
					if len(maxVersionArr) == 1 || minVersionArr[len(minVersionArr)-2] == maxVersionArr[len(maxVersionArr)-2] {
						for i := minVersionArr[len(minVersionArr)-1]; i <= maxVersionArr[len(maxVersionArr)-1]; i++ {
							intArr := append(minVersionArr[:len(minVersionArr)-1], i)
							var strArr []string
							for _, elem := range intArr {
								strArr = append(strArr, strconv.Itoa(elem))
							}
							versionsStr = append(versionsStr, strings.Join(strArr, "."))
						}
					} else {
						versionsStr = append(versionsStr, minVersion)
						for i := 0; i <= maxVersionArr[len(maxVersionArr)-1]; i++ {
							intArr := append(maxVersionArr[:len(maxVersionArr)-1], i)
							var strArr []string
							for _, elem := range intArr {
								strArr = append(strArr, strconv.Itoa(elem))
							}
							versionsStr = append(versionsStr, strings.Join(strArr, "."))
						}
					}

				}
			}

			if len(strings.Split(row[4], ", ")) == 1 {
				platforms := strings.Split(row[7], ", ")
				for _, elem := range platforms {
					programsAndPlatforms[programName] = append(programsAndPlatforms[programName], strings.Split(elem, " (")[0])
				}
			}

			programsAndVersions[programName] = append(programsAndVersions[programName], versionsStr...)
		}
	}

	for key, value := range programsAndVersions {
		var uniqArr []string
		for _, elem := range value {
			var exist bool
			for _, uniqElem := range uniqArr {
				if elem == uniqElem {
					exist = true
				}
			}
			if !exist {
				uniqArr = append(uniqArr, elem)
			}
		}

		programsAndVersions[key] = uniqArr
	}

	for key, value := range programsAndPlatforms {
		var uniqArr []string
		for _, elem := range value {
			var exist bool
			for _, uniqElem := range uniqArr {
				if elem == uniqElem {
					exist = true
				}
			}
			if !exist {
				uniqArr = append(uniqArr, elem)
			}
		}

		programsAndPlatforms[key] = uniqArr
	}

	var programRes []*model.SoftwareProgramItem
	var versionsRes []*model.VersionItem
	var platformRes []*model.PlatformItem
	i := int64(1)
	for key, value := range programsAndVersions {
		programRes = append(programRes, &model.SoftwareProgramItem{
			ID:   i,
			Name: key,
		})

		for _, elem := range value {
			if elem != "" && elem != "-" {
				versionsRes = append(versionsRes, &model.VersionItem{
					ID:                1,
					Name:              elem,
					SoftwareProgramID: i,
				})
			}
		}

		if _, ok := programsAndPlatforms[key]; ok {
			for _, elem := range programsAndPlatforms[key] {
				if elem != "" && elem != "-" {
					platformRes = append(platformRes, &model.PlatformItem{
						ID:                1,
						Name:              elem,
						SoftwareProgramID: i,
					})
				}
			}
		}

		i++
	}

	return programRes, versionsRes, platformRes, nil
}

func ParseVul() ([]*model.VulnerabilityItem, *custom_errors.ErrHttp) {
	resp, err := http.Get("https://bdu.fstec.ru/files/documents/vullist.xlsx")
	if err != nil {
		return nil, custom_errors.New(http.StatusInternalServerError, "send request: "+err.Error())
	}

	if resp.StatusCode != 200 {
		return nil, custom_errors.New(http.StatusInternalServerError, "wrong response status: "+resp.Status)
	}

	f, err := os.Create("files/vullist.xlsx")
	if err != nil {
		return nil, custom_errors.New(http.StatusInternalServerError, "create file: "+err.Error())
	}
	if _, err = io.Copy(f, resp.Body); err != nil {
		return nil, custom_errors.New(http.StatusInternalServerError, "copy to file: "+err.Error())
	}
	f.Close()
	resp.Body.Close()

	file, err := excelize.OpenFile("files/vullist.xlsx")
	if err != nil {
		return nil, custom_errors.New(http.StatusInternalServerError, "open file: "+err.Error())
	}

	rows, err := file.GetRows("Sheet")
	if err != nil {
		return nil, custom_errors.New(http.StatusInternalServerError, "get rows: "+err.Error())
	}

	var vul []*model.VulnerabilityItem

	for idxRow, row := range rows {
		if idxRow == 0 || idxRow == 1 || idxRow == 2 {
			continue
		}

		var article string
		if len(strings.Split(row[18], ", ")) > 0 {
			article = strings.Split(row[18], ", ")[0]
		}

		var cweID string
		if len(row) >= 25 {
			cweID = row[24]
		}

		vul = append(vul, &model.VulnerabilityItem{
			ID:              1,
			BDUID:           row[0],
			Name:            row[1],
			Description:     row[2],
			RegDate:         row[9],
			DangerLevel:     row[12],
			Recommendations: row[13],
			Actuality:       row[16],
			ArticleNumber:   article,
			CWEID:           cweID,
			Component:       row[4],
			Platform:        row[7],
			Environment:     row[5],
		})
	}

	return vul, nil
}
