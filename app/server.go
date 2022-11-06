package main

import (
	routes "configserver/app/routes"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	r := mux.NewRouter()

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

	var Connection routes.DatabaseNeo4j = routes.DatabaseNeo4j{Driver: &driver, Context: &ctx}

	// Routes consist of a path and a handler function.
	r.HandleFunc("/config", Connection.HandleCreateConfig).Methods("POST")
	r.HandleFunc("/config", Connection.HandleCreateConfig).Methods("GET")
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
