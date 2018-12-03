package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"gitlab.com/stephenjaya99/tax-service/controller"
	"gitlab.com/stephenjaya99/tax-service/database"
	"gitlab.com/stephenjaya99/tax-service/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using to container ENV only.")
	}

	r := gin.Default()
	r.Use(cors.Default())

	dbUsername := os.Getenv("POSTGRES_USERNAME")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DATABASE_NAME")

	client, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable",
			dbHost,
			dbUsername,
			dbName,
			dbPassword,
		),
	)
	if err != nil {
		log.Println("Error connecting to DB")
		panic(err)
	}

	db := database.New(client)

	port := os.Getenv("WEB_PORT")

	controllerOpt := controller.ControllerOpt{
		Database: db,
	}

	controller := controller.New(controllerOpt)
	handler := handler.New(controller)

	r.GET("/", handler.Ping)
	r.POST("/tax", handler.CreateTax)
	r.GET("/tax", handler.RetrieveAllTaxes)

	r.Run(port)
	fmt.Println("Running Tax service on Port %s", port)
}
