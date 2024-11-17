package main

import (
	"github.com/lpernett/godotenv"
	"go.uber.org/zap"
	"products/internal/store"
	"shared"
	"time"
)

//	@title			Products API
//	@version		1.0
//	@description	This is a sample server for a products API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Dawid Pereira
//	@contact.url	https://www.linkedin.com/in/pereiradawid/
//	@contact.email	pereiradawid@outlook.com

// @BasePath	/api/v1
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
		addr:    shared.GetString("PORT", ":8080"),
		env:     shared.GetString("ENV", "development"),
		version: shared.GetString("VERSION", "1.0"),
		rateLimiter: shared.Config{
			RequestPerTimeFrame: shared.GetInt("RATE_LIMIT_MAX_REQUESTS", 100),
			TimeFrame:           shared.GetDuration("RATE_LIMIT_WINDOW", 1*time.Minute),
			Enabled:             shared.GetBool("RATE_LIMIT_ENABLED", false),
		},
	}

	app := &application{
		config:      cfg,
		store:       storage,
		logger:      logger,
		rateLimiter: shared.NewFixedWindowRateLimiter(cfg.rateLimiter, logger),
	}

	mux := app.mount()

	err = app.run(mux)
	if err != nil {
		logger.Fatal(err)
	}
}
