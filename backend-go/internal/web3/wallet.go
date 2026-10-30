package web3

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gagliardetto/solana-go"
)

func GetOrCreateAgentWallet() (solana.PrivateKey, error) {
	keyPath := filepath.Join("internal", "web3", "agent_keypair.json")

	if _, err := os.Stat(keyPath); err == nil {
		fileBytes, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}

		var secretKeyBytes []byte
		if err := json.Unmarshal(fileBytes, &secretKeyBytes); err != nil {
			return nil, err
		}

		return solana.PrivateKey(secretKeyBytes), nil
	}

	account, err := solana.NewRandomPrivateKey()
	if err != nil {
		return nil, err
	}

	secretKeyBytes := account[:]
	jsonBytes, err := json.Marshal(secretKeyBytes)
	if err != nil {
		return nil, err
	}

	_ = os.MkdirAll(filepath.Dir(keyPath), os.ModePerm)
	if err := os.WriteFile(keyPath, jsonBytes, 0644); err != nil {
		return nil, err
	}

	return account, nil
}
