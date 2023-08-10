package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", "root:mysql_password1@tcp(localhost:3306)/todo")
	if err != nil {
		panic(err)
	}
	//defer DB.Close()

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}
