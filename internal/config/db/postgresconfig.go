package db

import (
	"flag"
	"os"
)

type DBConfig struct {
	ConnectionString string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		ConnectionString: "host=localhost port=5432 user=postgres password=postgres dbname=short_url sslmode=disable",
	}
}

func (cfg *DBConfig) DBConfigInit() {
	defaultConnectionStr := cfg.ConnectionString

	if flag.Lookup("d") == nil {
		connectionStringFlag := flag.String("d", defaultConnectionStr, "строка подключения к БД")

		flag.Parse()

		cfg.ConnectionStringSet(connectionStringFlag)
	} else {
		cfg.ConnectionStringSet(&cfg.ConnectionString)
	}
}

func (cfg *DBConfig) ConnectionStringSet(connectionStringFlag *string) {
	switch {
	case os.Getenv("DATABASE_DSN") != "":
		cfg.ConnectionString = os.Getenv("DATABASE_DSN")
	case *connectionStringFlag != cfg.ConnectionString:
		cfg.ConnectionString = *connectionStringFlag
	}
}
