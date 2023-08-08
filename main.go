package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
)

var (
	tpl *template.Template
	db  *sql.DB
)

type Task struct {
	ID              int
	Name            string
	Comment         string
	CreateDate      string
	Deadline        string
	AppointmentDate string
}

func main() {
	tpl, _ = template.ParseGlob("templates/*.html")

	var err error
	db, err = sql.Open("mysql", "root:mysql_password1@tcp(localhost:3306)/todo")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", startHandler)
	http.HandleFunc("/all_tasks", showTasksHandler)
	http.HandleFunc("/add_new_task", addNewTaskHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
	log.Println("run server")
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("{START}")
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "start.html", nil)
		return
	}
	http.Redirect(w, r, "/all_tasks", http.StatusFound)
}

func showTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("{MAIN POST}")
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var createDate, deadline, appointmentDate []uint8
		err := rows.Scan(&t.ID, &t.Name, &t.Comment, &createDate, &deadline, &appointmentDate)
		if err != nil {
			panic(err)
		}
		t.CreateDate = string(createDate)
		t.Deadline = string(deadline)
		t.AppointmentDate = string(appointmentDate)
		tasks = append(tasks, t)
	}

	tpl.ExecuteTemplate(w, "main_page.html", tasks)
}

func addNewTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("add new task GET")
		tpl.ExecuteTemplate(w, "add_new_task.html", nil)
		return
	}
	log.Println("add new task POST")
	r.ParseForm()
	name := r.FormValue("name")
	comment := r.FormValue("comment")
	deadline := r.FormValue("deadline")
	appointmentDate := r.FormValue("appointmentDate")
	log.Println(name, comment, deadline, appointmentDate)

	stmt, err := db.Prepare("INSERT INTO test (name, comment, createDate, deadline, appointmentDate) VALUES (?, ?, NOW(), ?, ?);")
	if err != nil {
		log.Println("stmt error")
		panic(err)
	}
	defer stmt.Close()
	log.Println("1")

	res, err := stmt.Exec(name, comment, deadline, appointmentDate)
	if err != nil {
		log.Println("error insert: ", err)
		panic(err)
	}
	log.Println("2")
	rowsAf, _ := res.RowsAffected()
	if err != nil || rowsAf != 1 {
		log.Println("Error insert:", err)
		tpl.ExecuteTemplate(w, "add_new_task.html", "ERROR")
		return
	}

	tpl.ExecuteTemplate(w, "add_new_task.html", "Success")
}
