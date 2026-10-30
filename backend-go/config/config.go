package config

import "os"

type Settings struct {
	ChromaBaseURL        string
	OllamaBaseURL        string
	EmbeddingModel       string
	ChromaCollectionName string
	DBDir                string
}

var GlobalSettings = Settings{
	ChromaBaseURL:        getEnv("CHROMA_BASE_URL", "http://localhost:8000"), // Chroma default port
	OllamaBaseURL:        getEnv("OLLAMA_BASE_URL", "http://localhost:11434"), // Ollama default port
	EmbeddingModel:       getEnv("EMBEDDING_MODEL", "nomic-embed-text"),
	ChromaCollectionName: getEnv("CHROMA_COLLECTION_NAME", "dao_governance"),
	DBDir:                getEnv("DB_DIR", "./chroma_data"),
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
