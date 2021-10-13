package main

import (
	"emp-simple/config"
	"emp-simple/routes"
)

// func logging(f http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.URL.Path)
// 		f(w, r)
// 	}
// }

func main() {
	config.ConnectDb()
	defer config.Db.Close()
	routes.Router()
	// http.HandleFunc("/", service.GetAllEmployees)
	// http.HandleFunc("/test", logging(service.GetAllEmployees))
	// fmt.Println("listening on port 8080")
	// http.ListenAndServe(":8080", nil)

}
