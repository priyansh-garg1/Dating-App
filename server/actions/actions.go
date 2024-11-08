package actions

import (
	"context"
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Database struct {
	driver neo4j.DriverWithContext
}

func NewDatabase(driver neo4j.DriverWithContext) *Database {
	return &Database{driver: driver}
}

type User struct {
	ApplicationId string `json:"application_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
}

func (db *Database) FetchUserByApplicationId(applicationId string) (*User, error) {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	query := `
		MATCH (u:User {applicationId: $applicationId})
		RETURN u`
	parameters := map[string]interface{}{
		"applicationId": applicationId,
	}

	result, err := session.Run(context.Background(), query, parameters)
	if err != nil {
		return nil, err
	}

	if !result.Next(context.Background()) {
		return nil, errors.New("user not found")
	}

	node := result.Record().Values[0].(neo4j.Node)

	user := &User{
		Name:          fmt.Sprintf("%v", node.Props["name"]),
		Email:         fmt.Sprintf("%v", node.Props["email"]),
		ApplicationId: fmt.Sprintf("%v", node.Props["applicationId"]),
	}

	return user, nil
}

func (db *Database) InsertUser(user User) error {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	query := `
		CREATE (u: User {name: $name, applicationId: $applicationId, email: $email})
		RETURN u`
	parameters := map[string]interface{}{
		"name":          user.Name,
		"applicationId": user.ApplicationId,
		"email":         user.Email,
	}

	result, err := session.Run(context.Background(), query, parameters)
	if err != nil {
		return err
	}

	if !result.Next(context.Background()) {
		return errors.New("failed to insert user")
	}

	return nil
}

func (db *Database) GetUsersWithNoConnection(applicationId string) ([]User, error) {
	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	query := `
		MATCH (cu: User {applicationId: $applicationId}) 
		MATCH (ou: User) 
		WHERE NOT (cu)-[:LIKE|:DILIKE]->(ou) AND cu <> ou 
		RETURN ou`
	parameters := map[string]interface{}{
		"applicationId": applicationId,
	}

	result, err := session.Run(context.Background(), query, parameters)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	var users []User
	for result.Next(context.Background()) {
		node := result.Record().Values[0].(neo4j.Node)

		user := User{
			Email:         fmt.Sprintf("%v", node.Props["email"]),
			Name:          fmt.Sprintf("%v", node.Props["name"]),
			ApplicationId: fmt.Sprintf("%v", node.Props["applicationId"]),
		}

		users = append(users, user)
	}

	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("error processing query result: %v", err)
	}

	return users, nil
}


func (db *Database) Neo4jSwipe(id string,swipe string,userId string) (bool, error) {
	var swipeType string
	if swipe == "left" {
		swipeType = "DISLIKE"
	} else {
		swipeType = "LIKE"
	} 

	session := db.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())
	fmt.Println("Swipe created successfully")

	query := fmt.Sprintf(`
        MATCH (cu: User { applicationId: $id }), (ou: User { applicationId: $userId })
        CREATE (cu)-[:%s]->(ou)`, swipeType)

	parameters := map[string]interface{}{
		"id":     id,
		"userId": userId,
	}

	_, err := session.Run(context.Background(), query, parameters)
	if err != nil {
		return false, fmt.Errorf("error executing query: %v", err)
	}

	if swipeType == "LIKE" {
		query := `MATCH (cu: User{applicationId: $id}), (ou: User{applicationId: $userId}) WHERE (ou)-[:LIKE]->(cu) RETURN ou as MATCH`
		parameters := map[string]interface{}{
			"id":     id,
			"userId": userId,
		}
	
		result, err := session.Run(context.Background(), query, parameters)
		if err != nil {
			return false, fmt.Errorf("error executing query: %v", err)
		}

		if result.Next(context.Background()) {
			return true, nil 
		}
	}


	return false, nil

	
}