package controllers

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"rest-api-gin-mongo/models"
	"rest-api-gin-mongo/utils"
)

func GetInstructions(c *gin.Context) {
	instructions, err := models.FindAll(0, 5)

	if err == nil {
		if len(instructions) <= 0 {
			utils.Success(c, nil)
		} else {
			utils.Success(c, instructions)
		}
	} else {
		utils.NotFound(c)
	}

	// curl -i http://localhost:8080/api/v1/instructions
}

func GetInstruction(c *gin.Context) {
	id := c.Param("id")

	if bson.IsObjectIdHex(id) {
		instructions, err := models.FindById(id)

		if err == nil {
			utils.Success(c, instructions)
		} else {
			utils.NotFound(c)
		}
	} else {
		utils.BadRequest(c, "Invalid format")
	}

	// curl -i http://localhost:8080/api/v1/Instructions/1
}

func PostInstruction(c *gin.Context) {
	var instruction models.Instruction
	errs := c.Bind(&instruction)

	if errs == nil {
		instruction_id := bson.NewObjectId()

		instructions, err := models.Insert(instruction_id, instruction.EventStatus, instruction.EventName)

		if err == nil {
			utils.Success(c, instructions)
		} else {
			utils.InternalError(c, "Insert failed", err)
		}
	} else {
		utils.BadRequest(c, "Parameters in complete") //errs.Error()
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"event_status\": \"83\", \"event_name\": \"100\" }" http://localhost:8080/api/v1/instructions
}

func UpdateInstruction(c *gin.Context) {
	id := c.Params.ByName("id")

	if bson.IsObjectIdHex(id) {
		_, err := models.FindById(id)

		if err == nil {
			var instruction models.Instruction

			errs := c.Bind(&instruction)

			if errs == nil {
				instructions, err := models.Update(id, instruction.EventStatus, instruction.EventName)

				if err == nil {
					utils.Success(c, instructions)
				} else {
					utils.InternalError(c, "Updated failed", err)
				}
			} else {
				utils.BadRequest(c, "Parameters in complete")
			}
		} else {
			utils.NotFound(c)
		}
	} else {
		utils.BadRequest(c, "Invalid format")
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"event_status\": \"83\", \"event_name\": \"100\" }" http://localhost:8080/api/v1/instructions/1
}

func DeleteInstruction(c *gin.Context) {
	id := c.Params.ByName("id")

	if bson.IsObjectIdHex(id) {
		_, err := models.FindById(id)

		if err == nil {
			err := models.Delete(id)

			if err == nil {
				utils.Success(c, nil)
			} else {
				utils.InternalError(c, "Delete failed", err)
			}
		} else {
			utils.NotFound(c)
		}
	} else {
		utils.BadRequest(c, "Invalid format")
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/instructions/1
}
