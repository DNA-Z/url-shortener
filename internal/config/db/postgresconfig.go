package db

import (
	"flag"
	"os"
)

type DBConfig struct {
	ConnectionString string
}

func DBConfigInit() *DBConfig {
	cfg := &DBConfig{}

	// Значение по умолчанию (опционально, можно убрать)
	defaultDSN := "host=localhost port=5432 user=postgres password=lightning dbname=short_url sslmode=disable"

	// Определяем флаг -d
	var dsnFlag string
	flag.StringVar(&dsnFlag, "d", defaultDSN, "строка подключения к БД")
	flag.Parse()

	// Приоритет: ENV > флаг -d
	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		cfg.ConnectionString = envDSN
	} else {
		cfg.ConnectionString = dsnFlag
	}

	return cfg
}

func (cfg *DBConfig) GetConnectionString() string {
	return cfg.ConnectionString
}
