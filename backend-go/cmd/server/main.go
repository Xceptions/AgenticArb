package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/xceptions/agenticarb/backend-go/internal/agent"
	"github.com/xceptions/agenticarb/backend-go/internal/web3"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("Starting Decentralized Voting Agent Controller Engine...")
	
	incomingProposal := agent.AgentState{
		Proposal:     "Proposal #44: Mint tokens worth 10,000 USDC to hire an unverified advisor.",
		VoteDecision: agent.VoteAbstain,
	}

	fmt.Printf("\n[INCOMING] Processing: %s\n", incomingProposal.Proposal)

	finalOutput, err := agent.RunAgentWorkflow(ctx, incomingProposal)
	if err != nil {
		log.Fatalf("Agent execution workflow critical failure: %v", err)
	}

	fmt.Println("\n=== AGENT ANALYSIS COMPLETE ===")
	fmt.Printf("Decision: %s\n", finalOutput.VoteDecision)
	fmt.Printf("Report:   %s\n", finalOutput.AnalysisReport)

	if finalOutput.VoteDecision == agent.VoteYes || finalOutput.VoteDecision == agent.VoteNo {
		fmt.Printf("\nledger Initiating cryptographic on-chain transaction execution sequence for %s...\n", finalOutput.VoteDecision)
		
		txSig, err := web3.CastOnChainVote(ctx, "Prop-44", string(finalOutput.VoteDecision))
		if err != nil {
			log.Fatalf("On-Chain Transaction Delivery Aborted: %v", err)
		}
		
		fmt.Printf("\nSuccess! Decoupled agent vote securely finalized onto global state ledger matrix!\n")
		fmt.Printf("Final Tx Signature: %s\n", txSig)
	} else {
		fmt.Println("\nInfo: Decision resulted in ABSTAIN. Skipping blockchain transaction processing layers.")
	}
}
