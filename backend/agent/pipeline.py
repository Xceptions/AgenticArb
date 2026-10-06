from agent.database import get_vector_store

def retrieve_proposal_context(proposal_text: str, k: int = 2) -> str:
    """
    Queries the vector database using the text of a new proposal 
    to retrieve relevant historical data.
    """
    db = get_vector_store()
    docs = db.similarity_search(query=proposal_text, k=k)
    context = "\n---\n".join([doc.page_content for doc in docs])
    print("----------CONTEXT----------")
    print(context)
    print("----------ENDCONTEXT-------")
    print(f"[RAG PIPELINE] Retrieved {len(docs)} relevant historical contexts.")
    return context
