import json
from pathlib import Path
from unittest.mock import MagicMock, patch
import pytest
from agent.pipeline import retrieve_proposal_context

DATA_SET_PATH = Path(__file__).parent / "data" / "golden_dataset.json"
with open(DATA_SET_PATH, "r") as f:
    golden_cases = json.load(f)

test_ids = [case["id"] for case in golden_cases]

class TestPipelineGolden:

    @pytest.mark.parametrize("case", golden_cases, ids=test_ids)
    @patch("agent.pipeline.get_vector_store")
    def test_retrieve_proposal_context_golden(self, mock_get_store, case):
        """Validates that the pipeline accurately fetches and joins documents."""
        
        mock_docs = []
        for doc_text in case["mocked_historical_docs"]:
            mock_doc = MagicMock()
            mock_doc.page_content = doc_text
            mock_docs.append(mock_doc)

        mock_vector_store = MagicMock()
        mock_vector_store.similarity_search.return_value = mock_docs
        mock_get_store.return_value = mock_vector_store
        context_result = retrieve_proposal_context(case["proposal"], k=2)
        mock_vector_store.similarity_search.assert_called_once_with(
            query=case["proposal"], k=2
        )
        for doc_text in case["mocked_historical_docs"]:
            assert doc_text in context_result
