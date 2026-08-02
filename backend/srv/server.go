// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package srv

import (
	"bitshifted/fundslock-be/log"
	"bitshifted/fundslock-be/model"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func Start() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	// Sets 'Content-Type: application/json' on all responses
	router.Use(render.SetContentType(render.ContentTypeJSON))
	// CORS config
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))
	configLoader := model.NewConfigurationLoader()
	err := configLoader.Load()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to load configuration")
		return err
	}
	agreementClient := newAgreementClient(model.AppConfig.GraphUrl, model.AppConfig.GraphApiKey)

	router.Group(func(r chi.Router) {
		r.Get("/api/v1/agreements", agreementClient.getAgreements)
		r.Post("/api/v1/agreements", agreementClient.createAgreement)
	})

	router.Group(func(r chi.Router) {
		r.Get("/api/v1/auth/nonce", createNonce)
		r.Post("/api/v1/auth/verify", verifySIWEMessage)
	})

	server := http.Server{
		Addr:         ":3000",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return server.ListenAndServe()
}
