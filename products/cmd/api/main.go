package main

import "log"

func main() {
	cfg := config{
		addr: ":8080",
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()

	err := app.run(mux)
	if err != nil {
		log.Fatal(err)
	}
}
