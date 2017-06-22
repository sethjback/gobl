package leveldb

import (
	"encoding/json"

	"github.com/sethjback/gobl/gobldb/errors"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/model"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func (l *Leveldb) SaveUser(u model.User) error {
	ubyte, err := json.Marshal(u)
	if err != nil {
		return goblerr.New("Unable to save user", errors.ErrCodeMarshal, err)
	}
	err = l.Connection.Put([]byte(keyTypeUser+u.Email), ubyte, nil)
	if err != nil {
		return goblerr.New("Unable to save user", errors.ErrCodeSave, err)
	}
	return nil
}

func (l *Leveldb) GetUser(id string) (*model.User, error) {
	ubyte, err := l.Connection.Get([]byte(keyTypeUser+id), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, goblerr.New("No user with that email", errors.ErrCodeNotFound, err)
		}
		return nil, goblerr.New("Unable to get user", errors.ErrCodeGet, err)
	}

	var u model.User
	err = json.Unmarshal(ubyte, &u)
	if err != nil {
		return nil, goblerr.New("Unable to get user", errors.ErrCodeUnMarshal, err)
	}

	return &u, nil
}

func (l *Leveldb) UserList() ([]model.User, error) {
	var ulist []model.User
	iter := l.Connection.NewIterator(util.BytesPrefix([]byte(keyTypeUser)), nil)
	for iter.Next() {
		var u model.User
		err := json.Unmarshal(iter.Value(), &u)
		if err != nil {
			continue
		}
		ulist = append(ulist, u)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil, goblerr.New("Unable to get user list", errors.ErrCodeGet, err)
	}
	return ulist, nil
}

func (l *Leveldb) DeleteUser(id string) error {
	err := l.Connection.Delete([]byte(keyTypeUser+id), nil)
	if err != nil {
		return goblerr.New("Unable to delete user", errors.ErrCodeDelete, err)
	}
	return nil
}
