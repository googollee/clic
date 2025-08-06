package clic_test

import (
	"context"
	"fmt"
	"os"

	"github.com/googollee/clic"
)

func Example() {
	// structs
	type Database struct {
		Driver string `clic:"driver,sqlite3,the driver of the database [sqlite3,mysql,postgres]"`
		URL    string `clic:"url,./database.sqlite,the url of the database"`
	}

	// args
	oldArgs := os.Args
	os.Args = []string{oldArgs[0], "-database.driver", "driver", "-database.url", "url"}
	defer func() {
		os.Args = oldArgs
	}()

	// main code
	ctx := context.Background()

	var db Database
	_ = clic.Register("database", &db)

	clic.Parse(ctx)

	fmt.Println("database:", db)

	// Output:
	// database: {driver url}
}
