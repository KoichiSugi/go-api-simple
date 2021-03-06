package config

import (
	"github.com/go-sql-driver/mysql"
)

func CreateConfig() mysql.Config {
	config := mysql.Config{
		User:   "tester",
		Passwd: "secret",
		//Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "IGD",
		ParseTime: true, /// Parse time values to time.Time
	}
	return config
}
