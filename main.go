package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/zolachel/gofinal/middleware"
	"github.com/zolachel/gofinal/task"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//db, err = sql.Open("postgres", "postgres://gosctihb:CqOz6dVYlooEBPY4quY9KHvySa2OmADZ@arjuna.db.elephantsql.com:5432/gosctihb")
	if err != nil {
		log.Fatal(err)
	}
}

func setupRouter() *gin.Engine {
	route := gin.Default()

	route.Use(middleware.Auth)

	handler := task.Handler{DB: db}

	handler.CreateCustomerTable()

	route.POST("/customers", handler.CreateCustomerHandler)
	route.GET("/customers/:id", handler.GetCustomerByIDHandler)
	route.GET("/customers", handler.GetCustomersHandler)
	route.PUT("/customers/:id", handler.UpdateCustomerHandler)
	route.DELETE("/customers/:id", handler.DeleteCustomerHandler)

	return route
}

func main() {
	route := setupRouter()

	route.Run(":2009") //run port ":2009"

	defer db.Close()
}
