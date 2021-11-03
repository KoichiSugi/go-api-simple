package routes

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type preferredResponse struct {
	code int                    //httpステータスコード
	body map[string]interface{} //帰ってくる文字列
}

func init() {
}

var format = "2006-01-02"
var mockDOB, _ = time.Parse(format, "2018-12-12")

func TestFindAllFailure(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)

	FindAll(c)
	t.Run("testing", func(t *testing.T) {
		assert.Equal(t, 200, w.Code)
	})
}

func TestFindAllFailure2(t *testing.T) {
	w := httptest.NewRecorder()
	router := gin.Default()
	gin.SetMode(gin.TestMode)

	router.GET("/employees", func(c *gin.Context) {
		FinderAll(c)
	})

	t.Run("get json data", func(t *testing.T) {
		assert.Equal(t, 200, w.Code)
	})
}

func TestFind(t *testing.T) {

	// tests := []struct {
	// 	emp  data.Employee
	// 	want preferredResponse
	// }{
	// 	{
	// 		data.Employee{
	// 			Id:           "1234",
	// 			FirstName:    "Vinodh",
	// 			MiddleName:   "K",
	// 			LastName:     "Landa",
	// 			Gender:       "Male",
	// 			Salary:       555.55,
	// 			DOB:          data.CustomDOB(mockDOB),
	// 			Email:        "vinod@gmail.com",
	// 			Phone:        "01298384",
	// 			AddressLine1: "Lonsdale",
	// 			AddressLine2: "street",
	// 			State:        "vic",
	// 			Postcode:     3000,
	// 			TFN:          "19944991",
	// 			SuperBalance: 4649,
	// 		},
	// 		preferredResponse{
	// 			code: http.StatusCreated,
	// 			body: map[string]interface{}{
	// 				"id":           "1234",
	// 				"firstname":    "Vinodh",
	// 				"middlename":   "K",
	// 				"lastname":     "Landa",
	// 				"gender":       "Male",
	// 				"salary":       555.55,
	// 				"dob":          "2018-12-12",
	// 				"email":        "vinod@gmail.com",
	// 				"phone":        "01298384",
	// 				"state":        "Vic",
	// 				"postcode":     3000,
	// 				"addressline1": "Lonsdale",
	// 				"addressline2": "street",
	// 				"tfn":          "19944991",
	// 				"superbalance": 4649,
	// 			},
	// 		},
	// 	},
	// 	{
	// 		data.Employee{
	// 			Id:           "1234",
	// 			FirstName:    "Vinodh",
	// 			MiddleName:   "K",
	// 			LastName:     "Landa",
	// 			Gender:       "Male",
	// 			Salary:       555.55,
	// 			DOB:          data.CustomDOB(mockDOB),
	// 			Email:        "vinod@gmail.com",
	// 			Phone:        "01298384",
	// 			AddressLine1: "Lonsdale",
	// 			AddressLine2: "street",
	// 			State:        "vic",
	// 			Postcode:     3000,
	// 			TFN:          "19944991",
	// 			SuperBalance: 4649,
	// 		},
	// 		preferredResponse{
	// 			code: http.StatusCreated,
	// 			body: map[string]interface{}{
	// 				"id":           "1234",
	// 				"firstname":    "Vinodh",
	// 				"middlename":   "K",
	// 				"lastname":     "Landa",
	// 				"gender":       "Male",
	// 				"salary":       555.55,
	// 				"dob":          "2018-12-12",
	// 				"email":        "vinod@gmail.com",
	// 				"phone":        "01298384",
	// 				"state":        "Vic",
	// 				"postcode":     3000,
	// 				"addressline1": "Lonsdale",
	// 				"addressline2": "street",
	// 				"tfn":          "19944991",
	// 				"superbalance": 4649,
	// 			},
	// 		},
	// 	},
	// }
	// url := "http://localhost:8080/employees"
	// var jsonStr = []byte(
	// 	`{"id":     "1234",
	// 		"firstname":    "Vinodh",
	// 		"middlename":   "K",
	// 		"lastname":     "Landa",
	// 		"gender":       "Male",
	// 		"salary":       555.55,
	// 		"dob":          "2018-12-12",
	// 		"email":        "vinod@gmail.com",
	// 		"phone":        "01298384",
	// 		"state":        "Vic",
	// 		"postcode":     3000,
	// 		"addressline1": "Lonsdale",
	// 		"addressline2": "street",
	// 		"tfn":          "19944991",
	// 		"superbalance": 4649,}`)

	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// if err != nil {
	// 	panic(err)
	// }
	// req.Header.Set("Content-Type", "application/json")
	// log.Println(tests)
	// for i, tt := range tests {
	// 	requestBody := strings.NewReader(
	// 		"id=" + tt.emp.Id +
	// 			"&firstname=" + tt.emp.FirstName +
	// 			"&middlename=" + tt.emp.MiddleName +
	// 			"&lastname=" + "Lnada" +
	// 			"&gender=" + "Male" +
	// 			"&salary=" + nil +
	// 			"&dob=" + "2018-12-12" +
	// 			"&email=" + "vinod@gmail.com" +
	// 			"&phone=" + "01298384" +
	// 			"&state=" + "Vic" +
	// 			"&postcode=" + 3000 +
	// 			"&addressline1=" + "Lonsdale" +
	// 			"&addressline2=" + "street" +
	// 			"&tfn=" + "19944991" +
	// 			"&superbalance=" + 4649,
	// 	)
	// }

}
