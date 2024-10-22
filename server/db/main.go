package db

import (
    "context"
    "fmt"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func ConnectToDb() (neo4j.DriverWithContext, error) {
    ctx := context.Background()
    dbUri := "neo4j+s://c1eb8f54.databases.neo4j.io"
    dbUser := "neo4j"
    dbPassword := "HxH4rFaYFskgFv5EHxYnGkb1NXEQWXFl8xzY09B7DCw"
    
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
