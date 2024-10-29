package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/priyansh-garg1/dating-app/actions"
	controllers "github.com/priyansh-garg1/dating-app/controllers"
	"github.com/priyansh-garg1/dating-app/db"
)

func main() {


	neo4jDriver, err := db.ConnectToDb()
	if err != nil {
		fmt.Println("Failed to connect to Neo4j:", err)
		return
	}
	defer neo4jDriver.Close(context.Background())

	dbInstance := actions.NewDatabase(neo4jDriver)

	
	userController := controllers.NewUserController(dbInstance)

	router := gin.Default()

	router.GET("/users/:applicationId", userController.UserHandler)
	router.GET("/notconnecteduser/:applicationId", userController.GetUsersWithNoConnectionHandler)
	router.GET("/connection/:applicationId/:userId", userController.SwipeHandler)

	router.Run(":8080")
	
}
