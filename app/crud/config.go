package crud

import (
	"configserver/app/models"
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateConfig(config models.Config, driver neo4j.DriverWithContext, ctx context.Context) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		_, err := tx.Run(ctx, "CREATE (service:Service { name: $name }) RETURN service;", map[string]any{
			"name": config.Service,
		})

		if err != nil {
			return nil, err
		}

		for _, data := range config.Data {
			for key, value := range data {
				tx.Run(ctx, fmt.Sprintf("MATCH (service:Service { name: $name }) CREATE (service)-[v:Version{n: %d }]->(key: %s { value: \"%s\" });", 1, key, value), map[string]any{
					"name": config.Service,
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
	})

	if err != nil {
		return
	}
	return
}
