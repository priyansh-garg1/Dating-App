package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/priyansh-garg1/dating-app/actions"
)

type UserController struct {
	DB *actions.Database
}

func NewUserController(db *actions.Database) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) UserHandler(c *gin.Context) {
	applicationId := c.Param("applicationId")

	user, _ := uc.DB.FetchUserByApplicationId(applicationId)

	if user == nil {
		newUser := actions.User{
			Email:         "a1@gmail.com",
			Name:          "hello",
			ApplicationId: applicationId,
		}

		err := uc.DB.InsertUser(newUser)
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
	return
}

func (uc *UserController) GetUsersWithNoConnectionHandler(c *gin.Context) {
	applicationId := c.Param("applicationId")
	users, err := uc.DB.GetUsersWithNoConnection(applicationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users with no connection"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uc *UserController)SwipeHandler(c *gin.Context) {
	id := c.Param("applicationId")
	userId := c.Param("userId")
	swipe := c.Query("swipe")  

	fmt.Println(id,swipe,userId)

	if swipe != "left" && swipe != "right" {
		c.JSON(http.StatusBadRequest, gin.H{"error": swipe})
		return
	}


	message, err := uc.DB.Neo4jSwipe(id, swipe, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "swipe created successfully",
		"swipe":   message,
	})





}