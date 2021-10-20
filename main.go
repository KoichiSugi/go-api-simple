package main

import (
	"git-clones/go-api-simple/config"
	"git-clones/go-api-simple/mysql"
	"git-clones/go-api-simple/routes"
	"log"
)

func main() {
	// config.ConnectDb()
	// defer config.Db.Close()
	//exp
	repo, err := mysql.NewMySQLRepository("mysql", config.CreateConfig(), 3, 3)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	r := routes.SetUpRouter(repo)
	r.Run()
	//exp above
	// r := routes.CreateRouter()
	// r.Run()
	// http.HandleFunc("/", service.GetAllEmployees)
	// http.HandleFunc("/test", logging(service.GetAllEmployees))
	// fmt.Println("listening on port 8080")
	// http.ListenAndServe(":8080", nil)

}
