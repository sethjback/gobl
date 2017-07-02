package leveldb

import (
	"encoding/json"

	"github.com/sethjback/gobl/gobldb/errors"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/model"
	"github.com/syndtr/goleveldb/leveldb"
)

func (l *Leveldb) SaveKey(kType string, caKey model.Key) error {
	abyte, err := json.Marshal(caKey)
	if err != nil {
		return goblerr.New("Unable to save CA key", errors.ErrCodeMarshal, err)
	}
	err = l.Connection.Put([]byte(keyTypeKey+kType), abyte, nil)

	if err != nil {
		return goblerr.New("Unable to save CA Key", errors.ErrCodeSave, err)
	}

	return nil
}

func (l *Leveldb) GetKey(kType string) (*model.Key, error) {
	abyte, err := l.Connection.Get([]byte(keyTypeKey+kType), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, goblerr.New("No saved key", errors.ErrCodeNotFound, err)
		}
		return nil, goblerr.New("Unable to get key", errors.ErrCodeGet, err)
	}

	var key model.Key
	err = json.Unmarshal(abyte, &key)
	if err != nil {
		return nil, goblerr.New("Unable to get key", errors.ErrCodeUnMarshal, err)
	}

	return &key, nil
}
