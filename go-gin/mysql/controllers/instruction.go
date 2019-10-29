package controllers

import (
	"rest-api-gin-mysql/db"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Instruction struct {
	Id          int64  `db:"id" json:"id"`
	EventStatus string `db:"event_status" json:"event_status"`
	EventName   string `db:"event_name" json:"event_name"`
}

type InstructionController struct{}

var dbmap = db.InitDb()

func (instruct *InstructionController) GetInstructions(c *gin.Context) {
	var instructions []Instruction
	var result gin.H

	_, err := dbmap.Select(&instructions, "SELECT * FROM instruction")
	if err == nil {
		//c.JSON(200, instructions)

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
		//c.JSON(404, gin.H{"error": "no instruction(s) into the table"})

		result = gin.H{
			"status": 404,
			"error":  "no instruction(s) into the table",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i http://localhost:8080/api/v1/instructions
}

func (instruct *InstructionController) GetInstruction(c *gin.Context) {
	id := c.Params.ByName("id")
	var instruction Instruction
	var result gin.H

	err := dbmap.SelectOne(&instruction, "SELECT * FROM instruction WHERE id=?", id)
	if err == nil {
		instruction_id, _ := strconv.ParseInt(id, 0, 64)

		content := &Instruction{
			Id:          instruction_id,
			EventStatus: instruction.EventStatus,
			EventName:   instruction.EventName,
		}

		//c.JSON(200, content)

		result = gin.H{
			"status": 200,
			"data":   content,
			"count":  1,
		}
	} else {
		//c.JSON(404, gin.H{"error": "instruction not found"})

		result = gin.H{
			"status": 404,
			"error":  "instruction not found",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i http://localhost:8080/api/v1/Instructions/1
}

func (instruct *InstructionController) PostInstruction(c *gin.Context) {
	var instruction Instruction
	var result gin.H
	c.Bind(&instruction)

	if instruction.EventStatus != "" && instruction.EventName != "" {
		if insert, _ := dbmap.Exec(`INSERT INTO instruction (event_status, event_name) VALUES (?, ?)`, instruction.EventStatus, instruction.EventName); insert != nil {
			instruction_id, err := insert.LastInsertId()
			if err == nil {
				content := &Instruction{
					Id:          instruction_id,
					EventStatus: instruction.EventStatus,
					EventName:   instruction.EventName,
				}

				//c.JSON(201, content)

				result = gin.H{
					"status": 200,
					"data":   content,
				}
			} else {
				//checkErr(err, "Insert failed")
				result = gin.H{
					"status": 400,
					"error":  "Insert failed",
				}
			}
		}
	} else {
		//c.JSON(422, gin.H{"error": "fields are empty"})
		result = gin.H{
			"status": 400,
			"error":  "fields are empty",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"event_status\": \"83\", \"event_name\": \"100\" }" http://localhost:8080/api/v1/instructions
}

func (instruct *InstructionController) UpdateInstruction(c *gin.Context) {
	id := c.Params.ByName("id")
	var instruction Instruction
	var result gin.H
	err := dbmap.SelectOne(&instruction, "SELECT * FROM instruction WHERE id=?", id)

	if err == nil {
		var json Instruction
		c.Bind(&json)
		instruction_id, _ := strconv.ParseInt(id, 0, 64)
		instruction := Instruction{
			Id:          instruction_id,
			EventStatus: json.EventStatus,
			EventName:   json.EventName,
		}

		if instruction.EventStatus != "" && instruction.EventName != "" {
			_, err = dbmap.Update(&instruction)

			if err == nil {
				//c.JSON(200, instruction)
				result = gin.H{
					"status": 200,
					"data":   instruction,
				}
			} else {
				//checkErr(err, "Updated failed")
				result = gin.H{
					"status": 400,
					"error":  "Updated failed",
				}
			}
		} else {
			//c.JSON(422, gin.H{"error": "fields are empty"})
			result = gin.H{
				"status": 400,
				"error":  "fields are empty",
			}
		}
	} else {
		//c.JSON(404, gin.H{"error": "instruction not found"})
		result = gin.H{
			"status": 404,
			"error":  "instruction not found",
		}
	}

	c.JSON(http.StatusOK, result)
	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"event_status\": \"83\", \"event_name\": \"100\" }" http://localhost:8080/api/v1/instructions/1
}

func (instruct *InstructionController) DeleteInstruction(c *gin.Context) {
	id := c.Params.ByName("id")
	var instruction Instruction
	var result gin.H

	err := dbmap.SelectOne(&instruction, "SELECT id FROM Instruction WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&instruction)

		if err == nil {
			//c.JSON(200, gin.H{"id #" + id: " deleted"})
			result = gin.H{
				"status": 200,
				"desc":   "id #" + id + " deleted",
			}
		} else {
			//checkErr(err, "Delete failed")
			result = gin.H{
				"status": 400,
				"error":  "Delete failed",
			}
		}
	} else {
		//c.JSON(404, gin.H{"error": "instruction not found"})
		result = gin.H{
			"status": 404,
			"error":  "instruction not found",
		}
	}

	c.JSON(http.StatusOK, result)

	// curl -i -X DELETE http://localhost:8080/api/v1/instructions/1
}
