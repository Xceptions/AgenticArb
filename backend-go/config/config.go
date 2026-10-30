package config

import "os"

type Settings struct {
	ChromaBaseURL        string
	OllamaBaseURL        string
	EmbeddingModel       string
	ChromaCollectionName string
	DBDir                string
	RPCURL               string
	MemoProgramID        string
}

var GlobalSettings = Settings{
	ChromaBaseURL:        getEnv("CHROMA_BASE_URL", "http://localhost:8000"),
	OllamaBaseURL:        getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
	EmbeddingModel:       getEnv("EMBEDDING_MODEL", "nomic-embed-text"),
	ChromaCollectionName: getEnv("CHROMA_COLLECTION_NAME", "dao_governance"),
	DBDir:                getEnv("DB_DIR", "./chroma_data"),
	RPCURL:               getEnv("RPC_URL", "https://api.devnet.solana.com"),
	MemoProgramID:        getEnv("MEMO_PROGRAM_ID", "Memomet1vbglC67g1A9wA6saChn6N5gB4Yw9B9"),
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
