package router

import (
	"configserver/app/crud"
	models "configserver/app/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type DatabaseNeo4j struct {
	Driver  *neo4j.DriverWithContext
	Context *context.Context
}

func (connection *DatabaseNeo4j) HandleCreateConfig(w http.ResponseWriter, r *http.Request) {
	var config models.Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	crud.CreateConfig(config, *connection.Driver, *connection.Context)

	// Do something with the Person struct...
	fmt.Fprintf(w, "Person: %+v", config)
	w.Write([]byte("Gorilla!\n"))
}
