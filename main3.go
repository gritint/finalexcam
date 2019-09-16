package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func getOneCustomer(c *gin.Context) {

	rowId := c.Params.ByName("id")
	var id int
	var name, email, status string

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		log.Fatal("can'tprepare query one row statment", err)
	}

	row := stmt.QueryRow(rowId)
	err = row.Scan(&id, &name, &email, &status)
	if err != nil {
		log.Fatal("can't Scan row into variables", err)
	}
	cus := Todo{
		ID:     id,
		Name:   name,
		Email:  email,
		Status: status,
	}
	c.JSON(http.StatusOK, cus)
}

func getAllCustomer(c *gin.Context) {

	todolist := []Todo{}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	fmt.Println("URL : ", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers")
	if err != nil {
		log.Fatal("can't prepare query all customers statment", err)
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("can't query all customer", err)
	}
	for rows.Next() {
		var id int
		var name, email, status string
		err := rows.Scan(&id, &name, &email, &status)
		if err != nil {
			log.Fatal("can't Scan row into variable", err)
		}
		fmt.Println(id, name, email, status)
		todo := Todo{id, name, email, status}
		todolist = append(todolist, todo)
		log.Println("Row id:", id)
	}

	c.JSON(http.StatusOK, todolist)
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

	row := db.QueryRow("INSERT INTO customers (name, email, status) values ($1, $2, $3) RETURNING id", t.Name, t.Email, t.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		fmt.Println("can't scan id", err)
		return
	}
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		log.Fatal("can'tprepare query one row statment", err)
	}
	row = stmt.QueryRow(id)

	var name, email, status string
	err = row.Scan(&id, &name, &email, &status)
	if err != nil {
		log.Fatal("can't Scan row into variables", err)
	}
	cus := Todo{
		ID:     id,
		Name:   name,
		Email:  email,
		Status: status,
	}

	c.JSON(http.StatusCreated, cus)

}

func updateCustomer(c *gin.Context) {
	rowId := c.Params.ByName("id")
	var t Todo
	var id int
	err := c.ShouldBindJSON(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		log.Fatal("can't prepare statment update", err)
	}
	if _, err = stmt.Exec(rowId, t.Name, t.Email, t.Status); err != nil {
		log.Fatal("error execute update ", err)
	}

	stmt, err = db.Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		log.Fatal("can'tprepare query one row statment", err)
	}
	row := stmt.QueryRow(rowId)
	var name, email, status string
	err = row.Scan(&id, &name, &email, &status)
	if err != nil {
		log.Fatal("can't Scan row into variables", err)
	}
	cus := Todo{
		ID:     id,
		Name:   name,
		Email:  email,
		Status: status,
	}

	c.JSON(http.StatusOK, cus)
}

func deleteCustomer(c *gin.Context) {
	rowId := c.Params.ByName("id")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM customers WHERE id=$1;")
	if err != nil {
		log.Fatal("can't prepare statment delete", err)
	}
	if _, err = stmt.Exec(rowId); err != nil {
		log.Fatal("error execute delete ", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})

}

func authMiddelware(c *gin.Context) {
	fmt.Println("This is a middelwear")
	token := c.GetHeader("Authorization")
	if token != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized."})
		c.Abort()
		return
	}
	c.Next()
	fmt.Println("after in middelware")
}

func main() {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()
	createTb := `
	CREATE TABLE IF NOT EXISTS customers (
	id SERIAL PRIMARY KEY,
	name TEXT,
	email TEXT,
	status TEXT

	);
	`
	_, err = db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
	fmt.Println("create table success")

	r := gin.Default()
	r.POST("/customers", createCustomer)
	r.GET("/customers/:id", getOneCustomer)
	r.GET("/customers", getAllCustomer)
	r.PUT("/customers/:id", updateCustomer)
	r.DELETE("/customers/:id", deleteCustomer)
	// r.GET("/customers", authMiddelware)
	r.Run(":2019")
}
