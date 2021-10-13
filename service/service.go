package service

import (
	"emp-simple/config"
	"emp-simple/data"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllEmployees(c *gin.Context) {
	rows, err := config.Db.Query("SELECT * FROM employee")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var emp data.Employee

		if err := rows.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName,
			&emp.LastName, &emp.Gender, &emp.Salary, &emp.DOB, &emp.Email,
			&emp.Phone, &emp.State, &emp.Postcode, &emp.AddressLine1, &emp.AddressLine2,
			&emp.TFN, &emp.SuperBalance); err != nil {
			log.Fatal(err)
			c.IndentedJSON(http.StatusInternalServerError, nil)
		}
		c.IndentedJSON(http.StatusOK, emp)

	}
	if c.Err(); err != nil {
		log.Fatal(c)
	}

}

// func GetAllEmployees(w http.ResponseWriter, r *http.Request) {
// 	rows, err := config.Db.Query("SELECT * FROM employee")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var emp data.Employee

// 		if err := rows.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName,
// 			&emp.LastName, &emp.Gender, &emp.Salary, &emp.DOB, &emp.Email,
// 			&emp.Phone, &emp.State, &emp.Postcode, &emp.AddressLine1, &emp.AddressLine2,
// 			&emp.TFN, &emp.SuperBalance); err != nil {
// 			log.Fatal(err)
// 		}
// 		json.NewEncoder(w).Encode(emp)
// 	}
// 	if rows.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// }
