package clic_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"

	"github.com/googollee/clic"
	"github.com/googollee/clic/source"
)

func ExampleSet() {
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

	// set with sources
	fset := flag.NewFlagSet("", flag.ExitOnError)
	set := clic.NewSet(fset,
		// The order means srouce priority, flag > config file > env
		source.Flag(source.FlagSplitter(".")),
		source.File(source.FilePathFlag("config"), source.FileFormat(source.JSON{})),
		source.Env(source.EnvSplitter("_")),
	)

	// args
	args := []string{"-log.level", "WARN", "-database.driver", "driver", "-database.url", "url", "other_cmd"}

	// main code
	ctx := context.Background()

	var db Database
	if err := set.RegisterValue("database", &db); err != nil {
		log.Fatal("register error:", err)
	}

	if err := set.RegisterCallback("log", initLog); err != nil {
		log.Fatal("register error:", err)
	}

	if err := set.Parse(ctx, args); err != nil {
		log.Fatal("parse error:", err)
	}

	fmt.Println("database:", db)
	fmt.Println("remain args:", fset.Args())

	// Output:
	// set log level: WARN
	// database: {driver url}
	// remain args: [other_cmd]
}
