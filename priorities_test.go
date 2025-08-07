package clic_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/googollee/clic"
)

func ExampleSet_sourcePriorities() {
	// prepare env
	for _, key := range []string{"DEMO_VALUE_FLAG", "DEMO_VALUE_ENV", "DEMO_VALUE_FILE"} {
		if err := os.Setenv(key, "value_from_env"); err != nil {
			log.Fatal("set env error:", err)
		}
	}

	// prepare config file
	cfgFile, err := os.CreateTemp("", "config_*.json")
	if err != nil {
		log.Fatal("create temp file error:", err)
	}
	defer os.Remove(cfgFile.Name())

	if _, err := cfgFile.WriteString(`{
		"demo": {
			"value_flag": "value_from_file",
			"value_file": "value_from_file"
		}
	}`); err != nil {
		log.Fatal("write temp file error:", err)
	}

	if err := cfgFile.Close(); err != nil {
		log.Fatal("close temp file error:", err)
	}

	// prepare args
	args := []string{"-config", cfgFile.Name(), "-demo.value_flag", "value_from_flag"}

	// code starts
	type Config struct {
		ValueFlag    string `clic:"value_flag,default,a test value in flag"`
		ValueEnv     string `clic:"value_env,default,a test value in env"`
		ValueFile    string `clic:"value_file,default,a test value in config file"`
		ValueDefault string `clic:"value_default,default,a test value by default"`
	}
	var cfg Config

	fset := flag.NewFlagSet("", flag.PanicOnError)
	set := clic.NewSet(fset, clic.DefaultSources...)

	set.RegisterValue("demo", &cfg)

	ctx := context.Background()
	if err := set.Parse(ctx, args); err != nil {
		log.Fatal("parse error:", err)
	}

	fmt.Println("ValueFlag:", cfg.ValueFlag)
	fmt.Println("ValueEnv:", cfg.ValueEnv)
	fmt.Println("ValueFile:", cfg.ValueFile)
	fmt.Println("ValueDefault:", cfg.ValueDefault)

	// Output:
	// ValueFlag: value_from_flag
	// ValueEnv: value_from_env
	// ValueFile: value_from_file
	// ValueDefault: default
}
