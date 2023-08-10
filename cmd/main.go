package main

import (
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"todo_db/database"
	"todo_db/handlers"
)

func main() {
	database.Init()
	defer database.DB.Close()

	handlers.TPL, _ = template.ParseGlob("templates/*.html")

	http.HandleFunc("/show", handlers.ShowTasksHandler)
	http.HandleFunc("/add", handlers.AddNewTaskHandler)
	http.HandleFunc("/update/", handlers.UpdateTaskHandler)
	http.HandleFunc("/update_result/", handlers.UpdateResultHandler)
	http.HandleFunc("/delete/", handlers.DeleteTaskHandler)
	http.HandleFunc("/sort", handlers.SortHandler)
	http.HandleFunc("/today", handlers.TodayHandler)
	http.HandleFunc("/done/", handlers.DoneHandler)
	http.HandleFunc("/show_completed/", handlers.ShowCompletedHandler)
	http.HandleFunc("/undo/", handlers.UndoHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
