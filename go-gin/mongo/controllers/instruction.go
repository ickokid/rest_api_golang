package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"rest-api-gin-mongo/models"
)

func GetInstructions(c *gin.Context) {
	var result gin.H

	instructions, err := models.FindAll(0, 5)

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
			"error":  "No instruction(s) into the table",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i http://localhost:8080/api/v1/instructions
}

func GetInstruction(c *gin.Context) {
	id := c.Param("id")
	var result gin.H

	if bson.IsObjectIdHex(id) {
		instructions, err := models.FindById(id)

		if err == nil {
			result = gin.H{
				"status": 200,
				"data":   instructions,
				"count":  1,
			}
		} else {
			result = gin.H{
				"status": 404,
				"error":  "Instruction not found",
			}
		}
	} else {
		result = gin.H{
			"status": 404,
			"error":  "Invalid format",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i http://localhost:8080/api/v1/Instructions/1
}

func PostInstruction(c *gin.Context) {
	var instruction models.Instruction
	var result gin.H
	errs := c.Bind(&instruction)

	if errs == nil {
		instruction_id := bson.NewObjectId()

		instructions, err := models.Insert(instruction_id, instruction.EventStatus, instruction.EventName)

		if err == nil {
			result = gin.H{
				"status": 200,
				"data":   instructions,
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
			"error":  errs.Error(),
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"event_status\": \"83\", \"event_name\": \"100\" }" http://localhost:8080/api/v1/instructions
}

func UpdateInstruction(c *gin.Context) {
	id := c.Params.ByName("id")
	var result gin.H

	if bson.IsObjectIdHex(id) {
		_, err := models.FindById(id)

		if err == nil {
			var instruction models.Instruction

			errs := c.Bind(&instruction)

			if errs == nil {
				instructions, err := models.Update(id, instruction.EventStatus, instruction.EventName)

				if err == nil {
					result = gin.H{
						"status": 200,
						"data":   instructions,
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
					"error":  "Fields are empty",
				}
			}
		} else {
			result = gin.H{
				"status": 404,
				"error":  "Instruction not found",
			}
		}
	} else {
		result = gin.H{
			"status": 404,
			"error":  "Invalid format",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"event_status\": \"83\", \"event_name\": \"100\" }" http://localhost:8080/api/v1/instructions/1
}

func DeleteInstruction(c *gin.Context) {
	id := c.Params.ByName("id")
	var result gin.H

	if bson.IsObjectIdHex(id) {
		_, err := models.FindById(id)

		if err == nil {
			err := models.Delete(id)

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
				"error":  "Instruction not found",
			}
		}
	} else {
		result = gin.H{
			"status": 404,
			"error":  "Invalid format",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i -X DELETE http://localhost:8080/api/v1/instructions/1
}
