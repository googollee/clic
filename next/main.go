package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/googollee/clic"
	"github.com/googollee/clic/source"
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
	type Config struct {
		Database Database `clic:"database"`
		Log      Log      `clic:"log"`
	}

	// main code
	ctx := context.Background()

	var config Config
	clic.Register("demo", &config)
	clic.Parse(ctx, os.Args[1:])
}

func Example_set() {
	// structs
	type Database struct {
		Driver string `clic:"driver,sqlite3,the driver of the database [sqlite3,mysql,postgres]"`
		URL    string `clic:"url,./database.sqlite,the url of the database"`
	}
	type Log struct {
		Level slog.Level `clic:"level,INFO,the level of the log [DEBUG,INFO,WARN,ERROR]"`
	}
	initLog := func(ctx context.Context, log *Log) error {
		slog.SetLogLoggerLevel(log.Level)
		return nil
	}

	// set with sources
	fset := flag.NewFlagSet("", flag.ExitOnError)
	flagSource := source.Flag(source.FlagWithSet(fset), source.FlagWithSplitter("."))
	set := clic.NewSet("demo",
		// The order means srouce priority, flag > config file > env
		flagSource,
		source.File(source.FileFlag(fset, "config", "the config file path"), source.FileFormat(source.JSON{})),
		source.Env(source.EnvSplitter("_")),
	)

	// main code
	ctx := context.Background()

	var db Database
	set.RegisterValue("database", &db)
	set.RegisterCallback("log", initLog)
	set.Parse(ctx, os.Args[1:])

	// remain
	fmt.Println("remain args:", flagSource.Args())
}
