// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package graph

import (
	"bitshifted/fundslock-be/log"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	agreementsQuery = `
query GetAgreementLogs($userAddr: String!){
 agreementLogs(
    where: {
      or: [
        { seller: $userAddr, },
        { buyer: $userAddr, }
      ]
    }
    orderBy: timestamp
    orderDirection: desc
  ) {
    id
    agreement_id
    seller
    buyer
    amount
    status
    timestamp
  }
}
`
)

type AgreementLog struct {
	Agreement_id string `json:"agreement_id"`
	Seller       string `json:"seller"`
	Buyer        string `json:"buyer"`
	Amount       string `json:"amount"`
	Status       int    `json:"status"`
	Timestamp    string `json:"timestamp"`
}

type AgreementsLogResponse struct {
	Data *LogsResponseData `json:"data"`
}

type LogsResponseData struct {
	AgreementLogs []AgreementLog `json:"agreementLogs"`
}

type QueryPayload struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operationName,omitempty"`
	Vars      map[string]interface{} `json:"variables,omitempty"`
}

type GraphqlClient interface {
	QueryAgreementsForAddress(string) (map[string][]AgreementLog, error)
}

type HttpGraphqlClient struct {
	GraphqlClient
	Endpoint  string
	AuthToken string
	Client    *http.Client
}

func (c *HttpGraphqlClient) QueryAgreementsForAddress(userAddress string) (map[string][]AgreementLog, error) {
	log.Logger.Debug().Msg("Running agreements query")
	query := QueryPayload{
		Query:     agreementsQuery,
		Operation: "Subgraphs",
		Vars: map[string]interface{}{
			"userAddr": userAddress,
		},
	}
	data, err := json.Marshal(query)
	if err != nil {
		log.Logger.Error().Msgf("Failed to marshal query string: %s", err)
		return nil, err
	}
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, c.Endpoint, bytes.NewReader(data))
	if err != nil {
		log.Logger.Error().Msgf("Failed to create graphQL request: %v", err)
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(request)
	if err != nil {
		log.Logger.Error().Msgf("Failed to get graphql response: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Error().Msgf("Failed to read response body: %v", err)
		return nil, err
	}
	log.Logger.Debug().Msgf("response body: %s", string(body))
	var result AgreementsLogResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Logger.Error().Msgf("Failed to unmarshal response body: %v", err)
		return nil, err
	}
	return convertResultToMap(result.Data.AgreementLogs), nil
}

func NewGraphqlClient(endpoint, authToken string) GraphqlClient {
	return &HttpGraphqlClient{
		Endpoint:  endpoint,
		AuthToken: authToken,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func convertResultToMap(logs []AgreementLog) map[string][]AgreementLog {
	out := make(map[string][]AgreementLog, 0)
	for _, l := range logs {
		lst, ok := out[l.Agreement_id]
		if ok {
			out[l.Agreement_id] = append(lst, l)
		} else {
			out[l.Agreement_id] = []AgreementLog{l}
		}
	}
	return out
}
