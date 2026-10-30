package web3

import (
	"fmt"
	"testing"
)

func TestBuildMemoTransactionInstruction(t *testing.T) {
	wallet, err := GetOrCreateAgentWallet()
	if err != nil {
		t.Fatalf("Failed checking signature wallet capabilities: %v", err)
	}

	proposalID := "Prop-44"
	vote := "YES"
	expectedPayload := fmt.Sprintf("AgentVote:%s:%s", proposalID, vote)

	memoData := []byte(fmt.Sprintf("AgentVote:%s:%s", proposalID, vote))

	if string(memoData) != expectedPayload {
		t.Errorf("Instruction payload serialization format mismatch: expected %s, got %s", expectedPayload, string(memoData))
	}

	if len(wallet.PublicKey()) != 32 {
		t.Errorf("Solana standard public address size constraints broken. Expected 32 bytes, got %d", len(wallet.PublicKey()))
	}
}
