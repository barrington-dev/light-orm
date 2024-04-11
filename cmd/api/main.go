package main

import (
	"fmt"
	"light-orm/internal/server"
)

func main() {

	server := server.NewServer()

	fmt.Printf("starting server on port %v\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
