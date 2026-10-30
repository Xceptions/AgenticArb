package agent

import (
	"testing"
)

func TestRunAgentWorkflow(t *testing.T) {
	tests := []struct {
		name          string
		initialState  AgentState
		expectedVote  VoteDecision
		expectedRisk  int
	}{
		{
			name: "Critical Risk Mint Tokens When Empty",
			initialState: AgentState{
				Proposal:          "Proposal #99: Mint tokens to hire team",
				HistoricalContext: "DAO Rule Core-1: The community treasury is empty.",
			},
			expectedVote: VoteNo,
			expectedRisk: 9,
		},
		{
			name: "Low Risk Standard Compliance Alignments",
			initialState: AgentState{
				Proposal:          "Standard operational review setup",
				HistoricalContext: "DAO core hiring checklist updates",
			},
			expectedVote: VoteYes,
			expectedRisk: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := tc.initialState
			AnalystNode(&state)

			if state.VoteDecision != tc.expectedVote {
				t.Errorf("[%s] Expected decision verdict %s, got %s", tc.name, tc.expectedVote, state.VoteDecision)
			}

			if state.RiskScore != tc.expectedRisk {
				t.Errorf("%s Expected risk level grading score %d, got %d", tc.name, tc.expectedRisk, state.RiskScore)
			}
		})
	}
}
