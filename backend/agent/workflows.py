from typing import TypedDict, Literal
from langgraph.graph import StateGraph, END
from agent.pipeline import retrieve_proposal_context


class AgentState(TypedDict):
    proposal: str
    historical_context: str
    risk_score: int
    analysis_report: str
    vote_decision: Literal["YES", "NO", "ABSTAIN"]


def research_node(state: AgentState) -> dict:
    print("[AGENT] Researcher executing RAG lookup...")
    context = retrieve_proposal_context(state["proposal"])
    return {"historical_context": context}


def analyst_node(state: AgentState) -> dict:
    print("[AGENT] Analyst reviewing proposal against historical context...")
    context = state["historical_context"]
    proposal = state["proposal"]
    
    if "mint tokens" in proposal.lower() and "treasury empty" in context.lower():
        risk_score = 9
        report = "CRITICAL RISK: Proposal attempts to mint tokens but history indicates the treasury is depleted."
        decision = "NO"
    else:
        risk_score = 2
        report = "Low Risk: Proposal aligns with historical governance parameters."
        decision = "YES"
        
    return {
        "risk_score": risk_score,
        "analysis_report": report,
        "vote_decision": decision
    }

def build_agent_graph():
    workflow = StateGraph(AgentState)
    workflow.add_node("researcher", research_node)
    workflow.add_node("analyst", analyst_node)
    workflow.set_entry_point("researcher")
    workflow.add_edge("researcher", "analyst")
    workflow.add_edge("analyst", END)
    
    return workflow.compile()

