package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gopkg.in/gorp.v1"
)

func InitDb() *gorp.DbMap {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB_HOST := os.Getenv("DB_HOST_MYSQL")
	DB_USER := os.Getenv("DB_USER_MYSQL")
	DB_PASSWORD := os.Getenv("DB_PASSWORD_MYSQL")
	DB_DATABASE := os.Getenv("DB_DATABASE_MYSQL")

	//db, err := sql.Open("mysql", "root:root@tcp(localhost:3307)/godb")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", DB_USER, DB_PASSWORD, DB_HOST, DB_DATABASE))
	checkErr(err, "sql.Open failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	//dbmap.AddTableWithName(Instruction{}, "instruction").SetKeys(true, "Id")

	err = db.Ping()
	checkErr(err, "db.Ping failed")

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create table failed")
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
