package repository

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
	"stride-wars-app/internal/ent"
)

func InitEnt() (*ent.Client, error) {
	client, err := ent.Open("sqlite3", "file:dev.db?_fk=1")
	if err != nil {
		return nil, err
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}
