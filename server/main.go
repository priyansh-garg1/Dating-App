package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/priyansh-garg1/dating-app/actions"
	"github.com/priyansh-garg1/dating-app/db"
)

var neo4jDriver neo4j.DriverWithContext

func UserHandler(c *gin.Context) {
    applicationId := c.Param("applicationId") 

    user, _ := actions.FetchUserByApplicationId(neo4jDriver, applicationId)

    if user == nil {
        newUser := actions.User{
            Email:         "a1@gmail.com",
            Name:          "hello",
            ApplicationId: applicationId,
        }

        err := actions.InsertUser(neo4jDriver, newUser)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": "User created successfully",
            "user":    newUser,
        })
        return
    }

    c.JSON(http.StatusOK, user)
}

func GetUsersWithnoConnectionHandler(c *gin.Context) {
    applicationId := c.Param("applicationId") 
	users,_ := actions.GetUsersWithNoConnection(neo4jDriver,applicationId)

	if users == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
        return
	}

	c.JSON(http.StatusOK, users)
	return

}

func main() {
    var err error

    neo4jDriver, err = db.ConnectToDb()
    if err != nil {
        fmt.Println("Failed to connect to Neo4j:", err)
        return
    }
    defer neo4jDriver.Close(context.Background()) 

    router := gin.Default()

    router.GET("/users/:applicationId", UserHandler)
    router.GET("/notconnecteduser/:applicationId", GetUsersWithnoConnectionHandler)

    router.Run(":8080") 
}
