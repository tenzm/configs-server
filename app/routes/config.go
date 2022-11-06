package router

import (
	"configserver/crud"
	models "configserver/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	data := crud.GetConfig(config.Service, -1, *(*connection).Driver, *(*connection).Context)

	if len(data) > 0 {
		http.Error(w, "The configuration is already in use\n", http.StatusAlreadyReported)
	} else {
		crud.CreateConfig(config, *(*connection).Driver, *(*connection).Context)
		w.Write([]byte(fmt.Sprintf("The configuration has been saved. Current version: 1\n")))
	}
}

func (connection *DatabaseNeo4j) HandleGetConfig(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	version_str := vars["version"]
	var version int = -1

	version, err := strconv.Atoi(version_str)

	if err != nil {
		version = -1
	}

	service := vars["service"]
	fmt.Println(vars["version"])

	result := crud.GetConfig(service, int64(version), *(*connection).Driver, *(*connection).Context)

	if len(result) == 0 {
		http.Error(w, "The configuration was not found\n", http.StatusNotFound)
	} else {

		jsonStr, err := json.Marshal(result)

		if err != nil {
			fmt.Println(err)
		}

		w.Write([]byte(jsonStr))
	}
}

func (connection *DatabaseNeo4j) HandleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var config models.Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	new_version := crud.UpdateConfig(config, *(*connection).Driver, *(*connection).Context)

	if new_version == -1 {
		w.Write([]byte(fmt.Sprintf("The configuration has been saved. Current version: %d\n", 1)))
	} else {
		w.Write([]byte(fmt.Sprintf("The configuration has been updated. Current version: %d\n", new_version)))
	}

}
