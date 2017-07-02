package leveldb

import (
	"errors"

	"github.com/sethjback/gobl/config"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

type key int

const Config key = 0

const (
	keyTypeAgent         = "ag-"
	keyTypeJobDefinition = "jd-"
	keyTypeJob           = "jb-"
	keyTypeFile          = "fl-"
	keyTypeFileDir       = "fd-"
	keyTypeIndex         = "in-"
	keyTypeShedule       = "sc-"
	keyTypeUser          = "us-"
	keyTypeKey           = "Key-"
)

type Leveldb struct {
	Connection *leveldb.DB
}

type dbConfig struct {
	path string
}

func (l *Leveldb) SaveConfig(cs config.Store, env map[string]string) error {
	dbc := &dbConfig{}
	for k, v := range env {
		switch k {
		case "DB_PATH":
			dbc.path = v
		}
	}

	cs.Add(Config, dbc)

	return nil
}

func configFromStore(cs config.Store) *dbConfig {
	if dbc, ok := cs.Get(Config); ok {
		return dbc.(*dbConfig)
	}
	return nil
}

func (l *Leveldb) Init(cs config.Store) error {
	dbc := configFromStore(cs)
	if dbc == nil {
		return errors.New("Unable to find leveldb config")
	}

	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}

	var err error

	if len(dbc.path) == 0 {
		l.Connection, err = leveldb.Open(storage.NewMemStorage(), o)
	} else {
		l.Connection, err = leveldb.OpenFile(dbc.path, o)
	}

	return err
}

func New(path string) (*Leveldb, error) {
	l := &Leveldb{}
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}

	var err error

	if len(path) == 0 {
		l.Connection, err = leveldb.Open(storage.NewMemStorage(), o)
	} else {
		l.Connection, err = leveldb.OpenFile(path, o)
	}

	return l, err
}

func (l *Leveldb) Close() error {
	return l.Connection.Close()
}
