from unittest.mock import patch
import pytest
from fastapi.testclient import TestClient
from api.main import app

client = TestClient(app)

def test_health_check_endpoint():
    """Validates the system sanity check endpoint returns 200 OK."""
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "healthy", "service": "solana-governance-agent"}

@patch("agent.workflows.retrieve_proposal_context")
def test_analyze_proposal_high_risk(mock_retrieve):
    """Validates that a matching risk keyword results in a deterministic NO vote."""

    mock_retrieve.return_value = "CRITICAL DATA: The community treasury empty status is active."

    payload = {"proposal_text": "Proposal #99: Request to mint tokens immediately."}
    
    response = client.post("/v1/proposals/analyze", json=payload)
    
    assert response.status_code == 200
    data = response.json()
    
    assert data["decision"] == "NO"
    assert data["risk_score"] == 9

def test_analyze_proposal_empty_payload():
    payload = {"proposal_text": "   "}
    response = client.post("/v1/proposals/analyze", json=payload)
    assert response.status_code == 400
