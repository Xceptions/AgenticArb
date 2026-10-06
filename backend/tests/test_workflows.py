import json
from pathlib import Path
from unittest.mock import patch
import pytest
from agent.workflows import build_agent_graph

DATA_SET_PATH = Path(__file__).parent / "data" / "golden_dataset.json"
with open(DATA_SET_PATH, "r") as f:
    golden_cases = json.load(f)

test_ids = [case["id"] for case in golden_cases]

class TestWorkflowGolden:

    @pytest.mark.parametrize("case", golden_cases, ids=test_ids)
    @patch("agent.workflows.retrieve_proposal_context")
    def test_workflow_state_transitions(self, mock_retrieve_context, case):
        """Validates the full LangGraph execution flow using the golden dataset."""
        # 1. Mock the pipeline return value inside workflows
        mocked_context = "\n---\n".join(case["mocked_historical_docs"])
        mock_retrieve_context.return_value = mocked_context

        # 2. Compile the graph app
        app = build_agent_graph()

        # 3. Form initial Agent State
        initial_state = {
            "proposal": case["proposal"],
            "historical_context": "",
            "risk_score": 0,
            "analysis_report": "",
            "vote_decision": "ABSTAIN"
        }

        # 4. Invoke graph
        final_output = app.invoke(initial_state)

        # 5. Assertions using your Golden Ground Truth variables
        assert final_output["risk_score"] == case["expected_risk_score"]
        assert final_output["vote_decision"] == case["expected_vote_decision"]
        assert case["expected_report_keyword"] in final_output["analysis_report"]
        assert final_output["historical_context"] == mocked_context
