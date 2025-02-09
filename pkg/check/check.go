package check

import (
	"diplom-vuls-server/internal/model"
	custom_errors "diplom-vuls-server/pkg/custom-errors"
	"errors"
	"github.com/restream/reindexer/v3"
	"github.com/rs/zerolog/log"
	"strings"
)

func CheckVul(programIDs, versionIDs, platformIDs []int64, db *reindexer.Reindexer) ([][]*model.VulnerabilityItem, *custom_errors.ErrHttp) {
	var res [][]*model.VulnerabilityItem
	for idxProgram, programID := range programIDs {
		var vuls []*model.VulnerabilityItem
		rec, found := db.Query("software_program").WhereInt64("id", reindexer.EQ, programID).Get()
		if !found {
			log.Error().Err(errors.New("segment check: check vul: software program with this id doesn't exist")).Msg("")
			continue
		}

		programRec := rec.(*model.SoftwareProgramItem)

		rec, found = db.Query("version").WhereInt64("id", reindexer.EQ, versionIDs[idxProgram]).Get()
		if !found {
			log.Error().Err(errors.New("segment check: check vul: version with this id doesn't exist")).Msg("")
			continue
		}

		versionRec := rec.(*model.VersionItem)

		rec, found = db.Query("platform").WhereInt64("id", reindexer.EQ, platformIDs[idxProgram]).Get()
		if !found {
			log.Error().Err(errors.New("segment check: check vul: platform with this id doesn't exist")).Msg("")
			continue
		}

		platformRec := rec.(*model.PlatformItem)

		if versionRec.SoftwareProgramID != versionRec.SoftwareProgramID {
			log.Error().Err(errors.New("segment check: check vul: wrong version")).Msg("")
			continue
		}

		if platformRec.SoftwareProgramID != platformRec.SoftwareProgramID {
			log.Error().Err(errors.New("segment check: check vul: wrong platform")).Msg("")
			continue
		}

		iterator := db.Query("vulnerability").Match("component", programRec.Name).Exec()
		if iterator.Error() != nil {
			log.Error().Err(errors.New("segment check: check vul: iterator vuls: " + iterator.Error().Error())).Msg("")
			continue
		}

		for iterator.Next() {
			vul := iterator.Object().(*model.VulnerabilityItem)
			if strings.Contains(vul.Component, programRec.Name) && strings.Contains(vul.Platform, platformRec.Name) {
				versions := strings.Split(vul.Environment, ", ")

				var add bool

				for _, version := range versions {
					versionAndName := strings.Split(version, " (")
					if len(versionAndName) != 2 {
						continue
					}

					if !(strings.Trim(versionAndName[1], ")") == programRec.Name) {
						continue
					}

					var minVersion, maxVersion, justVersion string
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

					if minVersion != "" && maxVersion != "" {
						if versionRec.Name > minVersion && versionRec.Name < maxVersion {
							add = true
						}
					} else if maxVersion != "" {
						if versionRec.Name < maxVersion {
							add = true
						}
					} else {
						if versionRec.Name == justVersion {
							add = true
						}
					}
				}
				if add {
					vuls = append(vuls, vul)
				}
			}
		}
		res = append(res, vuls)
	}

	return res, nil
}
