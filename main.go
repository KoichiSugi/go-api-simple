package main

import (
	"git-clones/go-api-simple/config"
	"git-clones/go-api-simple/routes"
)

func main() {
	config.ConnectDb()
	defer config.Db.Close()
	r := routes.CreateRouter()
	r.Run()
	// http.HandleFunc("/", service.GetAllEmployees)
	// http.HandleFunc("/test", logging(service.GetAllEmployees))
	// fmt.Println("listening on port 8080")
	// http.ListenAndServe(":8080", nil)

}
