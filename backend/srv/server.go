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
)

func Start() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	configLoader := model.NewConfigurationLoader()
	config, err := configLoader.Load()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to load configuration")
		return err
	}
	agreementClient := newAgreementClient(config.GraphUrl, config.GraphApiKey)

	router.Group(func(r chi.Router) {
		r.Get("/api/v1/agreements", agreementClient.getAgreements)
		r.Post("/api/v1/agreements", agreementClient.createAgreement)
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
