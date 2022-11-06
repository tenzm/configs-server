package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Config struct {
	Service string              `json:"service"`
	Data    []map[string]string `json:"data"`
}

func main() {

	jsonFile, err := os.Open("./app/data.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Hello")
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	var config Config

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &config)

	fmt.Println(config.Data)

	// Neo4j 4.0, defaults to no TLS therefore use bolt:// or neo4j://
	// Neo4j 3.5, defaults to self-signed certificates, TLS on, therefore use bolt+ssc:// or neo4j+ssc://
	dbUri := "neo4j://localhost:7687"
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "admin123", ""))
	if err != nil {
		panic(err)
	}
	// Starting with 5.0, you can control the execution of most driver APIs
	// To keep things simple, we create here a never-cancelling context
	// Read https://pkg.go.dev/context to learn more about contexts
	ctx := context.Background()
	// Handle driver lifetime based on your application lifetime requirements  driver's lifetime is usually
	// bound by the application lifetime, which usually implies one driver instance per application
	// Make sure to handle errors during deferred calls
	defer driver.Close(ctx)
	item, err := insertItem(ctx, driver, config.Service, config.Data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", item)
}

func insertItem(ctx context.Context, driver neo4j.DriverWithContext, service string, data []map[string]string) (*Item, error) {
	// Sessions are short-lived, cheap to create and NOT thread safe. Typically create one or more sessions
	// per request in your web application. Make sure to call Close on the session when done.
	// For multi-database support, set sessionConfig.DatabaseName to requested database
	// Session config will default to write mode, if only reads are to be used configure session for
	// read mode.
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	result, err := session.ExecuteWrite(ctx, createItemFn(ctx, service, data))
	if err != nil {
		return nil, err
	}
	return result.(*Item), nil
}

func createItemFn(ctx context.Context, serviceName string, data []map[string]string) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {

		_, err := tx.Run(ctx, "CREATE (service:Service { name: $name }) RETURN service;", map[string]any{
			"name": serviceName,
		})

		if err != nil {
			return nil, err
		}

		for _, d := range data {
			for key, value := range d {
				tx.Run(ctx, fmt.Sprintf("MATCH (service:Service { name: $name }) CREATE (service)-[v:Version{n: %d }]->(key: %s { value: \"%s\" });", 1, key, value), map[string]any{
					"name": serviceName,
				})
				if err != nil {
					return nil, err
				}
			}
		}

		// In face of driver native errors, make sure to return them directly.
		// Depending on the error, the driver may try to execute the function again.

		/*
			record, err := records.Single(ctx)
			if err != nil {
				return nil, err
			}
			// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
			return &Item{
				Name: record.Values[0].(string),
			}, nil
		*/
		return "ok", nil
	}
}

type Item struct {
	Id   int64
	Name string
}
