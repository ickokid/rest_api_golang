package db

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
)

type DBConnection struct {
	session *mgo.Session
}

func NewConnection() (conn *DBConnection) {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB_HOST := os.Getenv("DB_HOST_MONGO")
	DB_USER := os.Getenv("DB_USER_MONGO")
	DB_PASSWORD := os.Getenv("DB_PASSWORD_MONGO")
	DB_DATABASE := os.Getenv("DB_DATABASE_MONGO")

	//var session, err := mgo.Dial("mongodb://golang:golang123@localhost/belajar_golang")
	info := &mgo.DialInfo{
		Addrs:     []string{DB_HOST},
		Timeout:   5 * time.Second,
		Database:  DB_DATABASE,
		Username:  DB_USER,
		Password:  DB_PASSWORD,
		Mechanism: "SCRAM-SHA-1",
	}
	session, err := mgo.DialWithInfo(info)
	checkErr(err, "Connection failed")
	conn = &DBConnection{session}

	return conn
}

func (conn *DBConnection) Use(dbName, tableName string) (collection *mgo.Collection) {
	return conn.session.DB(dbName).C(tableName)

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
