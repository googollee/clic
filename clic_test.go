package clic_test

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/googollee/clic"
)

func Example() {
	// structs
	type Database struct {
		Driver string `clic:"driver,sqlite3,the driver of the database [sqlite3,mysql,postgres]"`
		URL    string `clic:"url,./database.sqlite,the url of the database"`
	}
	type Log struct {
		Level slog.Level `clic:"level,INFO,the level of the log [DEBUG,INFO,WARN,ERROR]"`
	}
	initLog := func(ctx context.Context, log *Log) error {
		fmt.Println("set log level:", log.Level)
		return nil
	}

	// args
	oldArgs := os.Args
	os.Args = []string{oldArgs[0], "-database.driver", "driver", "-database.url", "url", "sub_command"}
	defer func() {
		os.Args = oldArgs
	}()

	// main code
	ctx := context.Background()

	var db Database
	_ = clic.Register("database", &db)
	_ = clic.RegisterCallback("log", initLog)

	clic.Parse(ctx)

	fmt.Println("database:", db)
	fmt.Println("remain args:", flag.Args())

	// Output:
	// set log level: INFO
	// database: {driver url}
	// remain args: [sub_command]
}
