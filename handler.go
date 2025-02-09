package main

import (
	api_db "diplom-vuls-server/internal/api-db"
	"diplom-vuls-server/internal/model"
	"diplom-vuls-server/pkg/check"
	string_builder "diplom-vuls-server/pkg/string-builder"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/tealeg/xlsx"
	"github.com/valyala/fasthttp"
	"time"
)

func CheckHandler(ctx *fasthttp.RequestCtx) (message string, code int) {
	defer func() {
		resFinal := model.Response{Data: nil, Code: code, Message: message}
		if code != fasthttp.StatusOK {
			log.Error().Err(errors.New(message)).Msg("")
		}
		jsonRes, _ := json.Marshal(resFinal)
		ctx.Response.SetStatusCode(resFinal.Code)
		ctx.Response.SetBody(jsonRes)
	}()

	log.Info().Msg("Проверка машины...")

	if !ctx.IsPost() {
		return "handler: wrong method", fasthttp.StatusMethodNotAllowed
	}

	var req model.CheckReq

	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return "handler: unmarshal request: " + err.Error(), fasthttp.StatusUnprocessableEntity
	}

	ip := ctx.RemoteIP().String()

	var programIDs, versionIDs, platformIDs []int64
	for idx, program := range req.Programs {
		programID, errCustom := api_db.CheckProgramDB(program, db)
		if errCustom != nil {
			return "handler: check program: " + errCustom.Error(), errCustom.Code
		}
		programIDs = append(programIDs, programID)

		versionID, errCustom := api_db.CheckVersionDB(req.Versions[idx], programID, db)
		if errCustom != nil {
			return "handler: check version: " + errCustom.Error(), errCustom.Code
		}
		versionIDs = append(versionIDs, versionID)

		platformID, errCustom := api_db.CheckPlatformDB(req.Platform, programID, db)
		if errCustom != nil {
			return "handler: check platform: " + errCustom.Error(), errCustom.Code
		}
		platformIDs = append(platformIDs, platformID)
	}

	fileName := string_builder.BuildStrings([]string{"files/Отчет_от_", time.Now().Format("02.01.2006"), "_адрес машины_", ip, ".xlsx"})

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Отчет")
	if err != nil {
		return "handler: add sheet: " + err.Error(), fasthttp.StatusInternalServerError
	}
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Merge(12, 1)
	cell.SetString(fileName)

	vuls, errCustom := check.CheckVul(programIDs, versionIDs, platformIDs, db)
	if errCustom != nil {
		return "handler: check vul: " + errCustom.Error(), errCustom.Code
	}

	for idx, programVuls := range vuls {
		if len(programVuls) == 0 {
			continue
		}

		row = sheet.AddRow()
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Merge(12, 1)
		cell.SetString(string_builder.BuildStrings([]string{"Наименование ПО: ", req.Programs[idx], ". Версия ПО: ", req.Versions[idx]}))

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Merge(12, 1)
		cell.SetString("Возможные уязвимости:")

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetString("Идентификатор БДУ")
		cell = row.AddCell()
		cell.SetString("Наименование уязвимости")
		cell = row.AddCell()
		cell.SetString("Описание уязвимости")
		cell = row.AddCell()
		cell.SetString("Название ПО")
		cell = row.AddCell()
		cell.SetString("Версия ПО")
		cell = row.AddCell()
		cell.SetString("Платформа")
		cell = row.AddCell()
		cell.SetString("Актуальность")
		cell = row.AddCell()
		cell.SetString("Дата выявления")
		cell = row.AddCell()
		cell.SetString("Уровень опасности")
		cell = row.AddCell()
		cell.SetString("Рекомендации")
		cell = row.AddCell()
		cell.SetString("Идентификаторы других систем")
		cell = row.AddCell()
		cell.SetString("Тип ошибки CWE")

		for _, vul := range programVuls {
			row = sheet.AddRow()
			cell = row.AddCell()
			cell.SetString(vul.BDUID)
			cell = row.AddCell()
			cell.SetString(vul.Name)
			cell = row.AddCell()
			cell.SetString(vul.Description)
			cell = row.AddCell()
			cell.SetString(vul.Component)
			cell = row.AddCell()
			cell.SetString(vul.Environment)
			cell = row.AddCell()
			cell.SetString(vul.Platform)
			cell = row.AddCell()
			cell.SetString(vul.Actuality)
			cell = row.AddCell()
			cell.SetString(vul.RegDate)
			cell = row.AddCell()
			cell.SetString(vul.DangerLevel)
			cell = row.AddCell()
			cell.SetString(vul.Recommendations)
			cell = row.AddCell()
			cell.SetString(vul.ArticleNumber)
			cell = row.AddCell()
			cell.SetString(vul.CWEID)

		}
	}

	if err = file.Save(fileName); err != nil {
		return "handler: save file: " + err.Error(), fasthttp.StatusInternalServerError
	}

	log.Info().Msg("Проверка машины окончена. Отчет сформирован.")

	return "OK", fasthttp.StatusOK
}
