package clic_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/googollee/clic"
)

func Example_embedStruct() {
	// prepare env
	for _, key := range []string{"DEMO_VALUE"} {
		if err := os.Setenv(key, "value_from_env"); err != nil {
			log.Fatal("set env error:", err)
		}
	}

	// code starts
	type Inner struct {
		Value string `clic:"value,default,the value in the inner struct"`
	}

	type Config struct {
		Value string `clic:"value,default,the value in the outer struct"`
		Inner Inner  `clic:"inner"`
	}

	fset := flag.NewFlagSet("", flag.PanicOnError)
	set := clic.NewSet(fset)

	var cfg Config
	_ = set.RegisterValue("demo", &cfg)

	ctx := context.Background()
	_ = set.Parse(ctx, []string{"-demo.inner.value", "value_from_flag"})

	fmt.Println("Value:", cfg.Value)
	fmt.Println("Inner.Value:", cfg.Inner.Value)

	// Output:
	// Value: value_from_env
	// Inner.Value: value_from_flag
}
