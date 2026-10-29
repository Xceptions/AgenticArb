from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from agent.workflows import build_agent_graph

router = APIRouter(prefix="/v1/proposals", tags=["proposals"])

class ProposalRequest(BaseModel):
    proposal_text: str

class ProposalResponse(BaseModel):
    decision: str
    report: str
    risk_score: int

agent_app = build_agent_graph()

@router.post("/analyze", response_model=ProposalResponse)
def analyze_proposal_endpoint(payload: ProposalRequest):
    """Exposes the Agentic RAG evaluation loop over HTTP."""
    if not payload.proposal_text.strip():
        raise HTTPException(status_code=400, detail="Proposal text cannot be empty.")
    
    try:
        initial_state = {
            "proposal": payload.proposal_text,
            "historical_context": "",
            "risk_score": 0,
            "analysis_report": "",
            "vote_decision": "ABSTAIN"
        }
        
        result = agent_app.invoke(initial_state)
        
        return ProposalResponse(
            decision=result["vote_decision"],
            report=result["analysis_report"],
            risk_score=result["risk_score"]
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Agent workflow crashed: {str(e)}")
