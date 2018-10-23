package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username  string
	Password  string
	Email     string
	Telephone string
	Login     bool
}

var DBSql *sql.DB

func init() {
	var err error
	DBSql, err = sql.Open("mysql", "root:7home7Tmade@tcp(127.0.0.1:3306)/agenda?charset=utf8")
	if err != nil {
		panic(err)
	}
}
