package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git-clones/go-api-simple/config"
	"git-clones/go-api-simple/data"
	"git-clones/go-api-simple/errorhandling"
	"git-clones/go-api-simple/mysql"
	"git-clones/go-api-simple/repository"
	"git-clones/go-api-simple/routes/middlewares"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

//simple json body validation
func PostValidate() gin.HandlerFunc {

	return func(c *gin.Context) {
		var emp data.Employee
		bodyCopy := new(bytes.Buffer)
		// Read the whole body
		_, err := io.Copy(bodyCopy, c.Request.Body)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
		}
		bodyData := bodyCopy.Bytes()
		// Replace the body with a reader that reads from the buffer
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bodyData))
		// validate bodyCopy
		log.Println("body is copied unmarshalled -> ", emp)
		v := validator.New()
		if err = json.Unmarshal(bodyCopy.Bytes(), &emp); err != nil {
			err = v.Struct(emp)
			for _, e := range err.(validator.ValidationErrors) {
				log.Println(e)
			}
			c.JSON(int(errorhandling.BadRequest), &errorhandling.RequestError{Context: " ValidatePostPut middleware json.Unmarshal", Code: errorhandling.BadRequest, Message: err.Error()})
		}

	}
}

func SetUpRouter(repo repository.Repository) *gin.Engine {
	log.Println("test repo:", repo)
	router := gin.Default()
	middlewares.OutputLog()          //output logs
	router.Use(middlewares.Logger()) //call middleware
	group1 := router.Group("/employees")
	{
		group1.GET("/", func(c *gin.Context) {
			//FindAll(c)
			FinderAll(c)
		})
		group1.GET("/:id", func(c *gin.Context) {
			Find(c)
		})
		group1.POST("/", PostValidate(), Create)

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
			c.JSON(404, &errorhandling.RequestError{Context: "Find getSuper calling getEmployeeById", Code: errorhandling.NotFound, Message: err.Error()})
			return
		}
	}
	super, err := getSuper(emp.Id) //get super balance
	if err != nil {
		c.JSON(500, &errorhandling.RequestError{Context: "Find getSuper calling getSuper func", Code: errorhandling.Internal, Message: err.Error()})
		return
	}
	if super == nil {
		c.JSON(500, &errorhandling.RequestError{Context: "Find getSuper : failed retrieving super", Code: errorhandling.Internal, Message: "error retrieving super"})
		return
	}
	emp.SuperBalance = *super
	c.JSON(200, emp)
}

func FindAll(c *gin.Context) {
	emps, err := Repo.GetAllEmployees()
	if err != nil {
		c.JSON(int(errorhandling.BadRequest), err)
		return
	}
	var wg sync.WaitGroup
	for i := 0; i < len(emps); i++ {
		index := i
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			super, err := getSuper(emps[index].Id)
			if err != nil {
				emps[index].SuperBalance = 0
			} else {
				emps[index].SuperBalance = *super
			}
		}(index)
	}
	wg.Wait()
	c.JSON(http.StatusOK, emps)
}

func FinderAll(c *gin.Context) {
	var resultC = make(chan data.SuperDetails)
	emps, err := Repo.GetAllEmployees()
	if err != nil {
		c.JSON(int(errorhandling.BadRequest), err)
		return
	}

	go workerPool(emps, resultC)
	smap := make(map[string]float64)
	for e := range resultC {
		log.Println("log: ", e)
		smap[e.EmpId] = e.SuperBalance
	}
	log.Println("AM i here??", smap)
	//may need brocker?
	for i := 0; i < len(emps); i++ {
		_, ok := smap[emps[i].Id]
		if ok {
			emps[i].SuperBalance = smap[emps[i].Id]
		} else {
			emps[i].SuperBalance = 0
		}
	}
	c.JSON(http.StatusOK, emps)

}
func workerPool(emps []data.Employee, resultC chan data.SuperDetails) {
	var wg sync.WaitGroup
	for i := 0; i < len(emps); i++ {
		wg.Add(1)
		go worker(emps[i].Id, &wg, resultC)
	}
	wg.Wait()
	close(resultC)
}

func worker(id string, wg *sync.WaitGroup, resultC chan data.SuperDetails) {
	res, err := getSuper(id)
	if err != nil {
		resultC <- data.SuperDetails{EmpId: id, SuperBalance: 0}
	} else {
		resultC <- data.SuperDetails{EmpId: id, SuperBalance: *res}
	}
	wg.Done()
}

func getSuper(id string) (*float64, error) {
	var url = "http://localhost:3000/ato/employee/?/balance"
	url = strings.Replace(url, "?", id, 1)
	log.Println("Sending request to this url: ", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, &errorhandling.RequestError{Context: "getSuper calling ato api", Code: errorhandling.Internal, Message: err.Error()}
	}
	body, err := resposeToByte(resp)
	if err != nil {
		log.Println("err in coverting response to byte[]:", err)
		return nil, &errorhandling.RequestError{Context: "getsuper resposeToByte", Code: errorhandling.Internal, Message: err.Error()}
	}
	superData, err := UnmarshalSuperDetails(body)
	if err != nil {
		return nil, &errorhandling.RequestError{Context: "UnmarshalSuperDetails())", Code: errorhandling.Internal, Message: err.Error()}
	}
	log.Println("super details ", superData)

	return &superData.SuperBalance, nil
}

func resposeToByte(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body) //convert reponse to byte[]
	if err != nil {
		log.Println("err in coverting response to byte[]:", err)
		return []byte{}, &errorhandling.RequestError{Context: "getSuper.ioutil.ReadAll(resp.Body)", Code: errorhandling.Internal, Message: err.Error()}
	}
	return body, nil
}

func UnmarshalSuperDetails(b []byte) (data.SuperDetails, error) {
	var superData data.SuperDetails
	err := json.Unmarshal([]byte(b), &superData)
	if err != nil {
		return data.SuperDetails{}, &errorhandling.RequestError{Context: "getSuperjson.Unmarshal([]byte(body), &superData)", Code: errorhandling.Internal, Message: err.Error()}
	}
	return superData, nil
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
