// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0

package srv

import (
	"bitshifted/fundslock-be/graph"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_GetAgreementsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := graph.NewMockGraphqlClient(ctrl)

	agreementCLient := &agreementClient{
		client: mockClient,
	}

	agLog := []graph.AgreementLog{
		{
			Agreement_id: "123",
			Seller:       "0x1231231323344",
			Buyer:        "0x3453453453453",
			Amount:       "1000",
			Status:       1,
			Timestamp:    "123455667",
		},
	}

	mockClient.EXPECT().QueryAgreementsForAddress(gomock.Any()).Return(agLog, nil)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/agreements", nil)
	rec := httptest.NewRecorder()

	agreementCLient.getAgreements(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "[{\"agreement_id\":\"123\",\"seller\":\"0x1231231323344\",\"buyer\":"+
		"\"0x3453453453453\",\"amount\":\"1000\",\"status\":1,\"timestamp\":\"123455667\"}]\n", string(body))
}

func Test_GetAgreementsQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := graph.NewMockGraphqlClient(ctrl)
	mockClient.EXPECT().QueryAgreementsForAddress(gomock.Any()).Return(nil, assert.AnError)

	agreementCLient := &agreementClient{
		client: mockClient,
	}
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/agreements", nil)
	rec := httptest.NewRecorder()

	agreementCLient.getAgreements(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
