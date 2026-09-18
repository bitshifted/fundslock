// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package srv

import (
	"bitshifted/fundslock-be/graph"
	"bitshifted/fundslock-be/log"
	"bitshifted/fundslock-be/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
)

type agreementClient struct {
	chain  string
	client graph.GraphqlClient
}

var graphClients = make(map[string]*agreementClient)

func initGraphqlClients() {
	graphConfig := model.AppConfig.GraphConfig
	for _, config := range graphConfig {
		// if _, exists := graphClients[config.Chain]; !exists {
		graphClients[config.Chain] = &agreementClient{
			chain:  config.Chain,
			client: graph.NewGraphqlClient(config.GraphUrl, config.ApiKey),
		}
		log.Logger.Info().Msgf("Initialized agreement client for chain %s", config.Chain)
		// }
	}
	// return &agreementClient{
	// 	client: graph.NewGraphqlClient(endpoint, authToken),
	// }
}

func getAgreements(w http.ResponseWriter, r *http.Request) {
	chain := r.URL.Query().Get("chain")
	if chain == "" {
		log.Logger.Error().Msg("Missing 'chain' query parameter")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Missing 'chain' query parameter"))
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to write response")
		}
		return
	}
	// extract user wallet address from auth token
	_, claims, _ := jwtauth.FromContext(r.Context())
	walletAddr := fmt.Sprintf("%v", claims["wallet_address"])

	log.Logger.Debug().Msgf("Getting agreements for address %s", walletAddr)
	agreements, err := graphClients[chain].client.QueryAgreementsForAddress(walletAddr, chain)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Failed to query agreement logs for chain %s", chain)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Failed to query agreement logs"))
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to write response")
		}
		return
	}

	err = json.NewEncoder(w).Encode(agreements)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		_, err1 := w.Write([]byte("Failed to encode response"))
		if err1 != nil {
			log.Logger.Error().Err(err1).Msg("Failed to write response")
		}
		return
	}
}
