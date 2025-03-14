package test_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/models"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/server"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalcHandler(t *testing.T) {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)
	orchestrator := server.New(ctx)

	tests := []struct {
		name           string
		method         string
		requestBody    models.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "Valid POST request",
			method: http.MethodPost,
			requestBody: models.Request{
				Expression: "2 + 2",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name:           "Invalid method GET",
			method:         http.MethodGet,
			requestBody:    models.Request{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "You can use only POST method",
		},
		{
			name:   "Invalid expression",
			method: http.MethodPost,
			requestBody: models.Request{
				Expression: "2 + ",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "Bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest(tt.method, "/api/v1/calculate", bytes.NewBuffer(requestBody))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			handler := server.CalcHandler(orchestrator)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError != "" {
				var response models.BadResponse
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				var response models.ID
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.ID)
			}
		})
	}
}
