package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func ConnectDb() {
	config := mysql.Config{
		User:      "root",
		Passwd:    "pass",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "IGD",
		ParseTime: true, /// Parse time values to time.Time
	}
	var err error
	Db, err = sql.Open("mysql", config.FormatDSN()) //convert confing into DSN format for a connection
	if err != nil {
		log.Fatal("err")
	}
	pingErr := Db.Ping() //Ping verifies a connection to the database is still alive
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Database Connected")
}

func CreateConfig() mysql.Config {
	config := mysql.Config{
		User:      "root",
		Passwd:    "pass",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "IGD",
		ParseTime: true, /// Parse time values to time.Time
	}
	return config
}
