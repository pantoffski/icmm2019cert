package database

import (
	"icmm2019cert/cfg"

	"github.com/globalsign/mgo"
)

var mongoDbName string
var mongoSession *mgo.Session

// GetDB return db instalce
func GetDB() *mgo.Database {
	if mongoSession != nil {
		return mongoSession.Copy().DB(mongoDbName)
	}
	mongoURL := cfg.Getenv("CERT2019_MONGO_URL")
	mongoDbName := cfg.Getenv("CERT2019_MONGO_DB_NAME")
	mongoSession, err := mgo.Dial(mongoURL)
	if err != nil {
		panic("error connecting to mongo")
	}
	return mongoSession.Copy().DB(mongoDbName)
}
