package db

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

type DBConnection struct {
	session *mgo.Session
}

var (
	config = Config{}
)

func NewConnection() (conn *DBConnection) {
	config.Read("config.toml")

	//var session, err := mgo.Dial("mongodb://golang:golang123@localhost/belajar_golang")
	info := &mgo.DialInfo{
		Addrs:     config.Mongo_addrs,
		Timeout:   60 * time.Second,
		Database:  config.Mongo_database,
		Username:  config.Mongo_username,
		Password:  config.Mongo_password,
		Mechanism: "SCRAM-SHA-1",
	}
	session, err := mgo.DialWithInfo(info)
	checkErr(err, "Connection failed")
	conn = &DBConnection{session}

	return conn
}

func (conn *DBConnection) Use(tableName string) (collection *mgo.Collection) {
	config.Read("config.toml")

	return conn.session.DB(config.Mongo_database).C(tableName)

}

func (conn *DBConnection) Close() {
	conn.session.Close()
	return
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
