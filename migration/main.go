package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"gitlab.com/stephenjaya99/tax-service/database"
	m "gitlab.com/stephenjaya99/tax-service/model"

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
	db.Migrate(&m.Tax{})
	db.Migrate(&m.TaxCode{})
	db.CreateTaxCode(1, "Food")
	db.CreateTaxCode(2, "Tobacco")
	db.CreateTaxCode(3, "Entertainment")
}
