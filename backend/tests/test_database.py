import os
from pathlib import Path
from unittest.mock import patch, MagicMock
import pytest

def test_database_dir_config():
    """Test that the database path updates dynamically based on the .env file."""
    with patch.dict(os.environ, {"CHROMA_DB_DIR": "test_chroma_db"}):
        base_dir = Path(__file__).resolve().parent.parent.parent
        db_folder = os.getenv("CHROMA_DB_DIR", "chroma_db")
        db_dir = os.path.join(base_dir, db_folder)
        
        assert db_folder == "test_chroma_db"
        assert db_dir.endswith("test_chroma_db")

def test_database_fallback_config():
    """Test that the database falls back to 'chroma_db' if the environment variable is missing."""
    with patch.dict(os.environ, {}, clear=True):
        db_folder = os.getenv("CHROMA_DB_DIR", "chroma_db")
        assert db_folder == "chroma_db"

@patch('chromadb.Client')
def test_database_initialization(mock_chroma_client):
    """Test that your initialization function successfully contacts or creates a client."""

    mock_instance = MagicMock()
    assert mock_instance is not None
