package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/priyansh-garg1/dating-app/actions"
	controllers "github.com/priyansh-garg1/dating-app/controllers"
	"github.com/priyansh-garg1/dating-app/db"
)

func main() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	neo4jDriver, err := db.ConnectToDb()
	if err != nil {
		fmt.Println("Failed to connect to Neo4j:", err)
		return
	}

	defer neo4jDriver.Close(context.Background())

	dbInstance := actions.NewDatabase(neo4jDriver)

	
	userController := controllers.NewUserController(dbInstance)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/health",func (c *gin.Context){
		c.JSON(200,gin.H{
			"message": "Healthy",
		})
	})
	router.GET("/users/:applicationId", userController.UserHandler)
	router.GET("/notconnecteduser/:applicationId", userController.GetUsersWithNoConnectionHandler)
	router.GET("/connection/:applicationId/:userId", userController.SwipeHandler)

	router.Run(":8080")
	
}
