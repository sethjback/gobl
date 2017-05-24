package leveldb

import "github.com/sethjback/gobl/model"

type job struct {
	ID      string               `json:"id"`
	Def     *model.JobDefinition `json:"def"`
	AgentId string               `json:"agentid"`
	Meta    *model.JobMeta       `json:"meta"`
}
