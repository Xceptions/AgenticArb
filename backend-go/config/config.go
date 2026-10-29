package config

import "os"

type Settings struct {
	OllamaBaseURL        string
	EmbeddingModel       string
	ChromaCollectionName string
	DBDir                string
}

var GlobalSettings = Settings{
	OllamaBaseURL:        getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
	EmbeddingModel:       getEnv("EMBEDDING_MODEL", "nomic-embed-text"),
	ChromaCollectionName: getEnv("CHROMA_COLLECTION_NAME", "dao_governance"),
	DBDir:                getEnv("DB_DIR", "./chroma_data"),
}

// getEnv checks if an environment variable exists, otherwise returns a fallback string
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
