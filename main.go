package main

import (
	"git-clones/go-api-simple/routes"
)

func main() {
	repo := routes.GetMySqlRepo()
	defer repo.Close()
	r := routes.SetUpRouter(repo)
	r.Run()
}
