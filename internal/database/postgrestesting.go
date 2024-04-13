package database

import (
	"context"
	"database/sql"
	"fmt"
	"light-orm/internal/config"
	"log"
	"time"
)

type PostgresTestingService struct {
	db *sql.DB
}

var (
	dbTestInstance *PostgresTestingService
)

func NewPostgresTestingService() *PostgresTestingService {
	// Reuse Connection
	if dbTestInstance != nil {
		return dbTestInstance
	}
	conf := config.NewPostgresTestingConfig()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbTestInstance = &PostgresTestingService{
		db: db,
	}
	return dbTestInstance
}

func (s *PostgresTestingService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("testing db down: %v", err))
	}

	return map[string]string{
		"message": "Testing db is healthy",
	}
}

func (s *PostgresTestingService) Instance() *sql.DB {
	return s.db
}
