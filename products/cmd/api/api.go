package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"products/internal/store"
	"time"
)

type config struct {
	addr    string
	env     string
	version string
}

type application struct {
	config config
	store  store.Storage
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

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
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server is listening on %s", app.config.addr)

	return srv.ListenAndServe()
}
