package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	. "rest-api/models"
)

//use "."" if without call model like models.Student_Mongo / models.Response_Mongo

var DB_DATABASE string

func handleError(err error) {
	if err != nil {
		//log.Fatal(err)
		println("Error:", err.Error())
	}
}

func connect() (*mgo.Session, error) {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB_HOST := os.Getenv("DB_HOST_MONGO")
	DB_USER := os.Getenv("DB_USER_MONGO")
	DB_PASSWORD := os.Getenv("DB_PASSWORD_MONGO")

	//var session, err = mgo.Dial("mongodb://golang:golang123@localhost/belajar_golang")
	//var session, err = mgo.Dial("localhost:27017")
	info := &mgo.DialInfo{
		Addrs:     []string{DB_HOST},
		Timeout:   5 * time.Second,
		Database:  DB_DATABASE_MONGO,
		Username:  DB_USER,
		Password:  DB_PASSWORD,
		Mechanism: "SCRAM-SHA-1",
	}
	session, err := mgo.DialWithInfo(info)

	if err != nil {
		panic(err.Error())
	}

	return session, nil
}

func findOne(id string, table string) (Students []Student_Mongo) {
	var session, err = connect()
	handleError(err)
	defer session.Close()

	var result = Student_Mongo{}
	var collection = session.DB(DB_DATABASE_MONGO).C(table)
	//var selector = bson.M{"_id": bson.ObjectIdHex(id)}
	//err = collection.Find(selector).One(&student)
	err = collection.FindId(bson.ObjectIdHex(id)).Select(bson.M{"_id": 1, "name": 1, "grade": 1}).One(&result)
	if err == nil {
		Students = append(Students, result)
	}

	return Students
}

func find(table string, sortfield string, sort string, limit int, skip int) (Students []Student_Mongo) {
	var session, err = connect()
	handleError(err)
	defer session.Close()

	var collection = session.DB(DB_DATABASE_MONGO).C(table)
	var result []Student_Mongo
	var selector = bson.M{}
	//err = collection.Find(selector).All(&result)
	err = collection.Find(selector).Sort(sort + sortfield).Skip(skip).Limit(limit).All(&result)
	handleError(err)

	for _, each := range result {
		//fmt.Printf("Name : %s\t Grade : %d\n", each.Name, each.Grade)
		Students = append(Students, each)
	}

	return Students
}

func update(id string, table string, name string, grade int) (stat bool) {
	var session, err = connect()
	handleError(err)
	defer session.Close()

	var collection = session.DB(DB_DATABASE_MONGO).C(table)
	var selector = bson.M{"_id": bson.ObjectIdHex(id)}
	var changes = Student_Mongo{"", name, grade}
	//var changes = bson.M{"$set": bson.M{"name": "Poltak", "grade": 2}}
	err = collection.Update(selector, changes)
	handleError(err)

	stat = true
	return stat
}

func delete(id string, table string) (stat bool) {
	var session, err = connect()
	handleError(err)
	defer session.Close()

	var collection = session.DB(DB_DATABASE_MONGO).C(table)
	var selector = bson.M{"_id": bson.ObjectIdHex(id)}
	err = collection.Remove(selector)
	handleError(err)

	stat = true
	return stat
}

func insert(table string, name string, grade int) (ID bson.ObjectId) {
	var session, err = connect()
	handleError(err)
	defer session.Close()

	var collection = session.DB(DB_DATABASE_MONGO).C(table)
	var student Student_Mongo
	student.Id = bson.NewObjectId()
	err = collection.Insert(&Student_Mongo{student.Id, name, grade})
	handleError(err)

	return student.Id
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Rest API Student (Mongo)")
}

func getAllStudent(w http.ResponseWriter, r *http.Request) {
	var response Response_Mongo

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = 5
	}

	skip, err := strconv.Atoi(r.FormValue("skip"))
	if err != nil {
		skip = 0
	}

	student := find("student", "_id", "-", limit, skip)

	response.Status = 1
	response.Message = "success"
	response.Data = student

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	var response Response_Mongo

	var key = r.FormValue("id")
	if len(key) == 0 {
		response.Status = 0
		response.Message = "field id is required"
	} else {
		student := findOne(key, "student")
		response.Status = 1
		response.Message = "success"
		response.Data = student
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func insertStudent(w http.ResponseWriter, r *http.Request) {
	var response Response_Mongo
	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	name := r.FormValue("name")
	grade, err := strconv.Atoi(r.FormValue("grade"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resIns := insert("student", name, grade)

	student := []Student_Mongo{
		Student_Mongo{Id: resIns, Name: name, Grade: grade},
	}

	response.Status = 1
	response.Message = "success"
	response.Data = student

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	var response Response_Mongo
	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	var id = r.FormValue("id")
	if len(id) == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	name := r.FormValue("name")
	grade, err := strconv.Atoi(r.FormValue("grade"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	checkstudent := findOne(id, "student")

	if len(checkstudent) != 0 {
		resUpd := update(id, "student", name, grade)
		var student = []Student_Mongo{}

		if resUpd {
			response.Status = 1
			response.Message = "success"

			student = []Student_Mongo{
				Student_Mongo{Id: bson.ObjectIdHex(id), Name: name, Grade: grade},
			}

			response.Data = student
		} else {
			response.Status = 0
			response.Message = "failed update student"
			response.Data = []Student_Mongo{}
		}
	} else {
		response.Status = 0
		response.Message = "not found"
		response.Data = []Student_Mongo{}
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	var response Response_Mongo
	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	var id = r.FormValue("id")
	if len(id) == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	checkstudent := findOne(id, "student")

	if len(checkstudent) != 0 {

		resDel := delete(id, "student")

		if resDel {
			response.Status = 1
			response.Message = "success"
			response.Data = []Student_Mongo{}
		} else {
			response.Status = 0
			response.Message = "failed delete student"
			response.Data = []Student_Mongo{}
		}
	} else {
		response.Status = 0
		response.Message = "not found"
		response.Data = []Student_Mongo{}
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/students", getAllStudent).Methods("POST")
	myRouter.HandleFunc("/student/get", getStudent).Methods("POST")
	myRouter.HandleFunc("/student/add", insertStudent).Methods("POST")
	myRouter.HandleFunc("/student/update", updateStudent).Methods("POST")
	myRouter.HandleFunc("/student/delete", deleteStudent).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))

	/*student := findOne("5d2415bdae76b390c4206f42", "student")
	fmt.Println(student)*/

	/*allstudent := find("student", "_id", "-", 0, 2)
	fmt.Println(allstudent)*/

	//resUpd := update("5d2454c5ae76b390c42077b6", "student", "Poltak Tulus", 2)
	//fmt.Println(resUpd)

	//resDel := delete("5d2454c5ae76b390c42077b6", "student")
	//fmt.Println(resDel)

	//resIns := insert("student", "Evan Dimas", 5)
	//fmt.Println(resIns)
}
