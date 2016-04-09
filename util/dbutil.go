package util

import (
//	"encoding/json"
	"gopkg.in/mgo.v2"
	"github.com/larspensjo/config"
//	"gopkg.in/mgo.v2/bson"
)

func GetDB( inicfg *config.Config) *mgo.Database {
	ip,_:= inicfg.String("db", "ip")
	dbname,_:= inicfg.String("db", "db")
	session, err := mgo.Dial(ip)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB(dbname)
}

