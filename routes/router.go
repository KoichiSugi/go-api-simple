package routes

import (
	"git-clones/go-api-simple/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type rootHandler func(http.ResponseWriter, *http.Request) error

func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r) //call root handler
	if err != nil {
		return
	}
}

func SetUpRouter(repo repository.Repository) *gin.Engine {
	log.Println("test repo:", repo)
	router := gin.Default()
	group1 := router.Group("/employee")
	{
		group1.GET("/", func(ctx *gin.Context) {
			repo.GetAllEmployees(ctx)
		})
		group1.GET("/:id", func(c *gin.Context) {
			repo.GetEmployeeById(c)
		})
		group1.POST("/", func(c *gin.Context) {
			repo.CreateEmployee(c)
		})
		group1.DELETE("/:id", func(c *gin.Context) {
			repo.DeleteEmployee(c)
		})
		group1.PUT("/:id", func(c *gin.Context) {
			repo.UpdateEmployee(c)
		})
	}
	return router
}
