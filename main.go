package main

import (
	"git-clones/go-api-simple/routes"
)

func main() {
	// repo, err := mysql.NewMySQLRepository("mysql", config.CreateConfig(), 3, 3)
	// if err != nil {
	// 	panic(err)
	// }
	repo := routes.GetMySqlRepo()
	defer repo.Close()
	r := routes.SetUpRouter(repo)
	r.Run()
}
