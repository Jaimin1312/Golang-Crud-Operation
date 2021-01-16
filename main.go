package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id   int
	Name string
	City string
}

var tmpl = template.Must(template.ParseGlob("template/*"))

func dbConn() *sql.DB {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "EmployeeDatabase"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database is connected")
	}
	return db
}

func createform(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "create.html", nil)
}

func insertdata(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city)
		log.Println("INSERT: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Employee ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}
	file, _ := json.MarshalIndent(res, "", " ")
	_ = ioutil.WriteFile("Employee.json", file, 0644)

	tmpl.ExecuteTemplate(w, "index.html", res)
	defer db.Close()
}

func delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	empid := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", empid)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err = selDB.Scan(&id, &name, &city)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	tmpl.ExecuteTemplate(w, "Edit.html", emp)
	defer db.Close()
}

func update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("empid")
		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city, id)
		log.Println("UPDATE: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	fmt.Println("Server started at 9000")
	http.HandleFunc("/", index)
	http.HandleFunc("/submitdata", insertdata)
	http.HandleFunc("/create", createform)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/update", update)
	http.ListenAndServe(":9000", nil)

}
