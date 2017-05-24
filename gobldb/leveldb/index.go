package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	indexTypeFileState  = "fs-"
	indexTypeFileParent = "fp-"
	indexTypeJobDate    = "jd-"
	indexTypeJobState   = "js-"
	indexTypeJobAgent   = "ja-"
)

type index struct {
	itype string
	key   string
	value string
}

func (l *Leveldb) NewIndex(i index) error {
	return l.Connection.Put(append([]byte(keyTypeIndex+i.itype), []byte(i.key)...), []byte(i.value), nil)
}

func (l *Leveldb) GetIndex(itype string, key string) (*index, error) {
	value, err := l.Connection.Get(append([]byte(keyTypeIndex+itype), key...), nil)
	if err != nil {
		return nil, err
	}
	return &index{itype, key, string(value)}, nil
}

func (l *Leveldb) DeleteIndex(i index) error {
	return l.Connection.Delete([]byte(keyTypeIndex+i.itype+i.key), nil)
}

func indexBatchPut(is ...index) *leveldb.Batch {
	b := new(leveldb.Batch)
	for _, i := range is {
		b.Put(append([]byte(keyTypeIndex+i.itype), []byte(i.key)...), []byte(i.value))
	}
	return b
}

func (l *Leveldb) indexQuery(itype string, prefix string) ([]index, error) {
	var is []index
	iter := l.Connection.NewIterator(util.BytesPrefix(append([]byte(keyTypeIndex+itype), prefix...)), nil)
	for iter.Next() {
		is = append(is, index{itype, string(iter.Key()[len(keyTypeIndex)+len(itype):]), string(iter.Value())})
	}
	iter.Release()
	return is, iter.Error()
}

func (l *Leveldb) indexRange(itype string, start, limit string) ([]index, error) {
	var is []index

	iter := l.Connection.NewIterator(&util.Range{
		Start: append([]byte(keyTypeIndex+itype), start...),
		Limit: append([]byte(keyTypeIndex+itype), limit...)}, nil)

	for iter.Next() {
		is = append(is, index{itype, string(iter.Key()[len(keyTypeIndex)+len(itype):]), string(iter.Value())})
	}
	iter.Release()
	return is, iter.Error()
}
