package clic_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/googollee/clic"
)

func Example_embedStruct() {
	// prepare env
	for _, key := range []string{"CLIC_DEMO_VALUE"} {
		if err := os.Setenv(key, "value_from_env"); err != nil {
			log.Fatal("set env error:", err)
		}
	}

	// prepare flags
	os.Args = append(os.Args, "-demo.inner.value", "value_from_flag")

	// code starts
	type Inner struct {
		Value string `clic:"value,default,the value in the inner struct"`
	}

	type Config struct {
		Value string `clic:"value,default,the value in the outer struct"`
		Inner Inner  `clic:"inner"`
	}

	loadConfig := clic.RegisterAndGet[Config]("demo")

	ctx := context.Background()
	clic.Init(ctx)

	cfg := loadConfig()

	fmt.Println("Value", cfg.Value)
	fmt.Println("Inner.Value:", cfg.Inner.Value)

	// Output:
	// Value value_from_env
	// Inner.Value: value_from_flag
}
