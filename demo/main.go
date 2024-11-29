package main

import (
	"context"

	"github.com/googollee/clic"
	"github.com/googollee/clic/demo/database"
	"github.com/googollee/clic/demo/log"
)

var getDB = clic.RegisterAndGet[database.Config]("database")

/*
Try:

 - go run . -h
    Show the usage
 - CLIC_DATABASE_ADDRESS=urienv_addr go run .
    Provide parameters through env variables.
 - CLIC_DATABASE_ADDRESS=urienv_addr go run . -config ./config.json
    Priority order: config file > env variables
 - go run . -config ./config.json -database.address "uri:new_addr"
    Priority order: flag arguments > config file
*/

func main() {
	ctx := context.Background()
	clic.Init(ctx)

	db, err := database.New(getDB())
	if err != nil {
		log.Error("create database fails", "error", err)
	}

	db.Connect()
}
