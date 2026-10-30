package agent

import (
	"io"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xceptions/agenticarb/backend-go/config"
)

type OllamaEmbeddingFunction struct {
	BaseURL string
	Model   string
}

func (o *OllamaEmbeddingFunction) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	var embeddings [][]float32
	for _, text := range texts {
		emb, err := o.EmbedQuery(ctx, text)
		if err != nil {
			return nil, err
		}
		embeddings = append(embeddings, emb)
	}
	return embeddings, nil
}

func (o *OllamaEmbeddingFunction) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"model":  o.Model,
		"prompt": text,
	})
	req, err := http.NewRequestWithContext(ctx, "POST", o.BaseURL+"/api/embeddings", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Embedding, nil
}

func GetOrCreateCollectionID(ctx context.Context, name string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/collections", config.GlobalSettings.OllamaBaseURL)

	payload, _ := json.Marshal(map[string]interface{}{
		"name":          name,
		"get_or_create": true,
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err == nil {
		if id, ok := result["id"].(string); ok {
			return id, nil
		}
	}

	var arrayResult []map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &arrayResult); err == nil && len(arrayResult) > 0 {
		if id, ok := arrayResult[0]["id"].(string); ok {
			return id, nil
		}
	}

	return "", fmt.Errorf("failed to retrieve valid collection ID from Chroma server response")
}

func SeedHistoricalData(ctx context.Context, texts []string) error {
	collectionID, err := GetOrCreateCollectionID(ctx, config.GlobalSettings.ChromaCollectionName)
	if err != nil {
		return err
	}

	embFn := &OllamaEmbeddingFunction{
		BaseURL: config.GlobalSettings.OllamaBaseURL,
		Model:   config.GlobalSettings.EmbeddingModel,
	}

	embeddings, err := embFn.EmbedDocuments(ctx, texts)
	if err != nil {
		return fmt.Errorf("failed generating embeddings: %w", err)
	}

	ids := make([]string, len(texts))
	for i := range texts {
		ids[i] = fmt.Sprintf("doc_%d", i)
	}

	url := fmt.Sprintf("%s/api/v1/collections/%s/add", config.GlobalSettings.OllamaBaseURL, collectionID)
	payload, _ := json.Marshal(map[string]interface{}{
		"ids":        ids,
		"embeddings": embeddings,
		"documents":  texts,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("chroma server returned unexpected status code: %d", resp.StatusCode)
	}

	fmt.Printf("Database Seeded %d documents into %s using direct HTTP communications\n", len(texts), config.GlobalSettings.ChromaCollectionName)
	return nil
}
