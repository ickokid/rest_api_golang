package main

import "fmt"
import "log"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "net/http"
import "encoding/json"
import "github.com/gorilla/mux"
import "strconv"

import . "rest-api/models" //use "."" if without call model like models.Student / models.Response

/*type Student struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

type Response struct {
	Status  uint8     `json:"status"`
	Message string    `json:"message"`
	Data    []Student `json:"data"`
}*/

var Students []Student

func handleError(err error) {
	if err != nil {
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred
			//fmt.Println("No rows were returned!")
		} else {
			//log.Fatal(err)
			println("Error:", err.Error())
		}
	}
}

func connect() (*sql.DB, error) {
	//db, err := sql.Open("mysql", "root:root@tcp(localhost:3307)/db_belajar_golang")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", DB_USER, DB_PASSWORD, DB_HOST, DB_DATABASE))
	if err != nil {
		panic(err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	return db, nil
}

func insert(table string, name string, grade int) (ID int64) {
	db, err := connect()
	handleError(err)
	defer db.Close()

	stmt, err := db.Prepare("INSERT " + table + " SET name=?,grade=?")
	handleError(err)

	res, err := stmt.Exec(name, grade)
	handleError(err)

	ID, err = res.LastInsertId()
	handleError(err)

	return ID
}

func update(id int64, table string, name string, grade int) (stat bool) {
	db, err := connect()
	handleError(err)
	defer db.Close()

	_, err = db.Exec("UPDATE "+table+" SET name = ?, grade = ? WHERE id = ?", name, grade, id)
	handleError(err)

	stat = true
	return stat
}

func delete(id int64, table string) (stat bool) {
	db, err := connect()
	handleError(err)
	defer db.Close()

	_, err = db.Exec("DELETE FROM "+table+" WHERE id = ?", id)
	handleError(err)

	stat = true
	return stat
}

func findOne(id int64, table string) (Students []Student) {
	db, err := connect()
	handleError(err)
	defer db.Close()

	var student = Student{}

	err = db.QueryRow("SELECT id, name, grade FROM "+table+" WHERE id = ?", id).Scan(&student.Id, &student.Name, &student.Grade)
	handleError(err)
	if err == nil {
		Students = append(Students, student)
	}

	return Students
}

func find(table string, sortfield string, sort string, limit int, offset int) (Students []Student) {
	db, err := connect()
	handleError(err)
	defer db.Close()

	rows, err := db.Query("SELECT id, name, grade FROM "+table+" ORDER BY "+sortfield+" "+sort+" limit ? offset ?", limit, offset)
	handleError(err)
	defer rows.Close()

	for rows.Next() {
		each := Student{}
		rows.Scan(&each.Id, &each.Name, &each.Grade)
		Students = append(Students, each)
	}

	return Students
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Rest API Student (MySQL)")
}

func getAllStudent(w http.ResponseWriter, r *http.Request) {
	var response Response

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		/*http.Error(w, err.Error(), http.StatusInternalServerError)
		return*/
		limit = 5
	}

	skip, err := strconv.Atoi(r.FormValue("skip"))
	if err != nil {
		/*http.Error(w, err.Error(), http.StatusInternalServerError)
		return*/
		skip = 0
	}

	student := find("student", "id", "DESC", limit, skip)

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
	var response Response

	/*vars := mux.Vars(r)
	key, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/
	//fmt.Fprintf(w, "Key: "+key)

	key, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return

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
	var response Response
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

	student := []Student{
		Student{Id: resIns, Name: name, Grade: grade},
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
	var response Response
	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
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
		var student = []Student{}

		if resUpd {
			response.Status = 1
			response.Message = "success"

			student = []Student{
				Student{Id: id, Name: name, Grade: grade},
			}

			response.Data = student
		} else {
			response.Status = 0
			response.Message = "failed update student"
			response.Data = []Student{}
		}
	} else {
		response.Status = 0
		response.Message = "not found"
		response.Data = []Student{}
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
	var response Response
	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	checkstudent := findOne(id, "student")

	if len(checkstudent) != 0 {

		resDel := delete(id, "student")

		if resDel {
			response.Status = 1
			response.Message = "success"
			response.Data = []Student{}
		} else {
			response.Status = 0
			response.Message = "failed delete student"
			response.Data = []Student{}
		}
	} else {
		response.Status = 0
		response.Message = "not found"
		response.Data = []Student{}
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

	//resIns := insert("student", "Witan Sulaiman", 1)
	//fmt.Println(resIns)

	//resDel := delete(15, "student")
	//fmt.Println(resDel)

	//resUpd := update(4, "student", "Rutfianah Theodora", 1)
	//fmt.Println(resUpd)

	//student := findOne(3, "student")
	//fmt.Println(student)

	//students := find("student", "id", "DESC", 5, 0)
	//fmt.Println(students)

	/*http.HandleFunc("/", homePage)
	http.HandleFunc("/students", allStudents)
	log.Fatal(http.ListenAndServe(":10000", nil))*/
}
