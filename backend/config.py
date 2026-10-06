import os
from dotenv import load_dotenv
from pydantic_settings import BaseSettings

ENV_PATH = os.path.join(os.path.dirname(os.path.abspath(__file__)), ".env")
load_dotenv(dotenv_path=ENV_PATH)

class Settings(BaseSettings):
    OLLAMA_BASE_URL: str
    EMBEDDING_MODEL: str
    CHROMA_COLLECTION_NAME: str
    SOLANA_RPC_URL: str
    
    BASE_DIR: str = os.path.dirname(os.path.abspath(__file__))
    DB_FOLDER_NAME = os.getenv("CHROMA_DB_DIR", "chroma_db")
    DB_DIR: str = os.path.join(BASE_DIR, DB_FOLDER_NAME)

settings = Settings()
