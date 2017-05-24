package leveldb

import (
	"encoding/json"

	"github.com/sethjback/gobl/gobldb/errors"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/util/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func (l *Leveldb) SaveAgent(a model.Agent) error {
	abyte, err := json.Marshal(a)
	if err != nil {
		return goblerr.New("Unable to save agent", errors.ErrCodeMarshal, err)
	}
	err = l.Connection.Put([]byte(keyTypeAgent+a.ID), abyte, nil)
	if err != nil {
		return goblerr.New("Unable to save agent", errors.ErrCodeSave, err)
	}
	return nil
}

func (l *Leveldb) GetAgent(id string) (*model.Agent, error) {
	abyte, err := l.Connection.Get([]byte(keyTypeAgent+id), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, goblerr.New("No agent with that ID", errors.ErrCodeNotFound, err)
		}
		return nil, goblerr.New("Unable to get agent", errors.ErrCodeGet, err)
	}

	var a model.Agent
	err = json.Unmarshal(abyte, &a)
	if err != nil {
		return nil, goblerr.New("Unable to get agent", errors.ErrCodeUnMarshal, err)
	}

	return &a, nil
}

func (l *Leveldb) AgentList() ([]model.Agent, error) {
	var alist []model.Agent
	iter := l.Connection.NewIterator(util.BytesPrefix([]byte(keyTypeAgent)), nil)
	for iter.Next() {
		var a model.Agent
		err := json.Unmarshal(iter.Value(), &a)
		if err != nil {
			log.Errorf("leveldb", "Unable to unmarshal agent: %+v", err)
			continue
		}
		alist = append(alist, a)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, goblerr.New("Unable to get agent list", errors.ErrCodeGet, err)
	}
	return alist, nil
}
