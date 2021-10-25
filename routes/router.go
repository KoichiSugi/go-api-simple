package routes

import (
	"fmt"
	"git-clones/go-api-simple/config"
	"git-clones/go-api-simple/data"
	"git-clones/go-api-simple/errorhandling"
	"git-clones/go-api-simple/mysql"
	"git-clones/go-api-simple/repository"
	"git-clones/go-api-simple/routes/middlewares"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var Repo repository.Repository

func init() {
	repo, err := mysql.NewMySQLRepository("mysql", config.CreateConfig(), 3, 3)
	if err != nil {
		panic(err)
	}
	Repo = repo
}

func GetMySqlRepo() repository.Repository {
	return Repo
}

func SetUpRouter(repo repository.Repository) *gin.Engine {
	log.Println("test repo:", repo)
	router := gin.Default()
	middlewares.OutputLog()          //output logs
	router.Use(middlewares.Logger()) //call middleware
	group1 := router.Group("/employees")
	{
		group1.GET("/", func(c *gin.Context) {
			FindAll(c)
		})
		group1.GET("/:id", func(c *gin.Context) {
			Find(c)
		})
		group1.POST("/", func(c *gin.Context) {
			Create(c)
		})
		group1.DELETE("/:id", func(c *gin.Context) {
			Delete(c)
		})
		group1.PUT("/:id", func(c *gin.Context) {
			Update(c)
		})
	}
	return router
}

func Find(c *gin.Context) {
	id := c.Params.ByName("id")
	emp, err := Repo.GetEmployeeById(id)
	if err != nil {
		if err != nil {
			c.JSON(500, &errorhandling.RequestError{Context: "FindSuper calling getEmployeeById", Code: errorhandling.Internal, Message: err.Error()})
			return
		}
	}
	super, err := getSuper(emp.Id) //get super balance
	if err != nil {
		c.JSON(500, &errorhandling.RequestError{Context: "FindSuper calling getSuper func", Code: errorhandling.Internal, Message: err.Error()})
		return
	}
	if super == 0 {
		c.JSON(500, &errorhandling.RequestError{Context: "FindSuper: failed retrieving super", Code: errorhandling.Internal, Message: "error retrieving super"})
		return
	}
	emp.SuperBalance = super
	c.JSON(200, emp)
}

func FindAll(c *gin.Context) {
	emps, err := Repo.GetAllEmployees()
	if err != nil {
		c.JSON(int(errorhandling.BadRequest), err)
		return
	}
	var wg sync.WaitGroup
	for i := 0; i < len(emps)-1; i++ {
		index := i
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			super, err := getSuper(emps[index].Id)
			if err != nil {
				emps[index].SuperBalance = 0
			}
			emps[index].SuperBalance = super
		}(index)
	}
	wg.Wait()
	c.JSON(http.StatusOK, emps)
}

func getSuper(id string) (float64, error) {
	var url = "http://localhost:3000/ato/employee/?/balance"
	url = strings.Replace(url, "?", id, 1)
	log.Println("Sending request to this url: ", url)
	resp, err := http.Get(url)
	if err != nil {
		return 0.0, &errorhandling.RequestError{Context: "getSuper calling ato api", Code: errorhandling.Internal, Message: err.Error()}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) //convert reponse to byte[]
	if err != nil {
		log.Println("err in coverting response to byte[]:", err)
		return 0.0, &errorhandling.RequestError{Context: "getSuper.ioutil.ReadAll(resp.Body)", Code: errorhandling.Internal, Message: err.Error()}
	}

	super, err := strconv.ParseFloat(string(body), 64) //convert stirng to float
	if err != nil {
		log.Println("err in coversion:", err)
		return 0.0, &errorhandling.RequestError{Context: string(body), Code: errorhandling.Internal, Message: err.Error()}
	}
	return super, nil
}

// func GetAll(c *gin.Context) {
// 	emps, err := Repo.GetAllEmployees()
// 	if err != nil {
// 		c.JSON(int(errorhandling.BadRequest), err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, emps)
// }

// func Find(c *gin.Context) {
// 	id := c.Params.ByName("id")
// 	emp, err := Repo.GetEmployeeById(id)
// 	if err != nil {
// 		c.JSON(int(errorhandling.NotFound), errorhandling.WrapError("mysql.GetEmployeeById", errorhandling.NotFound, err.Error()))
// 		return
// 	}
// 	c.JSON(http.StatusOK, emp)
// }

func Create(c *gin.Context) {
	var emp data.Employee
	v := validator.New()
	if err := c.BindJSON(&emp); err != nil {
		err := v.Struct(emp)
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e)
		}
		c.JSON(int(errorhandling.BadRequest), &errorhandling.RequestError{Context: " CreateEmployee c.BindJson", Code: errorhandling.BadRequest, Message: err.Error()})
		return
	}

	emp, err := Repo.CreateEmployee(emp)
	if err != nil {
		c.JSON(int(errorhandling.BadRequest), errorhandling.WrapError("mysql.CreateEmployee", errorhandling.BadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, emp)

}

func Delete(c *gin.Context) {
	id := c.Params.ByName("id")
	err := Repo.DeleteEmployee(id)
	if err != nil {
		if customErr, ok := err.(*errorhandling.RequestError); ok {
			c.JSON(int(customErr.Code), errorhandling.WrapError(customErr.Context, customErr.Code, customErr.Message))
			return
		}
		c.JSON(404, err.Error()) //generic error
		return
	}
	c.String(200, fmt.Sprintf("emp id %v has been deleted", id))
}

func Update(c *gin.Context) {
	id := c.Params.ByName("id")
	var emp data.Employee
	if err := c.BindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, errorhandling.WrapError("router.Update", errorhandling.BadRequest, err.Error()))
		return
	}
	emp.Id = id
	emp, err := Repo.UpdateEmployee(emp)
	if err != nil {
		if customErr, ok := err.(*errorhandling.RequestError); ok {
			c.JSON(int(customErr.Code), errorhandling.WrapError(customErr.Context, customErr.Code, customErr.Message))
			return
		}
		c.JSON(404, err.Error()) //generic error
		return
	}
	c.JSON(http.StatusOK, emp)
}
