package api

import (
	"encoding/json"
	"net/http"
	"github.com/xceptions/agenticarb/backend-go/internal/agent"
	"github.com/xceptions/agenticarb/backend-go/internal/web3"
)

type RequestPayload struct {
	Proposal string `json:"proposal"`
}

type ResponsePayload struct {
	Status         string             `json:"status"`
	Decision       agent.VoteDecision `json:"decision"`
	AnalysisReport string             `json:"analysis_report"`
	TxSignature    string             `json:"tx_signature,omitempty"`
}

func ProposalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var reqPayload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&reqPayload)
	if err != nil || reqPayload.Proposal == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Malformed or empty proposal body"}`))
		return
	}

	initialState := agent.AgentState{
		Proposal:     reqPayload.Proposal,
		VoteDecision: agent.VoteAbstain,
	}

	finalState, err := agent.RunAgentWorkflow(r.Context(), initialState)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var txSig string
	if finalState.VoteDecision == agent.VoteYes || finalState.VoteDecision == agent.VoteNo {
		txSig, err = web3.CastOnChainVote(r.Context(), "Web-Prop", string(finalState.VoteDecision))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "On-chain broadcasting aborted: " + err.Error()})
			return
		}
	}

	response := ResponsePayload{
		Status:         "success",
		Decision:       finalState.VoteDecision,
		AnalysisReport: finalState.AnalysisReport,
		TxSignature:    txSig,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
