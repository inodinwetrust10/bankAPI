package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	// initialising the database by creating an account table and other required things
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	server := NewAPIServer(":3000", store)
	server.Run()
}
