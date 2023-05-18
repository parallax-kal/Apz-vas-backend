package main

import (
	"database/sql"
    _ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
	"fmt"
)



func main() {
	router := gin.Default()
	db, err := connectDB()
    if err != nil {
		fmt.Println(err)
    }
    defer db.Close()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	router.Run(":8080")
}


func connectDB() (*sql.DB, error) {
    connectionString := "postgres://postgres:KALISA123.@localhost:5432/apzVas?sslmode=disable"
    // Replace 'username', 'password', 'localhost', 'dbname' with your actual database credentials

    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, err
    }

    // Ping the database to ensure the connection is established
    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}
