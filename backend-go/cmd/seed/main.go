package main

import (
	"context"
	"fmt"
	"log"

	"github.com/xceptions/agenticarb/backend-go/config"
	"github.com/xceptions/agenticarb/backend-go/internal/agent"
)

func main() {
	ctx := context.Background()

	fmt.Println("Initializing standalone data loading script...")
	fmt.Printf("Targeting Ollama Node: %s\n", config.GlobalSettings.OllamaBaseURL)
	fmt.Printf("Embedding Engine Model: %s\n", config.GlobalSettings.EmbeddingModel)
	fmt.Printf("Collection Name: %s\n", config..GlobalSettings.ChromaCollectionName)

	defaultDaoDocuments := []string{
		"DAO Core Constitution (Rule 1): The community treasury must never fall below 5,000 SOL. Any proposal pushing it lower must be rejected automatically.",
		"DAO Core Constitution (Rule 2): Minting new tokens is completely frozen until the Q4 audits are finalized. No exceptions.",
		"DAO Hiring Guidelines: All external consultants and advisors must provide a minimum of 3 past approved references or security team checkups.",
		"Financial Report Q3: The current treasury status holds 12,000 SOL, 50,000 USDC, and 0 unallocated native tokens.",
		"Past Proposal #12 Summary: Approved hiring Dev Team Alpha for 500 SOL. Result: Successful deployment of the staking program.",
		"Past Proposal #40 Summary: Proposed spending 25,000 USDC on an unverified marketing vendor. Result: REJECTED due to lack of track record.",
	}

	err := agent.SeedHistoricalData(ctx, defaultDaoDocuments)
	if err != nil {
		log.Fatalf("Database seeding failed: %v", err)
	}

	fmt.Println("Database files successfully saved and updated locally.")
}