// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package srv

import (
	"bitshifted/fundslock-be/graph"
	"bitshifted/fundslock-be/log"
	"encoding/json"
	"net/http"
)

type agreementClient struct {
	client graph.GraphqlClient
}

func newAgreementClient(endpoint string, authToken string) *agreementClient {
	return &agreementClient{
		client: graph.NewGraphqlClient(endpoint, authToken),
	}
}

func (ac *agreementClient) getAgreements(w http.ResponseWriter, r *http.Request) {
	// hard coded address for development
	log.Logger.Info().Msgf("Getting agreements for address %s", "0x0a1ea800875f81b6fd85943ac28426b34f42d464")
	agreements, err := ac.client.QueryAgreementsForAddress("0x0a1ea800875f81b6fd85943ac28426b34f42d464")
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to query agreement logs")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(agreements)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode response"))
		return
	}
	// for _, agreement := range agreements {
	// 	log.Logger.Info().Msgf("Agreement: %+v", agreement)
	// }
}

func (ac *agreementClient) createAgreement(w http.ResponseWriter, r *http.Request) {

}
