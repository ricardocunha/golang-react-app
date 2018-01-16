package common

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//SESSION is an DB field for use like session store
var SESSION *gorm.DB

//InitSession Opening a database and save the reference to `Database` struct.
func InitSession() *gorm.DB {
	SESSION, err := gorm.Open("sqlite3", "./../sesssion.db")
	if err != nil {
		fmt.Println("db error: ", err)
	}
	SESSION.DB().SetMaxIdleConns(2)
	return SESSION
}

//GetSession Using this function to get a connection, you can create your connection pool here.
func GetSession() *gorm.DB {
	return DB
}
