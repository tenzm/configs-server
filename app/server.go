package main

import (
	routes "configserver/routes"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	r := mux.NewRouter()

	dbUri := "neo4j://neo4j:7687"
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "admin123", ""))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	defer driver.Close(ctx)

	var Connection routes.DatabaseNeo4j = routes.DatabaseNeo4j{Driver: &driver, Context: &ctx}

	r.HandleFunc("/config", Connection.HandleCreateConfig).Methods("POST")
	r.HandleFunc("/config", Connection.HandleUpdateConfig).Methods("PUT")
	r.HandleFunc("/config", Connection.HandleGetConfig).Methods("GET").Queries("service", "{service}", "version", "{version}")
	r.HandleFunc("/config", Connection.HandleGetConfig).Methods("GET").Queries("service", "{service}")
	r.HandleFunc("/config/delete", Connection.HandleDeleteConfig).Methods("DELETE").Queries("service", "{service}")

	log.Fatal(http.ListenAndServe(":8000", r))
}
