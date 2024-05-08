package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	Instance() *sql.DB
	Close() error
}
