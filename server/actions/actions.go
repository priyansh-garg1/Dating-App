package actions

import (
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type User struct {
    ApplicationId string `json:"application_id"`
    Name          string `json:"name"`
    Email         string `json:"email"`
}

func FetchUserByApplicationId(driver neo4j.DriverWithContext, applicationId string) (*User, error) {

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
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
		Name:           fmt.Sprintf("%v", node.Props["name"]),
        Email:          fmt.Sprintf("%v", node.Props["email"]),
        ApplicationId:  fmt.Sprintf("%v", node.Props["applicationId"]),
    }

    return user, nil
}


func InsertUser(driver neo4j.DriverWithContext, user User) error {
    

    session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
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


func GetUsersWithNoConnection(driver neo4j.DriverWithContext, applicationId string) ([]User,error) {

    session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
    defer session.Close(context.Background())

    query := `
        MATCH (cu: User {applicationId: $applicationId}) MATCH (ou: User) WHERE NOT (cu)-[:LIKE|:DILIKE]->(ou) AND cu <> ou RETURN ou`
    parameters := map[string]interface{}{
        "applicationId": applicationId,
    }

    result, err := session.Run(context.Background(), query, parameters)
    if err != nil {
        fmt.Println("Error executing Get user query with no connection")
    }

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
