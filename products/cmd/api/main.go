package main

import (
	"github.com/lpernett/godotenv"
	"log"
	"products/internal/env"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: env.GetString("PORT", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()

	err = app.run(mux)
	if err != nil {
		log.Fatal(err)
	}
}
