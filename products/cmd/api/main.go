package main

import (
	"github.com/lpernett/godotenv"
	"log"
	"products/internal/env"
	"products/internal/store"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	storage := store.NewStorage()

	cfg := config{
		addr: env.GetString("PORT", ":8080"),
	}

	app := &application{
		config: cfg,
		store:  storage,
	}

	mux := app.mount()

	err = app.run(mux)
	if err != nil {
		log.Fatal(err)
	}
}
