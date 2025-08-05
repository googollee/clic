package clic_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/googollee/clic"
)

func ExampleInit_valuePriorities() {
	// prepare env
	for _, key := range []string{"CLIC_DEMO_VALUE_FLAG", "CLIC_DEMO_VALUE_ENV", "CLIC_DEMO_VALUE_FILE"} {
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

	// prepare flags
	os.Args = append(os.Args, "-config", cfgFile.Name(), "-demo.value_flag", "value_from_flag")

	// code starts
	type Config struct {
		ValueFlag    string `clic:"value_flag,default,a test value in flag"`
		ValueEnv     string `clic:"value_env,default,a test value in env"`
		ValueFile    string `clic:"value_file,default,a test value in config file"`
		ValueDefault string `clic:"value_default,default,a test value by default"`
	}

	loadConfig := clic.RegisterAndGet[Config]("demo")

	ctx := context.Background()
	clic.Init(ctx)

	cfg := loadConfig()

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
