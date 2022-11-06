package crud

import (
	"configserver/models"
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateConfig(config models.Config, driver neo4j.DriverWithContext, ctx context.Context) {
	const startVersion = 1
	session := (driver).NewSession(ctx, neo4j.SessionConfig{})
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		_, err := tx.Run(ctx, "CREATE (service:Service { name: $name, version: $ver  })", map[string]any{
			"name": config.Service,
			"ver":  startVersion,
		})

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for _, data := range config.Data {
			for key, value := range data {
				tx.Run(ctx, fmt.Sprintf("MATCH (service:Service { name: $name, version: $ver }) CREATE (service)-[e:Edge]->(item:Item{key: \"%s\", value: \"%s\"})", key, value), map[string]any{
					"name": config.Service,
					"ver":  startVersion,
				})
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
			}
		}

		return "ok", nil
	})
	session.Close(ctx)

	if err != nil {
		return
	}
	return
}

func DeleteConfig(service string, driver neo4j.DriverWithContext, ctx context.Context) error {
	session := (driver).NewSession(ctx, neo4j.SessionConfig{})
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		_, err := tx.Run(ctx, "MATCH (service:Service {name: $name }) DETACH DELETE service", map[string]any{
			"name": service,
		})

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		return "ok", nil
	})
	session.Close(ctx)

	if err != nil {
		return err
	}
	return nil
}

func GetConfig(service string, version int64, driver neo4j.DriverWithContext, ctx context.Context) models.Output {

	if version == -1 {
		session := (driver).NewSession(ctx, neo4j.SessionConfig{})
		ver, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

			records, err := tx.Run(ctx, "MATCH (service:Service { name: $name }) RETURN service.version ORDER BY (service.version) DESC LIMIT 1", map[string]any{
				"name": service,
			})

			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			record, err := records.Single(ctx)
			if err != nil {
				return nil, err
			}

			return record.Values[0], nil
		})
		session.Close(ctx)
		version = 0

		if err == nil {
			version = ver.(int64)
		}
	}
	session := (driver).NewSession(ctx, neo4j.SessionConfig{})
	output, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		records, err := tx.Run(ctx, "MATCH (service:Service { name: $name, version: $ver })-[e:Edge]->(item: Item) RETURN service.name, item.key, item.value", map[string]any{
			"name": service,
			"ver":  version,
		})

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		output := make(models.Output)

		for records.Next(ctx) {
			record := records.Record()
			output[record.Values[1].(string)] = record.Values[2].(string)
		}

		return output, nil

	})
	session.Close(ctx)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return output.(models.Output)
}

func UpdateConfig(config models.Config, driver neo4j.DriverWithContext, ctx context.Context) int64 {
	session := (driver).NewSession(ctx, neo4j.SessionConfig{})
	version, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		records, err := tx.Run(ctx, "MATCH (service:Service { name: $name }) RETURN service.version ORDER BY (service.version) DESC LIMIT 1", map[string]any{
			"name": config.Service,
		})

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		record, err := records.Single(ctx)
		if err != nil {
			return nil, err
		}

		return record.Values[0], nil
	})
	session.Close(ctx)

	if err != nil {
		fmt.Println(err)
		CreateConfig(config, driver, ctx)
		return -1
	}

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		_, err := tx.Run(ctx, "CREATE (service:Service { name: $name, version: $ver  })", map[string]any{
			"name": config.Service,
			"ver":  version.(int64) + 1,
		})

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for _, data := range config.Data {
			for key, value := range data {
				fmt.Println(fmt.Sprintf("MERGE (item:Item{key: \"%s\", value:\"%s\"}) MATCH (service:Service { name: \"%s\", version: %d}) CREATE (service)-[e1:Edge]->(item)", key, value, config.Service, version.(int64)+1))
				tx.Run(ctx, fmt.Sprintf("MATCH (service:Service { name: \"%s\", version: %d}) MERGE (item:Item{key: \"%s\", value:\"%s\"}) CREATE (service)-[e1:Edge]->(item)", config.Service, version.(int64)+1, key, value), map[string]any{})
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
			}
		}

		return "ok", nil
	})
	session.Close(ctx)

	if err != nil {
		return -1
	}
	return version.(int64) + 1
}
