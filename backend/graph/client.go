// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package graph

import (
	"bitshifted/fundslock-be/log"
	"context"
	"fmt"
	"net/http"

	"github.com/hasura/go-graphql-client"
)

type AgreementLogFilter struct {
	Seller string `json:"seller,omitempty"`
	Buyer  string `json:"buyer,omitempty"`
}

type AgreementLog struct {
	Agreement_id string `json:"agreement_id"`
	Seller       string `json:"seller"`
	Buyer        string `json:"buyer"`
	Amount       string `json:"amount"`
	Status       int    `json:"status"`
	Timestamp    string `json:"timestamp"`
}

type AgreementLogsQuery struct {
	//nolint:lll
	AgreementLogs []AgreementLog `graphql:"agreementLogs(where: { or: [{ seller: $userAddress }, { buyer: $userAddress }] }, orderBy: timestamp, orderDirection: desc, first: $first, skip: $skip)"`
}

type GraphqlClient interface {
	QueryAgreementsForAddress(string) ([]AgreementLog, error)
}

type HasuraGraphqlClient struct {
	GraphqlClient
	client *graphql.Client
}

func (g *HasuraGraphqlClient) QueryAgreementsForAddress(userAddress string) ([]AgreementLog, error) {
	var query AgreementLogsQuery
	pageSize := 10
	pageNumber := 0

	// Construct the 'where' argument payload
	variables := map[string]interface{}{
		"userAddress": userAddress,
		"first":       graphql.Int(pageSize),
		"skip":        graphql.Int(pageNumber * pageSize),
	}
	// Execute the query
	err := g.client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to query agreement logs")
		return nil, err
	}

	return query.AgreementLogs, nil
}

func NewGraphqlClient(endpoint, authToken string) GraphqlClient {
	client := graphql.NewClient(endpoint, http.DefaultClient).WithRequestModifier(func(r *http.Request) {
		r.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", authToken))
	})
	return &HasuraGraphqlClient{
		client: client,
	}
}
