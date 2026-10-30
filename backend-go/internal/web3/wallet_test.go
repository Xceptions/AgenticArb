package web3

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetOrCreateAgentWallet(t *testing.T) {
	testKeyPath := filepath.Join("internal", "web3", "agent_keypair.json")
	
	_ = os.Remove(testKeyPath)
	defer os.Remove(testKeyPath)

	wallet1, err := GetOrCreateAgentWallet()
	if err != nil {
		t.Fatalf("Failed to initialize randomized agent keypair instance: %v", err)
	}

	if len(wallet1) == 0 {
		t.Fatal("Generated cryptographic private key slice array is empty")
	}

	wallet2, err := GetOrCreateAgentWallet()
	if err != nil {
		t.Fatalf("Failed reloading cached local signature keys from disk: %v", err)
	}

	if !wallet1.PublicKey().Equals(wallet2.PublicKey()) {
		t.Error("Cryptographic alignment check failed: Keypair reloaded from disk does not match original generated keys")
	}
}
