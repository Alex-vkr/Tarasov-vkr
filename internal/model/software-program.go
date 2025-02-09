package model

type SoftwareProgramItem struct {
	ID   int64    `json:"id" reindex:"id,hash,pk"`
	Name string   `json:"name" reindex:"name,hash"`
	_    struct{} `reindex:"name=search,text,composite"`
}
