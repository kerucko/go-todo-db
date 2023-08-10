package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ilyakaznacheev/cleanenv"
)

var DB *sql.DB

type config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

var c config

func Init() {
	var err error
	readConfig()
	//log.Println(c)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.User, c.Password, c.Host, c.Port, c.Database)
	//log.Println(connection)
	DB, err = sql.Open("mysql", connection)
	if err != nil {
		panic(err)
	}
	//defer DB.Close()

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}

func readConfig() {
	err := cleanenv.ReadConfig("config.yml", &c)
	if err != nil {
		panic(err)
	}
}
