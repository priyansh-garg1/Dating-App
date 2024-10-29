package db

import (
	"context"
	"fmt"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func ConnectToDb() (neo4j.DriverWithContext, error) {
    ctx := context.Background()
    dbUri := os.Getenv("DB_URI")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")

    fmt.Println(dbUri,dbUser)

    driver, err := neo4j.NewDriverWithContext(
        dbUri,
        neo4j.BasicAuth(dbUser, dbPassword, ""))
    if err != nil {
        return nil, err
    }

    err = driver.VerifyConnectivity(ctx)
    if err != nil {
        defer driver.Close(ctx)
        return nil, err
    }

    fmt.Println("Connection established.")
    return driver, nil
}
