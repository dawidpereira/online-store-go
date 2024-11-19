package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dawidpereira/online-store-go/products/docs"
	"github.com/dawidpereira/online-store-go/products/internal/store"
	"github.com/dawidpereira/online-store-go/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	addr        string
	env         string
	version     string
	rateLimiter shared.Config
}

type application struct {
	config      config
	store       store.Storage
	logger      *zap.SugaredLogger
	rateLimiter shared.RateLimiter
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(app.rateLimiter.RateLimiterMiddleware())

	//Test workflow
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL)))

		r.Route("/products", func(r chi.Router) {
			r.Get("/", app.listProductsHandler)
			r.Get("/{id}", app.getProductHandler)
			r.Post("/", app.createProductHandler)
			r.Put("/{id}", app.updateProductHandler)
			r.Delete("/{id}", app.deleteProductHandler)

		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = app.config.version
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost%s", app.config.addr)

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
