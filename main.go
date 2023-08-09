package main

import (
	"database/sql"
	"fmt"
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

	//http.HandleFunc("/", startHandler)
	http.HandleFunc("/show", showTasksHandler)
	http.HandleFunc("/add", addNewTaskHandler)
	http.HandleFunc("/update/", updateTaskHandler)
	http.HandleFunc("/update_result/", updateResultHandler)
	http.HandleFunc("/delete/", deleteTaskHandler)
	http.HandleFunc("/sort", SortHandler)
	http.HandleFunc("/today", TodayHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func showTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("{MAIN}")
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

	err = tpl.ExecuteTemplate(w, "main_page.html", tasks)
	if err != nil {
		panic(err)
	}
}

func addNewTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("{ADD GET}")
		tpl.ExecuteTemplate(w, "add_new_task.html", nil)
		return
	}
	log.Println("{ADD POST}")
	r.ParseForm()
	name := r.FormValue("name")
	if name == "" {
		err := tpl.ExecuteTemplate(w, "result.html", "Неправильные введенные данные")
		if err != nil {
			panic(err)
		}
		return
	}
	name = "'" + name + "'"
	comment := r.FormValue("comment")
	comment = "'" + comment + "'"
	deadline := r.FormValue("deadline")
	if deadline == "" {
		deadline = "NULL"
	} else {
		deadline = "'" + deadline + "'"
	}
	appointmentDate := r.FormValue("appointmentDate")
	if appointmentDate == "" {
		appointmentDate = "NULL"
	} else {
		appointmentDate = "'" + appointmentDate + "'"
	}
	log.Println(name, comment, deadline, appointmentDate)

	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO test (name, comment, createDate, deadline, appointmentDate) VALUES (%s, %s, NOW(), %s, %s);", name, comment, deadline, appointmentDate))
	if err != nil {
		log.Println("stmt error")
		panic(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		log.Println("error insert: ", err)
		panic(err)
	}
	rowsAf, _ := res.RowsAffected()
	if err != nil || rowsAf != 1 {
		log.Println("Error insert:", err)
		tpl.ExecuteTemplate(w, "result.html", "Ошибка")
		return
	}

	err = tpl.ExecuteTemplate(w, "result.html", "Задача добавлена успешно")
	if err != nil {
		panic(err)
	}
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	log.Println("{UPDATE", id, "}")

	row := db.QueryRow("SELECT * FROM test WHERE (id = ?);", id)

	var t Task
	var createDate, deadline, appointmentDate []uint8
	err := row.Scan(&t.ID, &t.Name, &t.Comment, &createDate, &deadline, &appointmentDate)
	if err != nil {
		panic(err)
	}
	t.CreateDate = string(createDate)
	t.Deadline = string(deadline)
	t.AppointmentDate = string(appointmentDate)
	log.Println(t)

	err = tpl.ExecuteTemplate(w, "update.html", t)
	if err != nil {
		panic(err)
	}
}

func updateResultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("{UPDATE RESULT}")
	r.ParseForm()
	id := r.FormValue("id")
	name := r.FormValue("name")
	name = "'" + name + "'"
	comment := r.FormValue("comment")
	comment = "'" + comment + "'"
	deadline := r.FormValue("deadline")
	if deadline == "" {
		deadline = "NULL"
	} else {
		deadline = "'" + deadline + "'"
	}
	appointmentDate := r.FormValue("appointmentDate")
	if appointmentDate == "" {
		appointmentDate = "NULL"
	} else {
		appointmentDate = "'" + appointmentDate + "'"
	}
	log.Println(id, name, comment, deadline, appointmentDate)

	stmt, err := db.Prepare(fmt.Sprintf("UPDATE test SET name=%s, comment=%s, deadline=%s, appointmentDate=%s WHERE id=%s;", name, comment, deadline, appointmentDate, id))
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		panic(err)
	}
	rowsAf, _ := res.RowsAffected()
	if rowsAf != 1 {
		log.Println("Error: ", err)
		tpl.ExecuteTemplate(w, "result.html", "Возникла ошибка, попробуйте еще раз")
		return
	}

	err = tpl.ExecuteTemplate(w, "result.html", "Задача успешно обновлена")
	if err != nil {
		panic(err)
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	log.Println("{DELETE", id, "}")

	stmt, err := db.Prepare("DELETE FROM test WHERE (id = ?);")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}
	rowsAf, _ := res.RowsAffected()
	if rowsAf != 1 {
		log.Println("ERROR: ", rowsAf)
	}

	err = tpl.ExecuteTemplate(w, "result.html", "Задача успешно удалена")
	if err != nil {
		panic(err)
	}
}

func SortHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filter := r.FormValue("sort")
	var stmt string
	if filter == "дедлайну" {
		stmt = "SELECT * FROM test WHERE deadline IS NOT NULL ORDER BY deadline;"
	} else if filter == "дате создания" {
		stmt = "SELECT * FROM test WHERE createDate IS NOT NULL ORDER BY createDate;"
	} else {
		stmt = "SELECT * FROM test WHERE appointmentDate IS NOT NULL ORDER BY appointmentDate;"
	}

	rows, err := db.Query(stmt)
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

	err = tpl.ExecuteTemplate(w, "main_page.html", tasks)
	if err != nil {
		panic(err)
	}
}

func TodayHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("{TODAY}")
	rows, err := db.Query("SELECT * FROM test WHERE appointmentDate=DATE(NOW());")
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

	err = tpl.ExecuteTemplate(w, "main_page.html", tasks)
	if err != nil {
		panic(err)
	}
}
