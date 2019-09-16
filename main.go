package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	//"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func createCustomer(c *gin.Context) {
	var t Todo
	err := c.ShouldBindJSON(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	fmt.Println("URL : ", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	row := db.QueryRow("INSERT INTO Todo (name, email, status) values ($1, $2, $3) RETURNING ", t.Name, t.Email, t.Status)
	var id int
	err = row.Scan(&name)
	if err != nil {
		fmt.Println("can't scan Name", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func main() {
	r := gin.Default()
	r.POST("/customers", createCustomer)
	r.Run(":2019")
}
