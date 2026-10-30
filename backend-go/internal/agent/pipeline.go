package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/xceptions/agenticarb/backend-go/config"
)

func RetrieveProposalContext(ctx context.Context, proposalText string, k int32) (string, error) {
	collectionID, err := GetOrCreateCollectionID(ctx, config.GlobalSettings.ChromaCollectionName)
	if err != nil {
		return "", fmt.Errorf("pipeline context lookup failed to find collection: %w", err)
	}

	embFn := &OllamaEmbeddingFunction{
		BaseURL: config.GlobalSettings.OllamaBaseURL,
		Model:   config.GlobalSettings.EmbeddingModel,
	}
	queryVector, err := embFn.EmbedQuery(ctx, proposalText)
	if err != nil {
		return "", fmt.Errorf("pipeline context lookup failed to embed search text: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/collections/%s/query", config.GlobalSettings.ChromaBaseURL, collectionID)
	payload, _ := json.Marshal(map[string]interface{}{
		"query_embeddings": [][]float32{queryVector},
		"n_results":        k,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("chroma query endpoint returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var queryResponse struct {
		Documents [][]string `json:"documents"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&queryResponse); err != nil {
		return "", fmt.Errorf("failed to parse chroma match stream: %w", err)
	}

	if len(queryResponse.Documents) == 0 || len(queryResponse.Documents[0]) == 0 {
		return "", nil
	}

	joinedContext := strings.Join(queryResponse.Documents[0], "\n---\n")

	fmt.Println("Context:")
	fmt.Println(joinedContext)
	fmt.Println("-----------------")
	fmt.Printf("Pipeline Retrieved %d relevant historical contexts.\n", len(queryResponse.Documents[0]))

	return joinedContext, nil
}
