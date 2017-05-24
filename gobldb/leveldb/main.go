package leveldb

import (
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/util/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

const (
	keyTypeAgent         = "ag-"
	keyTypeJobDefinition = "jd-"
	keyTypeJob           = "jb-"
	keyTypeFile          = "fl-"
	keyTypeFileDir       = "fd-"
	keyTypeIndex         = "in-"
	keyTypeShedule       = "sc-"
	keyTypeUser          = "us-"
)

type Leveldb struct {
	Connection *leveldb.DB
}

func New(options config.DB) (*Leveldb, error) {
	l := &Leveldb{}
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}

	var err error

	if len(options.Path) == 0 {
		log.Warn("leveldb", "DB Path empty, this will create in-memory db: probably not what you wanted!")
		l.Connection, err = leveldb.Open(storage.NewMemStorage(), o)
	} else {
		l.Connection, err = leveldb.OpenFile(options.Path, o)
	}

	return l, err
}

func (l *Leveldb) Close() error {
	return l.Connection.Close()
}
