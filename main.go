package main

import (
	api_db "diplom-vuls-server/internal/api-db"
	"diplom-vuls-server/internal/model"
	consruct_namespaces "diplom-vuls-server/pkg/consruct-namespaces"
	"diplom-vuls-server/pkg/parser"
	"encoding/json"
	"errors"
	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"os"
)

var db *reindexer.Reindexer

var config model.Config

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	data, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Error().Err(errors.New("file open: " + err.Error())).Msg("")
	}

	if err = json.Unmarshal(data, &config); err != nil {
		log.Error().Err(errors.New("unmarshal config: " + err.Error())).Msg("")
	}

	db = reindexer.NewReindex(config.DB.Scheme + "://" + config.DB.Hostname + ":" + config.DB.Port + "/" + config.DB.Path)
	if err = db.Status().Err; err != nil {
		log.Error().Err(errors.New("reindexer connection: " + err.Error())).Msg("")
	}

	log.Info().Msg("Connection to diplom_vuls reindexer DB successful!")

	consruct_namespaces.ConstructNamespaces(db)

	server := &fasthttp.Server{
		Handler: initRoutes,
	}

	log.Info().Msg("Начало парсинга БДУ")

	programs, versions, platforms, errCustom := parser.ParseSoftwareProgram()
	if errCustom != nil {
		log.Error().Err(errors.New("parse programs: " + errCustom.Error())).Msg("")
	}

	for _, program := range programs {
		if exist := api_db.CheckExistingProgramDB(program, db); !exist {
			if errCustom = api_db.CreateSoftwareProgramDB(program, db); errCustom != nil {
				log.Error().Err(errors.New("create software program DB: " + errCustom.Error())).Msg("")
			}
		}
	}

	for _, version := range versions {
		if exist := api_db.CheckExistingVersionDB(version, db); !exist {
			if errCustom = api_db.CreateVersionDB(version, db); errCustom != nil {
				log.Error().Err(errors.New("create create version DB: " + errCustom.Error())).Msg("")

			}
		}
	}

	for _, platform := range platforms {
		if exist := api_db.CheckExistingPlatformDB(platform, db); !exist {
			if errCustom = api_db.CreatePlatformDB(platform, db); errCustom != nil {
				log.Error().Err(errors.New("create create platform DB: " + errCustom.Error())).Msg("")

			}
		}
	}

	vuls, errCustom := parser.ParseVul()
	if errCustom != nil {
		log.Error().Err(errors.New("parse vuls: " + errCustom.Error())).Msg("")
	}

	for _, vul := range vuls {
		if exist := api_db.CheckExistingVulnerabilityDB(vul, db); !exist {
			if errCustom = api_db.CreateVulDB(vul, db); errCustom != nil {
				log.Error().Err(errors.New("create create vulnerability DB: " + errCustom.Error())).Msg("")
			}
		}
	}
	log.Info().Msg("Парсинг окончен")

	c := cron.New()
	c.AddFunc("@every dat at 1am", func() {
		log.Info().Msg("Начало парсинга БДУ")
		programs, versions, platforms, errCustom := parser.ParseSoftwareProgram()
		if errCustom != nil {
			log.Error().Err(errors.New("parse programs: " + errCustom.Error())).Msg("")
		}

		for _, program := range programs {
			if exist := api_db.CheckExistingProgramDB(program, db); !exist {
				if errCustom = api_db.CreateSoftwareProgramDB(program, db); errCustom != nil {
					log.Error().Err(errors.New("create software program DB: " + errCustom.Error())).Msg("")
				}
			}
		}

		for _, version := range versions {
			if exist := api_db.CheckExistingVersionDB(version, db); !exist {
				if errCustom = api_db.CreateVersionDB(version, db); errCustom != nil {
					log.Error().Err(errors.New("create create version DB: " + errCustom.Error())).Msg("")

				}
			}
		}

		for _, platform := range platforms {
			if exist := api_db.CheckExistingPlatformDB(platform, db); !exist {
				if errCustom = api_db.CreatePlatformDB(platform, db); errCustom != nil {
					log.Error().Err(errors.New("create create platform DB: " + errCustom.Error())).Msg("")

				}
			}
		}

		vuls, errCustom := parser.ParseVul()
		if errCustom != nil {
			log.Error().Err(errors.New("parse vuls: " + errCustom.Error())).Msg("")
		}

		for _, vul := range vuls {
			if exist := api_db.CheckExistingVulnerabilityDB(vul, db); !exist {
				if errCustom = api_db.CreateVulDB(vul, db); errCustom != nil {
					log.Error().Err(errors.New("create create vulnerability DB: " + errCustom.Error())).Msg("")
				}
			}
		}
		log.Info().Msg("Парсинг окончен")
	})
	c.Start()

	if err = server.ListenAndServe(config.HTTP.Port); err != nil {
		log.Error().Err(errors.New("start server: " + err.Error())).Msg("")
	}
}

func initRoutes(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	log.Info().Str("path", string(ctx.Path())).Str("method", string(ctx.Method())).Msg("")
	switch string(ctx.Path()) {
	case "/check":
		CheckHandler(ctx)

	default:
		ctx.Error("Page not found", fasthttp.StatusNotFound)
	}
}
