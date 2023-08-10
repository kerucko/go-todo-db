package handlers

import (
	"database/sql"
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
	FinishDate      string
}

var TPL *template.Template

func ShowTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("{MAIN}")
	rows, err := database.DB.Query("SELECT * FROM todo.tasks")
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

	stmt, err := database.DB.Prepare(fmt.Sprintf("INSERT INTO todo.tasks (name, comment, createDate, deadline, appointmentDate) VALUES (%s, %s, NOW(), %s, %s);", name, comment, deadline, appointmentDate))
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
	from := r.FormValue("from")
	log.Println("{UPDATE", id, "}")

	row := database.DB.QueryRow(fmt.Sprintf("SELECT * FROM todo.%s WHERE (id = %s);", from, id))

	var t Task
	var createDate, deadline, appointmentDate, finishDate []uint8
	var err error
	if from == "tasks" {
		err = row.Scan(&t.ID, &t.Name, &t.Comment, &createDate, &deadline, &appointmentDate)
	} else if from == "completed" {
		err = row.Scan(&t.ID, &t.Name, &t.Comment, &createDate, &deadline, &appointmentDate, &finishDate)
	}
	if err != nil {
		panic(err)
	}
	t.CreateDate = string(createDate)
	t.Deadline = string(deadline)
	t.AppointmentDate = string(appointmentDate)
	t.FinishDate = string(finishDate)
	log.Println(t)

	err = TPL.ExecuteTemplate(w, "update.html", struct {
		Data Task
		From string
	}{t, from})
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

	table := r.FormValue("from")

	var stmt *sql.Stmt
	var err error
	if table == "tasks" {
		stmt, err = database.DB.Prepare(fmt.Sprintf("UPDATE %s SET name=%s, comment=%s, deadline=%s, appointmentDate=%s WHERE id=%s;", table, name, comment, deadline, appointmentDate, id))
	} else if table == "completed" {
		finishDate := r.FormValue("finishDate")
		if finishDate == "" {
			finishDate = "NULL"
		} else {
			finishDate = "'" + finishDate + "'"
		}

		stmt, err = database.DB.Prepare(fmt.Sprintf("UPDATE %s SET name=%s, comment=%s, deadline=%s, appointmentDate=%s, finishDate=%s WHERE id=%s;", table, name, comment, deadline, appointmentDate, finishDate, id))
	}
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
	from := r.FormValue("from")
	log.Println("{DELETE", id, "}")

	var table string
	if from == "1" {
		table = "todo.tasks"
	} else if from == "2" {
		table = "todo.completed"
	} else {
		panic("error table")
	}

	stmt, err := database.DB.Prepare(fmt.Sprintf("DELETE FROM %s WHERE (id = %s);", table, id))
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
		stmt = "SELECT * FROM todo.tasks WHERE deadline IS NOT NULL ORDER BY deadline;"
	} else if filter == "дате создания" {
		stmt = "SELECT * FROM todo.tasks WHERE createDate IS NOT NULL ORDER BY createDate;"
	} else {
		stmt = "SELECT * FROM todo.tasks WHERE appointmentDate IS NOT NULL ORDER BY appointmentDate;"
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
	rows, err := database.DB.Query("SELECT * FROM todo.tasks WHERE appointmentDate=DATE(NOW());")
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

func DoneHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	id := r.FormValue("id")

	row := database.DB.QueryRow("SELECT * FROM todo.tasks WHERE id=?;", id)
	var t Task
	var createDate, deadline, appointmentDate []uint8
	err = row.Scan(&t.ID, &t.Name, &t.Comment, &createDate, &deadline, &appointmentDate)
	if err != nil {
		panic(err)
	}
	t.CreateDate = string(createDate)
	t.Deadline = string(deadline)
	t.AppointmentDate = string(appointmentDate)
	if t.Deadline == "" {
		t.Deadline = "NULL"
	} else {
		t.Deadline = "'" + t.Deadline + "'"
	}
	if t.AppointmentDate == "" {
		t.AppointmentDate = "NULL"
	} else {
		t.AppointmentDate = "'" + t.AppointmentDate + "'"
	}
	log.Println(t)

	stmtInsert, err := database.DB.Prepare(fmt.Sprintf("INSERT INTO todo.completed (id, name, comment, createDate, deadline, appointmentDate, finishDate) VALUES (%s, '%s', '%s', '%s', %s, %s, NOW());", id, t.Name, t.Comment, t.CreateDate, t.Deadline, t.AppointmentDate))
	if err != nil {
		log.Println("stmt error")
		panic(err)
	}
	defer stmtInsert.Close()

	res, err := stmtInsert.Exec()
	if err != nil {
		log.Println("error insert in completed: ", err)
		panic(err)
	}
	rowsAf, _ := res.RowsAffected()
	if err != nil || rowsAf != 1 {
		log.Println("Error insert in completed:", err)
		err := TPL.ExecuteTemplate(w, "result.html", "Ошибка")
		if err != nil {
			panic(err)
		}
		return
	}

	stmtDelete, err := database.DB.Prepare("DELETE FROM todo.tasks WHERE (id = ?);")
	if err != nil {
		panic(err)
	}
	defer stmtDelete.Close()
	res, err = stmtDelete.Exec(id)
	if err != nil {
		panic(err)
	}
	rowsAf, _ = res.RowsAffected()
	if rowsAf != 1 {
		log.Println("ERROR: ", rowsAf)
	}

	err = TPL.ExecuteTemplate(w, "result.html", "Задача выполнена и перенесена в журнал")
	if err != nil {
		panic(err)
	}
}

func ShowCompletedHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("{COMPLETED}")
	rows, err := database.DB.Query("SELECT * FROM todo.completed")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var createDate, deadline, appointmentDate, finishDate []uint8
		err := rows.Scan(&t.ID, &t.Name, &t.Comment, &createDate, &deadline, &appointmentDate, &finishDate)
		if err != nil {
			panic(err)
		}
		t.CreateDate = string(createDate)
		t.Deadline = string(deadline)
		t.AppointmentDate = string(appointmentDate)
		t.FinishDate = string(finishDate)
		tasks = append(tasks, t)
	}

	err = TPL.ExecuteTemplate(w, "completed.html", tasks)
	if err != nil {
		panic(err)
	}
}

func UndoHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	id := r.FormValue("id")

	row := database.DB.QueryRow("SELECT * FROM todo.completed WHERE id=?;", id)
	var t Task
	var createDate, deadline, appointmentDate, finishDate []uint8
	err = row.Scan(&t.ID, &t.Name, &t.Comment, &createDate, &deadline, &appointmentDate, &finishDate)
	if err != nil {
		panic(err)
	}
	t.CreateDate = string(createDate)
	t.Deadline = string(deadline)
	t.AppointmentDate = string(appointmentDate)
	if t.Deadline == "" {
		t.Deadline = "NULL"
	} else {
		t.Deadline = "'" + t.Deadline + "'"
	}
	if t.AppointmentDate == "" {
		t.AppointmentDate = "NULL"
	} else {
		t.AppointmentDate = "'" + t.AppointmentDate + "'"
	}
	log.Println(t)

	stmtInsert, err := database.DB.Prepare(fmt.Sprintf("INSERT INTO todo.tasks (id, name, comment, createDate, deadline, appointmentDate) VALUES (%s, '%s', '%s', '%s', %s, %s);", id, t.Name, t.Comment, t.CreateDate, t.Deadline, t.AppointmentDate))
	if err != nil {
		log.Println("stmt error")
		panic(err)
	}
	defer stmtInsert.Close()

	res, err := stmtInsert.Exec()
	if err != nil {
		log.Println("error insert in completed: ", err)
		panic(err)
	}
	rowsAf, _ := res.RowsAffected()
	if err != nil || rowsAf != 1 {
		log.Println("Error insert in completed:", err)
		err := TPL.ExecuteTemplate(w, "result.html", "Ошибка")
		if err != nil {
			panic(err)
		}
		return
	}

	stmtDelete, err := database.DB.Prepare("DELETE FROM todo.completed WHERE (id = ?);")
	if err != nil {
		panic(err)
	}
	defer stmtDelete.Close()
	res, err = stmtDelete.Exec(id)
	if err != nil {
		panic(err)
	}
	rowsAf, _ = res.RowsAffected()
	if rowsAf != 1 {
		log.Println("ERROR: ", rowsAf)
	}

	err = TPL.ExecuteTemplate(w, "result.html", "Задача возвращена")
	if err != nil {
		panic(err)
	}
}
