package controllers

import (
	"rest-api-gin-mongo/db"

	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionInstruction = "instruction"
)

type Instruction struct {
	Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	EventStatus string        `json:"event_status" bson:"event_status"`
	EventName   string        `json:"event_name" bson:"event_name"`
}

var dbConnect = db.NewConnection()

func GetInstructions(c *gin.Context) {
	var instructions []Instruction
	var result gin.H
	var selector = bson.M{}

	var collection = dbConnect.Use("belajar_golang", CollectionInstruction)
	err := collection.Find(selector).All(&instructions)

	if err == nil {
		if len(instructions) <= 0 {
			result = gin.H{
				"status": 200,
				"data":   nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"status": 200,
				"data":   instructions,
				"count":  len(instructions),
			}
		}
	} else {
		result = gin.H{
			"status": 404,
			"error":  "no instruction(s) into the table",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i http://localhost:8080/api/v1/instructions
}

func GetInstruction(c *gin.Context) {
	id := c.Param("id")
	var instruction Instruction
	var result gin.H

	if bson.IsObjectIdHex(id) {
		var collection = dbConnect.Use("belajar_golang", CollectionInstruction)
		err := collection.FindId(bson.ObjectIdHex(id)).Select(bson.M{"_id": 1, "event_status": 1, "event_name": 1}).One(&instruction)

		if err == nil {
			content := &Instruction{
				Id:          bson.ObjectIdHex(id),
				EventStatus: instruction.EventStatus,
				EventName:   instruction.EventName,
			}

			result = gin.H{
				"status": 200,
				"data":   content,
				"count":  1,
			}
		} else {
			result = gin.H{
				"status": 404,
				"error":  "instruction not found",
			}
		}
	} else {
		result = gin.H{
			"status": 404,
			"error":  "invalif format",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i http://localhost:8080/api/v1/Instructions/1
}

func PostInstruction(c *gin.Context) {
	var instruction Instruction
	var result gin.H
	c.Bind(&instruction)

	if instruction.EventStatus != "" && instruction.EventName != "" {
		var collection = dbConnect.Use("belajar_golang", CollectionInstruction)
		instruction_id := bson.NewObjectId()
		err := collection.Insert(&Instruction{instruction_id, instruction.EventStatus, instruction.EventName})

		if err == nil {
			content := &Instruction{
				Id:          instruction_id,
				EventStatus: instruction.EventStatus,
				EventName:   instruction.EventName,
			}

			result = gin.H{
				"status": 200,
				"data":   content,
			}
		} else {
			result = gin.H{
				"status": 400,
				"error":  "Insert failed",
			}
		}
	} else {
		result = gin.H{
			"status": 400,
			"error":  "fields are empty",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"event_status\": \"83\", \"event_name\": \"100\" }" http://localhost:8080/api/v1/instructions
}

func UpdateInstruction(c *gin.Context) {
	id := c.Params.ByName("id")
	var instruction Instruction
	var result gin.H

	if bson.IsObjectIdHex(id) {
		var collection = dbConnect.Use("belajar_golang", CollectionInstruction)
		err := collection.FindId(bson.ObjectIdHex(id)).Select(bson.M{"_id": 1, "event_status": 1, "event_name": 1}).One(&instruction)

		if err == nil {
			var json Instruction
			c.Bind(&json)
			instruction := Instruction{
				Id:          bson.ObjectIdHex(id),
				EventStatus: json.EventStatus,
				EventName:   json.EventName,
			}

			if instruction.EventStatus != "" && instruction.EventName != "" {
				var selector = bson.M{"_id": bson.ObjectIdHex(id)}
				err := collection.Update(selector, &instruction)

				if err == nil {
					result = gin.H{
						"status": 200,
						"data":   instruction,
					}
				} else {
					result = gin.H{
						"status": 400,
						"error":  "Updated failed",
					}
				}
			} else {
				result = gin.H{
					"status": 400,
					"error":  "fields are empty",
				}
			}
		} else {
			result = gin.H{
				"status": 404,
				"error":  "instruction not found",
			}
		}
	} else {
		result = gin.H{
			"status": 404,
			"error":  "invalif format",
		}
	}

	c.JSON(http.StatusOK, result)
	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"event_status\": \"83\", \"event_name\": \"100\" }" http://localhost:8080/api/v1/instructions/1
}

func DeleteInstruction(c *gin.Context) {
	id := c.Params.ByName("id")
	var instruction Instruction
	var result gin.H

	if bson.IsObjectIdHex(id) {
		var collection = dbConnect.Use("belajar_golang", CollectionInstruction)
		err := collection.FindId(bson.ObjectIdHex(id)).Select(bson.M{"_id": 1, "event_status": 1, "event_name": 1}).One(&instruction)

		if err == nil {
			var selector = bson.M{"_id": bson.ObjectIdHex(id)}
			err := collection.Remove(selector)

			if err == nil {
				result = gin.H{
					"status": 200,
					"desc":   "id #" + id + " deleted",
				}
			} else {
				result = gin.H{
					"status": 400,
					"error":  "Delete failed",
				}
			}
		} else {
			result = gin.H{
				"status": 404,
				"error":  "instruction not found",
			}
		}
	} else {
		result = gin.H{
			"status": 404,
			"error":  "invalif format",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i -X DELETE http://localhost:8080/api/v1/instructions/1
}
