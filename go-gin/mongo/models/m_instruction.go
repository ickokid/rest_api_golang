package models

import (
	"rest-api-gin-mongo/db"

	"gopkg.in/mgo.v2/bson"
)

const (
	C_INSTRUCTION = "instruction"
)

type Instruction struct {
	Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	EventStatus string        `json:"event_status" bson:"event_status" form:"event_status" binding:"required"`
	EventName   string        `json:"event_name" bson:"event_name" form:"event_name" binding:"required"`
}

var dbConnect = db.NewConnection()

func FindAll(skip int, limit int) (instructions []Instruction, err error) {
	var selector = bson.M{}

	var collection = dbConnect.Use(C_INSTRUCTION)
	err = collection.Find(selector).Sort("-_id").Skip(skip).Limit(limit).All(&instructions)

	return instructions, err
}

func FindById(id string) (result *Instruction, err error) {
	var instruction Instruction

	var collection = dbConnect.Use(C_INSTRUCTION)
	err = collection.FindId(bson.ObjectIdHex(id)).Select(bson.M{"_id": 1, "event_status": 1, "event_name": 1}).One(&instruction)

	result = &Instruction{
		Id:          bson.ObjectIdHex(id),
		EventStatus: instruction.EventStatus,
		EventName:   instruction.EventName,
	}
	return result, err
}

func Insert(eventStatusId bson.ObjectId, eventStatus string, eventName string) (result *Instruction, err error) {
	var collection = dbConnect.Use(C_INSTRUCTION)
	err = collection.Insert(&Instruction{eventStatusId, eventStatus, eventName})

	result = &Instruction{
		Id:          eventStatusId,
		EventStatus: eventStatus,
		EventName:   eventName,
	}

	return result, err
}

func Update(eventStatusId string, eventStatus string, eventName string) (result *Instruction, err error) {
	var collection = dbConnect.Use(C_INSTRUCTION)
	var selector = bson.M{"_id": bson.ObjectIdHex(eventStatusId)}
	var changes = bson.M{"$set": bson.M{"event_status": eventStatus, "event_name": eventName}}
	err = collection.Update(selector, &changes)

	result = &Instruction{
		Id:          bson.ObjectIdHex(eventStatusId),
		EventStatus: eventStatus,
		EventName:   eventName,
	}

	return result, err
}

func Delete(eventStatusId string) (err error) {
	var collection = dbConnect.Use(C_INSTRUCTION)
	var selector = bson.M{"_id": bson.ObjectIdHex(eventStatusId)}
	err = collection.Remove(selector)

	return err
}
