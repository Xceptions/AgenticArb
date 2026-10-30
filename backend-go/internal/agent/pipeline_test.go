package agent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/xceptions/agenticarb/backend-go/config"
)

func TestRetrieveProposalContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		if strings.Contains(r.URL.Path, "/api/v1/collections") && !strings.Contains(r.URL.Path, "/query") {
			w.Write([]byte(`{"id": "mock_id_123", "name": "dao_governance"}`))
			return
		}
		
		if strings.Contains(r.URL.Path, "/query") {
			w.Write([]byte(`{"documents": [["DAO Rule Core-1: Treasury is empty.", "Rule 2: Frozen."]]}`))
			return
		}
	}))
	defer server.Close()

	config.GlobalSettings.ChromaBaseURL = server.URL
	config.GlobalSettings.OllamaBaseURL = server.URL

	ctx := context.Background()
	proposal := "Proposal to mint new tokens"
	
	result, err := RetrieveProposalContext(ctx, proposal, 2)
	if err != nil {
		t.Fatalf("Expected no context lookup errors, got: %v", err)
	}

	if !strings.Contains(result, "Treasury is empty") {
		t.Errorf("Expected parsed string fragments to include 'Treasury is empty', got: %s", result)
	}
}
