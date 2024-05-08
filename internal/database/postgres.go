package database

import (
	"context"
	"database/sql"
	"fmt"
	"light-orm/internal/config"
	"log"
	"time"
)

type PostgresService struct {
	db *sql.DB
}

var (
	dbInstance *PostgresService
)

func NewPostgresService() *PostgresService {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	conf := config.NewPostgresConfig()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &PostgresService{
		db: db,
	}
	return dbInstance
}

func (s *PostgresService) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

func (s *PostgresService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *PostgresService) Instance() *sql.DB {
	return s.db
}
