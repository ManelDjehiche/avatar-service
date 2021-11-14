package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/emplorium/auth-service/env"
	_ "github.com/lib/pq"
)

func StartDatabase() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		env.Settings.Database.PSQLConfig.Host,
		env.Settings.Database.PSQLConfig.Port,
		env.Settings.Database.PSQLConfig.Username,
		env.Settings.Database.PSQLConfig.Password,
		env.Settings.Database.PSQLConfig.DbName)

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Fatalf("Failed to connect to postgres database: %v", err)
	}
	return db, nil
}
