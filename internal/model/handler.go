package model

type Response struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CheckReq struct {
	Programs []string `json:"programs"`
	Versions []string `json:"versions"`
	Platform string   `json:"platform"`
}
