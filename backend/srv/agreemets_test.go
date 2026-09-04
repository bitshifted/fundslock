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

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_GetAgreementsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := graph.NewMockGraphqlClient(ctrl)

	agreementCLient := &agreementClient{
		client: mockClient,
	}

	items := []graph.AgreementResponseItem{
		{
			AgreementId: 123,
			Seller:      "0x1231231323344",
			Buyer:       "0x3453453453453",
			Amount:      1.2,
			Status:      1,
			StatusChanges: []graph.AgreementStatusChange{
				{
					Status:    1,
					Timestamp: "2023-01-01T12:00:00Z",
				},
			},
		},
	}

	mockClient.EXPECT().QueryAgreementsForAddress(gomock.Any()).Return(items, nil)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/agreements", nil)
	rec := httptest.NewRecorder()

	agreementCLient.getAgreements(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, "[{\"agreementId\":123,\"seller\":\"0x1231231323344\",\"buyer\":\"0x3453453453453\""+
		",\"amount\":1.2,\"status\":1,\"statusChanges\":[{\"status\":1,\"timestamp\":\"2023-01-01T12:00:00Z\"}]}]\n", string(body))
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
