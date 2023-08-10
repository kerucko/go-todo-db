package handlers

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"todo_db/database"
)

type Task struct {
	ID              int
	Name            string
	Comment         string
	CreateDate      string
	Deadline        string
	AppointmentDate string
}

var (
	TPL *template.Template
	//db  *sql.DB = database.DB
)

func ShowTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("{MAIN}")
	rows, err := database.DB.Query("SELECT * FROM todo.test")
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

	err = TPL.ExecuteTemplate(w, "main_page.html", tasks)
	if err != nil {
		panic(err)
	}
}

func AddNewTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("{ADD GET}")
		TPL.ExecuteTemplate(w, "add_new_task.html", nil)
		return
	}
	log.Println("{ADD POST}")
	r.ParseForm()
	name := r.FormValue("name")
	if name == "" {
		err := TPL.ExecuteTemplate(w, "result.html", "Неправильные введенные данные")
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

	stmt, err := database.DB.Prepare(fmt.Sprintf("INSERT INTO todo.test (name, comment, createDate, deadline, appointmentDate) VALUES (%s, %s, NOW(), %s, %s);", name, comment, deadline, appointmentDate))
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
		TPL.ExecuteTemplate(w, "result.html", "Ошибка")
		return
	}

	err = TPL.ExecuteTemplate(w, "result.html", "Задача добавлена успешно")
	if err != nil {
		panic(err)
	}
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	log.Println("{UPDATE", id, "}")

	row := database.DB.QueryRow("SELECT * FROM todo.test WHERE (id = ?);", id)

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

	err = TPL.ExecuteTemplate(w, "update.html", t)
	if err != nil {
		panic(err)
	}
}

func UpdateResultHandler(w http.ResponseWriter, r *http.Request) {
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

	stmt, err := database.DB.Prepare(fmt.Sprintf("UPDATE todo.test SET name=%s, comment=%s, deadline=%s, appointmentDate=%s WHERE id=%s;", name, comment, deadline, appointmentDate, id))
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
		TPL.ExecuteTemplate(w, "result.html", "Возникла ошибка, попробуйте еще раз")
		return
	}

	err = TPL.ExecuteTemplate(w, "result.html", "Задача успешно обновлена")
	if err != nil {
		panic(err)
	}
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	log.Println("{DELETE", id, "}")

	stmt, err := database.DB.Prepare("DELETE FROM todo.test WHERE (id = ?);")
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

	err = TPL.ExecuteTemplate(w, "result.html", "Задача успешно удалена")
	if err != nil {
		panic(err)
	}
}

func SortHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filter := r.FormValue("sort")
	var stmt string
	if filter == "дедлайну" {
		stmt = "SELECT * FROM todo.test WHERE deadline IS NOT NULL ORDER BY deadline;"
	} else if filter == "дате создания" {
		stmt = "SELECT * FROM todo.test WHERE createDate IS NOT NULL ORDER BY createDate;"
	} else {
		stmt = "SELECT * FROM todo.test WHERE appointmentDate IS NOT NULL ORDER BY appointmentDate;"
	}

	rows, err := database.DB.Query(stmt)
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

	err = TPL.ExecuteTemplate(w, "main_page.html", tasks)
	if err != nil {
		panic(err)
	}
}

func TodayHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("{TODAY}")
	rows, err := database.DB.Query("SELECT * FROM todo.test WHERE appointmentDate=DATE(NOW());")
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

	err = TPL.ExecuteTemplate(w, "main_page.html", tasks)
	if err != nil {
		panic(err)
	}
}
