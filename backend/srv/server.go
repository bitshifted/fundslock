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
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

func Start() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	// Sets 'Content-Type: application/json' on all responses
	router.Use(render.SetContentType(render.ContentTypeJSON))
	// CORS config
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))
	configLoader := model.NewConfigurationLoader()
	err := configLoader.Load()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to load configuration")
		return err
	}
	agreementClient := newAgreementClient(model.AppConfig.GraphUrl, model.AppConfig.GraphApiKey)
	// initialze JWT authentication middlwware
	jwtInit()

	router.Group(func(r chi.Router) {
		r.Get("/api/v1/auth/nonce", createNonce)
		r.Post("/api/v1/auth/verify", verifySIWEMessage)
		r.Get("/api/v1/auth/refresh", refreshAccessToken)
	})

	// JWT protected paths
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/api/v1/users/session", getSession)
		r.Get("/api/v1/agreements", agreementClient.getAgreements)
	})
	// res, _ := agreementClient.client.QueryAgreementsForAddress("0x92c5fd33E29B31672Ba59D1109F8281d74fB838B")
	// log.Logger.Info().Msgf("result: %v", res)

	server := http.Server{
		Addr:         ":3000",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return server.ListenAndServe()
}
