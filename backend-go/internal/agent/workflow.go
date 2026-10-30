package agent

import (
	"context"
	"fmt"
	"strings"
)

type VoteDecision string

const (
	VoteYes     VoteDecision = "YES"
	VoteNo      VoteDecision = "NO"
	VoteAbstain VoteDecision = "ABSTAIN"
)

type AgentState struct {
	Proposal          string
	HistoricalContext string
	RiskScore         int
	AnalysisReport    string
	VoteDecision      VoteDecision
}

func ResearchNode(ctx context.Context, state *AgentState) error {
	fmt.Println("Agent Researcher executing RAG lookup...")

	contextText, err := RetrieveProposalContext(ctx, state.Proposal, 2)
	if err != nil {
		return err
	}
	state.HistoricalContext = contextText
	return nil
}

func AnalystNode(state *AgentState) {
	fmt.Println("Agent Analyst reviewing proposal against historical context...")

	proposalLower := strings.ToLower(state.Proposal)
	contextLower := strings.ToLower(state.HistoricalContext)

	if strings.Contains(proposalLower, "mint tokens") && strings.Contains(contextLower, "treasury empty") {
		state.RiskScore = 9
		state.AnalysisReport = "CRITICAL RISK: Proposal attempts to mint tokens but history indicates the treasury is depleted."
		state.VoteDecision = VoteNo
	} else {
		state.RiskScore = 2
		state.AnalysisReport = "Low Risk: Proposal aligns with historical governance parameters."
		state.VoteDecision = VoteYes
	}
}

func RunAgentWorkflow(ctx context.Context, initialState AgentState) (AgentState, error) {
	state := initialState

	if err := ResearchNode(ctx, &state); err != nil {
		return state, err
	}

	AnalystNode(&state)

	return state, nil
}
