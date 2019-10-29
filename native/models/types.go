package models

import "gopkg.in/mgo.v2/bson"

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
