package main

import (
	"github.com/lpernett/godotenv"
	"go.uber.org/zap"
	"products/internal/env"
	"products/internal/store"
)

//	@title			Products API
//	@version		1.0
//	@description	This is a sample server for a products API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Dawid Pereira
//	@contact.url	https://www.linkedin.com/in/pereiradawid/
//	@contact.email	pereiradawid@outlook.com

// @BasePath	/v1
func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer func(logger *zap.SugaredLogger) {
		_ = logger.Sync()
	}(logger)

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal(err)
	}

	storage := store.NewStorage()

	cfg := config{
		addr: env.GetString("PORT", ":8080"),
	}

	app := &application{
		config: cfg,
		store:  storage,
		logger: logger,
	}

	mux := app.mount()

	err = app.run(mux)
	if err != nil {
		logger.Fatal(err)
	}
}
