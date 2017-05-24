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

func (l *Leveldb) SaveJobDefinition(jd model.JobDefinition) error {
	jdbyte, err := json.Marshal(jd)
	if err != nil {
		return goblerr.New("Unable to save job defintion", errors.ErrCodeMarshal, err)
	}
	err = l.Connection.Put([]byte(keyTypeJobDefinition+jd.ID), jdbyte, nil)
	if err != nil {
		return goblerr.New("Unable to save jobdefinition", errors.ErrCodeSave, err)
	}
	return nil
}

func (l *Leveldb) GetJobDefinition(id string) (*model.JobDefinition, error) {
	jdbyte, err := l.Connection.Get([]byte(keyTypeJobDefinition+id), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, goblerr.New("No job definition with that ID", errors.ErrCodeNotFound, err)
		}
		return nil, goblerr.New("Unable to get job definition", errors.ErrCodeGet, err)
	}

	var jd model.JobDefinition
	err = json.Unmarshal(jdbyte, &jd)
	if err != nil {
		return nil, goblerr.New("Unable to get job definition", errors.ErrCodeUnMarshal, err)
	}

	return &jd, nil
}

func (l *Leveldb) DeleteJobDefinition(id string) error {
	err := l.Connection.Delete([]byte(keyTypeJobDefinition+id), nil)
	if err != nil {
		return goblerr.New("Unable to delete job definition", errors.ErrCodeDelete, err)
	}
	return nil
}

func (l *Leveldb) JobDefinitionList() ([]model.JobDefinition, error) {
	var jdlist []model.JobDefinition
	iter := l.Connection.NewIterator(util.BytesPrefix([]byte(keyTypeJobDefinition)), nil)
	for iter.Next() {
		var jd model.JobDefinition
		err := json.Unmarshal(iter.Value(), &jd)
		if err != nil {
			log.Errorf("leveldb", "Unable to unmarshal job definition: %+v", err)
			continue
		}
		jdlist = append(jdlist, jd)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, goblerr.New("Unable to get job definition list", errors.ErrCodeGet, err)
	}
	return jdlist, nil
}
