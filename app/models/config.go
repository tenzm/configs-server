package models

type Config struct {
	Service string              `json:"service"`
	Data    []map[string]string `json:"data"`
}

type Output map[string]string
