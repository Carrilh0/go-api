package main

import (
	"net/http"

	"github.com/Carrilh0/aula/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Student struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Age  int       `json:"age"`
}

var Students = []Student{
	{ID: shared.GetUuid(), Name: "Vitor", Age: 24},
	{ID: shared.GetUuid(), Name: "Gabriel", Age: 27},
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

	student.ID = shared.GetUuid()
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

	idString := c.Params.ByName("id")
	id, err := shared.GetUuidByString(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Problema com seu id",
		})
		return
	}

	for _, studentElement := range Students {
		if studentElement.ID == id {
			studentLocal = studentElement
		}
	}

	if studentLocal.ID == shared.GetUuidEmpty() {
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
	idString := c.Params.ByName("id")
	id, err := shared.GetUuidByString(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Problema com seu id",
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
	idString := c.Params.ByName("id")
	id, err := shared.GetUuidByString(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Problema com seu id",
		})
		return
	}

	for _, studentElement := range Students {
		if studentElement.ID == id {
			student = studentElement
		}
	}

	if student.ID == shared.GetUuidEmpty() {
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
