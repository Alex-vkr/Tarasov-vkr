package model

type PlatformItem struct {
	ID                int64    `json:"id" reindex:"id,hash,pk"`
	Name              string   `json:"name" reindex:"name,hash"`
	SoftwareProgramID int64    `json:"software_program_id" reindex:"software_program_id,hash"`
	_                 struct{} `reindex:"name=search,text,composite"`
}
