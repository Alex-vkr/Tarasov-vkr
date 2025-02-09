package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

type Config struct {
	ServerDomain string `json:"server_domain"`
}

type CheckReq struct {
	Programs []string `json:"programs"`
	Versions []string `json:"versions"`
	Platform string   `json:"platform"`
}

func main() {
	var config Config
	var wg sync.WaitGroup

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Error().Err(errors.New("file open: " + err.Error())).Msg("")
	}

	if err = json.Unmarshal(data, &config); err != nil {
		log.Error().Err(errors.New("unmarshal config: " + err.Error())).Msg("")
	}

	Check(config)
	c := cron.New()
	c.AddFunc("24 15 * * *", func() {
		Check(config)
	})
	c.Start()

	wg.Add(1)

	go func() {

	}()

	wg.Wait()
}

func Check(config Config) {
	command := exec.Command("powershell.exe", "-Command", "Get-ItemProperty HKLM:\\Software\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\* | Select-Object DisplayName, DisplayVersion  | Export-Csv -Path .\\products.csv -NoTypeInformation -Encoding UTF8 -Delimiter ';'")

	_, err := command.Output()
	if err != nil {
		log.Error().Err(errors.New("command output: " + err.Error())).Msg("")
		return
	}

	var req CheckReq

	file, err := os.Open("products.csv")
	if err != nil {
		log.Error().Err(errors.New("open: " + err.Error())).Msg("")
	}

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true

	programs, err := reader.ReadAll()
	if err != nil {
		log.Error().Err(errors.New("read all: " + err.Error())).Msg("")
	}

	for idx, line := range programs {
		if idx == 0 {
			continue
		}

		if line[0] == "" && line[1] == "" {
			continue
		}

		req.Programs = append(req.Programs, line[0])
		req.Versions = append(req.Versions, line[1])
	}

	file.Close()

	req.Platform = "Windows"

	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Error().Err(errors.New("marshal: " + err.Error())).Msg("")
	}

	r := bytes.NewReader(reqBody)

	resp, err := http.Post(config.ServerDomain+"/check", "application/json", r)
	if err != nil {
		log.Error().Err(errors.New("http post: " + err.Error())).Msg("")
	}

	if resp.StatusCode != 200 {
		log.Error().Err(errors.New("response status: " + resp.Status)).Msg("")
	}

	log.Info().Msg("Сведения о машине отправлены.")
}
