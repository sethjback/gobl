package db

import (
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/gobldb"
)

var gDb gobldb.Database

func Get() gobldb.Database {
	return gDb
}

func SaveConfig(cs config.Store, cMap map[string]string) error {
	return gobldb.SaveConfig(cs, cMap)
}

func Init(c config.Store) error {
	var err error
	gDb, err = gobldb.Get(c)
	return err
}

func Shutdown() {
	gDb.Close()
}
