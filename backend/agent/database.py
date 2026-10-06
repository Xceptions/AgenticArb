import os
from langchain_chroma import Chroma
from langchain_ollama import OllamaEmbeddings
from config import settings

def get_vector_store() -> Chroma:
    """Initializes and returns a persistent local vector database connection using parameters from config."""
    embeddings = OllamaEmbeddings(
        model=settings.EMBEDDING_MODEL,
        base_url=settings.OLLAMA_BASE_URL
    )
    
    vector_store = Chroma(
        persist_directory=settings.DB_DIR,
        embedding_function=embeddings,
        collection_name=settings.CHROMA_COLLECTION_NAME
    )
    return vector_store

def seed_historical_data(texts: list[str]):
    """Seeds the local database with historical DAO documents, rules, or past proposals."""
    db = get_vector_store()
    ids = [f"doc_{i}" for i in range(len(texts))]
    db.add_texts(texts=texts, ids=ids)
    print(f"[DATABASE] Successfully seeded {len(texts)} documents into '{settings.CHROMA_COLLECTION_NAME}'.")

if __name__ == "__main__":
    print(f"Initializing local database schema...")
    print(f"Targeting Ollama Node: {settings.OLLAMA_BASE_URL}")
    print(f"Embedding Engine Model: {settings.EMBEDDING_MODEL}")
    print(f"Collection Name: {settings.CHROMA_COLLECTION_NAME}")
    print(f"Storage Path: {settings.DB_DIR}")
    
    default_dao_documents = [
        "DAO Core Constitution (Rule 1): The community treasury must never fall below 5,000 SOL. Any proposal pushing it lower must be rejected automatically.",
        "DAO Core Constitution (Rule 2): Minting new tokens is completely frozen until the Q4 audits are finalized. No exceptions.",
        "DAO Hiring Guidelines: All external consultants and advisors must provide a minimum of 3 past approved references or security team checkups.",
        "Financial Report Q3: The current treasury status holds 12,000 SOL, 50,000 USDC, and 0 unallocated native tokens.",
        "Past Proposal #12 Summary: Approved hiring Dev Team Alpha for 500 SOL. Result: Successful deployment of the staking program.",
        "Past Proposal #40 Summary: Proposed spending 25,000 USDC on an unverified marketing vendor. Result: REJECTED due to lack of track record.",
    ]
    
    seed_historical_data(default_dao_documents)
    print(f"Database files successfully saved locally.")
