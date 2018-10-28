package models

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type User struct {
	Username  string
	Password  string
	Email     string
	Telephone string
	Login     bool
}

type Meeting struct {
	Title        string `xorm:"unique"`
	Moderator    string
	Participants string
	StartTime    time.Time
	EndTime      time.Time
}

type UserMeeting struct {
	Participator string
	Title        string
	Originate    bool
}

var DBSql *sql.DB
var Logger *log.Logger
var logFile *os.File
var Engine *xorm.Engine

func init() {
	var err error
	DBSql, err = sql.Open("mysql", "root:7home7Tmade@tcp(127.0.0.1:3306)/agenda?charset=utf8")
	if err != nil {
		panic(err)
	}
	logFile, err = os.OpenFile("github.com/liuyh73/Go/agenda/log/logFile.txt", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	Logger = log.New(logFile, "", log.LstdFlags)
	Engine, err = xorm.NewEngine("mysql", "root:7home7Tmade@tcp(127.0.0.1:3306)/agenda?charset=utf8")
	if err != nil {
		panic(err)
	}
}
