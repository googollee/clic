package database

import (
	"fmt"

	"github.com/googollee/clic/demo/log"
)

type Config struct {
	Addr string `clic:"address,,the address of the database"`
}

type DB struct {
	addr string
}

func New(cfg *Config) (*DB, error) {
	if cfg.Addr == "" {
		return nil, fmt.Errorf("invalid db address: %q", cfg.Addr)
	}
	return &DB{
		addr: cfg.Addr,
	}, nil
}

func (db *DB) Connect() {
	log.Info("connect db", "addr", db.addr)
}
