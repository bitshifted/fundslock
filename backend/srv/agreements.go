// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package srv

import (
	"bitshifted/fundslock-be/graph"
	"bitshifted/fundslock-be/log"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
)

type agreementClient struct {
	client graph.GraphqlClient
}

func newAgreementClient(endpoint, authToken string) *agreementClient {
	return &agreementClient{
		client: graph.NewGraphqlClient(endpoint, authToken),
	}
}

func (ac *agreementClient) getAgreements(w http.ResponseWriter, r *http.Request) {
	// extract user wallet address from auth token
	_, claims, _ := jwtauth.FromContext(r.Context())
	walletAddr := fmt.Sprintf("%v", claims["wallet_address"])

	log.Logger.Debug().Msgf("Getting agreements for address %s", walletAddr)
	agreements, err := ac.client.QueryAgreementsForAddress(walletAddr)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to query agreement logs")
		w.WriteHeader(http.StatusInternalServerError)
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
