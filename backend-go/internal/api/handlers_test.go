package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/xceptions/agenticarb/backend-go/config"
)

func TestProposalHandlerTable(t *testing.T) {
	mockChromaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "/api/v1/collections") && !strings.Contains(r.URL.Path, "/query") {
			w.Write([]byte(`{"id": "mock_id_123", "name": "dao_governance"}`))
			return
		}
		if strings.Contains(r.URL.Path, "/query") {
			w.Write([]byte(`{"documents": "DAO Rule Core-1: The community treasury is empty."}`))
			return
		}
	}))
	defer mockChromaServer.Close()

	config.GlobalSettings.ChromaBaseURL = mockChromaServer.URL
	config.GlobalSettings.OllamaBaseURL = mockChromaServer.URL

	tests := []struct {
		name           string
		method         string
		payload        string
		expectedStatus int
		checkBody      string
	}{
		{
			name:           "Invalid Method GET Rejected",
			method:         http.MethodGet,
			payload:        "",
			expectedStatus: http.StatusMethodNotAllowed,
			checkBody:      "Method not allowed",
		},
		{
			name:           "Malformed Empty JSON Request Payload",
			method:         http.MethodPost,
			payload:        `{"proposal": ""}`,
			expectedStatus: http.StatusBadRequest,
			checkBody:      "Malformed or empty proposal body",
		},
		{
			name:           "Valid Proposal Processing Pipeline Execution",
			method:         http.MethodPost,
			payload:        `{"proposal": "Proposal #44: Mint tokens to hire team."}`,
			expectedStatus: http.StatusOK,
			checkBody:      "success",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(
				context.Background(),
				tc.method,
				"/api/proposals",
				bytes.NewBufferString(tc.payload),
			)
			if err != nil {
				t.Fatalf("Failed to initialize tracking request structure: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ProposalHandler)

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("%s Expected status code %d, but got %d", tc.name, tc.expectedStatus, rr.Code)
			}

			if !strings.Contains(rr.Body.String(), tc.checkBody) {
				t.Errorf("%s Expected body string to contain context '%s', but got: %s", tc.name, tc.checkBody, rr.Body.String())
			}

			if rr.Code == http.StatusOK {
				var response ResponsePayload
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Errorf("%s Failed to unmarshal successful handler response payload: %v", tc.name, err)
				}
				if response.Status != "success" {
					t.Errorf("%s Expected status key 'success', got '%s'", tc.name, response.Status)
				}
			}
		})
	}
}
