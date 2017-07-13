package apihandler

import "github.com/sethjback/gobl/gobldb"

var db gobldb.Database

func Init(dbS gobldb.Database) {
	db = dbS
}
