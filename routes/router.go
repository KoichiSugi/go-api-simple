package routes

import (
	"fmt"
	"git-clones/go-api-simple/config"
	"git-clones/go-api-simple/data"
	"git-clones/go-api-simple/errorhandling"
	"git-clones/go-api-simple/mysql"
	"git-clones/go-api-simple/repository"
	"log"
	"net/http"

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

type rootHandler func(http.ResponseWriter, *http.Request) error

func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r) //call root handler
	if err != nil {
		return
	}
}

func GetMySqlRepo() repository.Repository {
	return Repo
}

func SetUpRouter(repo repository.Repository) *gin.Engine {
	log.Println("test repo:", repo)
	router := gin.Default()
	group1 := router.Group("/employees")
	{
		group1.GET("/", func(ctx *gin.Context) {
			GetAll(ctx)
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

func GetAll(c *gin.Context) {
	emps, err := Repo.GetAllEmployees()
	if err != nil {
		c.JSON(int(errorhandling.BadRequest), err)
		return
	}
	c.JSON(http.StatusOK, emps)
}

func Find(c *gin.Context) {
	id := c.Params.ByName("id")
	emp, err := Repo.GetEmployeeById(id)
	if err != nil {
		c.JSON(int(errorhandling.NotFound), errorhandling.WrapError("mysql.GetEmployeeById", errorhandling.NotFound, err.Error()))
		return
	}
	c.JSON(http.StatusOK, emp)
}

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
