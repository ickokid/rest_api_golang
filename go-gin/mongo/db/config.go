package db

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Mongo_host     string
	Mongo_addrs    []string
	Mongo_database string
	Mongo_username string
	Mongo_password string
	Gin_mode       string
	Gin_port       string
	Key            string
}

func (c *Config) Read(configFile string) {
	if _, err := toml.DecodeFile(configFile, &c); err != nil {
		log.Fatal(err)
	}
}
