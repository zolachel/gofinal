package task

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//Customer is a struct
type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

//Handler ...
type Handler struct {
	DB *sql.DB
}

//CreateCustomerTable ...
func (h *Handler) CreateCustomerTable() {
	createTb := `
	CREATE TABLE IF NOT EXISTS customer (
	id SERIAL PRIMARY KEY,
	name TEXT,
	email TEXT,
	status TEXT
	);
	`

	if _, err := h.DB.Exec(createTb); err != nil {
		log.Fatal("can't create table", err)
	}
}

//CreateCustomerHandler ...
func (h *Handler) CreateCustomerHandler(context *gin.Context) {
	cus := Customer{}
	if err := context.ShouldBindJSON(&cus); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := h.DB.QueryRow("INSERT INTO customer (name, email, status) values ($1, $2, $3) RETURNING id", cus.Name, cus.Email, cus.Status)

	if err := row.Scan(&cus.ID); err != nil {
		context.JSON(http.StatusInternalServerError, err)
	} else {
		context.JSON(http.StatusCreated, cus)
	}
}

//GetCustomerByIDHandler ...
func (h *Handler) GetCustomerByIDHandler(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}

	stmt, err := h.DB.Prepare("SELECT name, email, status FROM customer where id=$1")
	if err != nil {
		context.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	var name, email, status string

	if err := row.Scan(&name, &email, &status); err != nil {
		context.JSON(http.StatusOK, Customer{})
	} else {
		context.JSON(http.StatusOK, Customer{ID: id, Name: name, Email: email, Status: status})
	}
}

//GetCustomersHandler ...
func (h *Handler) GetCustomersHandler(context *gin.Context) {
	stmt, err := h.DB.Prepare("SELECT id, name, email, status FROM customer")
	if err != nil {
		context.JSON(http.StatusInternalServerError, err)
		return
	}

	if rows, err := stmt.Query(); err != nil {
		context.JSON(http.StatusInternalServerError, err)
	} else {
		result := []*Customer{}

		for rows.Next() {
			var id int
			var name, email, status string
			err := rows.Scan(&id, &name, &email, &status)
			if err != nil {
				log.Fatal("can't Scan row into variable", err)
			} else {
				result = append(result, &Customer{ID: id, Name: name, Email: email, Status: status})
			}
		}

		context.JSON(http.StatusOK, result)
	}
}

//UpdateCustomerHandler ...
func (h *Handler) UpdateCustomerHandler(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}

	cus := Customer{}
	if err = context.ShouldBindJSON(&cus); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := h.DB.Prepare("UPDATE customer SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		context.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := stmt.Exec(id, cus.Name, cus.Email, cus.Status); err != nil {
		context.JSON(http.StatusInternalServerError, err)
	} else {
		cus.ID = id
		context.JSON(http.StatusOK, cus)
	}
}

//DeleteCustomerHandler ...
func (h *Handler) DeleteCustomerHandler(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}

	stmt, err := h.DB.Prepare("DELETE FROM customer WHERE id=$1;")
	if err != nil {
		context.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := stmt.Exec(id); err != nil {
		context.JSON(http.StatusInternalServerError, err)
	} else {
		context.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
	}
}
