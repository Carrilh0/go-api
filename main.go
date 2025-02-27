package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Student struct {
	ID   int    `json:id`
	Name string `json:name`
	Age  int    `json:age`
}

var Students = []Student{
	Student{ID: 1, Name: "Vitor", Age: 24},
	Student{ID: 2, Name: "Gabriel", Age: 27},
}

func routeHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func routeGetStudents(c *gin.Context) {
	c.JSON(http.StatusOK, Students)
}

func routePostStudent(c *gin.Context) {
	var student Student

	err := c.Bind(&student)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Nao foi possivel obter o payload",
		})
		return
	}

	student.ID = Students[len(Students)-1].ID + 1
	Students = append(Students, student)

	c.JSON(http.StatusCreated, student)
}

func routePutStudent(c *gin.Context) {
	var student Student
	var studentLocal Student
	var newStudents []Student

	err := c.Bind(&student)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Nao foi possivel obter o payload",
		})
		return
	}

	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Nao foi possivel obter o payload",
		})
		return
	}

	for _, studentElement := range Students {
		if studentElement.ID == id {
			studentLocal = studentElement
		}
	}

	if studentLocal.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Nao foi possivel encontrar o estudante",
		})
		return
	}

	studentLocal.Name = student.Name
	studentLocal.Age = student.Age

	for _, studentElement := range Students {
		if id == studentElement.ID {
			newStudents = append(newStudents, studentLocal)
		} else {
			newStudents = append(newStudents, studentElement)
		}
	}

	Students = newStudents

	c.JSON(http.StatusCreated, studentLocal)
}

func routeDeleteStudent(c *gin.Context) {
	var newStudents []Student
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Nao foi possivel obter o id",
		})
		return
	}

	for _, studentElement := range Students {
		if studentElement.ID != id {
			newStudents = append(newStudents, studentElement)
		}
	}

	Students = newStudents

	c.JSON(http.StatusOK, gin.H{
		"message": "Estudante excluido com sucesso",
	})
}

func routeGetStudent(c *gin.Context) {
	var student Student
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Nao foi possivel obter o id",
		})
		return
	}

	for _, studentElement := range Students {
		if studentElement.ID == id {
			student = studentElement
		}
	}

	if student.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Nao foi possivel obter o estudante",
		})
		return
	}

	c.JSON(http.StatusOK, student)
}

func main() {
	service := gin.Default()

	getRoutes(service)
	service.Run()
}

func getRoutes(c *gin.Engine) *gin.Engine {
	c.GET("/health", routeHealth)

	groupStudents := c.Group("/students")
	groupStudents.GET("/", routeGetStudents)
	groupStudents.POST("/", routePostStudent)
	groupStudents.PUT("/:id", routePutStudent)
	groupStudents.DELETE("/:id", routeDeleteStudent)
	groupStudents.GET("/:id", routeGetStudent)

	return c
}
