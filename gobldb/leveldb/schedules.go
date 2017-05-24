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

func (l *Leveldb) SaveSchedule(s model.Schedule) error {
	sbyte, err := json.Marshal(s)
	if err != nil {
		return goblerr.New("Unable to save schedule", errors.ErrCodeMarshal, err)
	}
	err = l.Connection.Put([]byte(keyTypeShedule+s.ID), sbyte, nil)
	if err != nil {
		return goblerr.New("Unable to save schedule", errors.ErrCodeSave, err)
	}
	return nil
}

func (l *Leveldb) GetSchedule(id string) (*model.Schedule, error) {
	abyte, err := l.Connection.Get([]byte(keyTypeShedule+id), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, goblerr.New("No schedule with that ID", errors.ErrCodeNotFound, err)
		}
		return nil, goblerr.New("Unable to get schedule", errors.ErrCodeGet, err)
	}

	var s model.Schedule
	err = json.Unmarshal(abyte, &s)
	if err != nil {
		return nil, goblerr.New("Unable to get schedule", errors.ErrCodeUnMarshal, err)
	}

	return &s, nil
}

func (l *Leveldb) DeleteSchedule(id string) error {
	err := l.Connection.Delete([]byte(keyTypeShedule+id), nil)
	if err != nil {
		return goblerr.New("Unable to delete schedule", errors.ErrCodeDelete, err)
	}
	return nil
}

func (l *Leveldb) ScheduleList() ([]model.Schedule, error) {
	var slist []model.Schedule
	iter := l.Connection.NewIterator(util.BytesPrefix([]byte(keyTypeShedule)), nil)
	for iter.Next() {
		var s model.Schedule
		err := json.Unmarshal(iter.Value(), &s)
		if err != nil {
			log.Errorf("leveldb", "Unable to unmarshal schedule: %+v", err)
			continue
		}
		slist = append(slist, s)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, goblerr.New("Unable to get schedule list", errors.ErrCodeGet, err)
	}
	return slist, nil
}
