package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"light-orm/internal/database"
)

type Server struct {
	port int
	db   database.Service
}

type JSONResponseError[T any] struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Message T      `json:"message"`
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db: database.NewPostgresService(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func NewJSONResponseError[T any](code int, data T) *JSONResponseError[T] {
	return &JSONResponseError[T]{
		Code:    strconv.Itoa(code),
		Status:  http.StatusText(code),
		Message: data,
	}
}

func NewJSONResponseSuccess[T any](data T) *map[string]T {
	return &map[string]T{
		"data": data,
	}
}
