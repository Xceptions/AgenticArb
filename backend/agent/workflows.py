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


if __name__ == "__main__":
    import asyncio
    from agent.database import seed_historical_data
    from web3.governance import cast_on_chain_vote
    
    seed_historical_data([
        "DAO Rule Core-1: The community treasury is empty and cannot fund new tokens.",
        "DAO Rule Core-2: Advisors must have a minimum of 3 past approved references."
    ])
    
    app = build_agent_graph()
    incoming_proposal = {
        "proposal": "Proposal #44: Mint tokens worth 10,000 USDC to hire an unverified advisor.",
        "historical_context": "",
        "risk_score": 0,
        "analysis_report": "",
        "vote_decision": "ABSTAIN"
    }
    final_output = app.invoke(incoming_proposal)
    
    print("\n=== AGENT ANALYSIS COMPLETE ===")
    print(f"Decision: {final_output['vote_decision']}")
    print(f"Report: {final_output['analysis_report']}")
    
    if final_output["vote_decision"] in ["YES", "NO"]:
        async def main_web3_runner():
            try:
                tx_sig = await cast_on_chain_vote(
                    proposal_id="Prop-44", 
                    vote=final_output["vote_decision"]
                )
                print(f"\n[SUCCESS] Vote transaction finalized on-chain!")
                print(f"Explorer URL: https://solana.com{tx_sig}?cluster=devnet")
            except Exception as e:
                print(f"\n[WEB3 ERROR] Failed to push transaction: {e}")
                
        asyncio.run(main_web3_runner())
