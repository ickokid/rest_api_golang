package models

import "gopkg.in/mgo.v2/bson"

const (
	DB_HOST     = string("localhost:3307")
	DB_USER     = string("root")
	DB_PASSWORD = string("root")
	DB_DATABASE = string("db_belajar_golang")
)

type Student struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

type Response struct {
	Status  uint8     `json:"status"`
	Message string    `json:"message"`
	Data    []Student `json:"data"`
}

const (
	DB_HOST_MONGO     = string("localhost:27017") //mongodb://127.0.0.1:27017
	DB_USER_MONGO     = string("golang")
	DB_PASSWORD_MONGO = string("golang123")
	DB_DATABASE_MONGO = string("belajar_golang")
)

type Student_Mongo struct {
	Id    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name  string        `json:"name"`
	Grade int           `json:"grade"`
}

type Response_Mongo struct {
	Status  uint8           `json:"status"`
	Message string          `json:"message"`
	Data    []Student_Mongo `json:"data"`
}
