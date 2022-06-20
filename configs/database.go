package configs

import (
	"fmt"
	"log"
	"os"
)

// DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

//get env variable
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

//set config database
func BuildDBConfig() *DBConfig {
	dbConfig := DBConfig{
		Host:     mustGetenv("DB_HOST"),
		Port:     mustGetenv("DB_PORT"),
		User:     mustGetenv("DB_USER"),
		Password: mustGetenv("DB_PASSWORD"),
		DBName:   mustGetenv("DB_NAME"),
	}
	return &dbConfig
}

//call database
func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}
