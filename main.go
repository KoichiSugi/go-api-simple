package main

import (
	"git-clones/go-api-simple/config"
	"git-clones/go-api-simple/mysql"
	"git-clones/go-api-simple/routes"
)

func main() {
	repo, err := mysql.NewMySQLRepository("mysql", config.CreateConfig(), 3, 3)
	if err != nil {
		panic(err)
	}
	defer repo.Close()
	r := routes.SetUpRouter(repo)
	r.Run()
}
