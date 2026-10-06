# Autonomous AI Governance Agent

An autonomous, stateful **AI Agent for decentralized governance**, engineered to intelligently evaluate DAO proposals and execute automated, on-chain voting. 

The agent utilizes **Retrieval-Augmented Generation (RAG)** to parse incoming proposals against historical DAO constitutions, rules, and past decisions, calculates a dynamic risk assessment, and securely signs and broadcasts voting transactions to the **Solana Devnet ledger**.

---

## 🚀 Key Features

*   **Stateful AI Workflows:** Orchestrated via **LangGraph**, utilizing a dedicated multi-node graph (`Researcher` -> `Analyst`) to handle complex cognitive reasoning loops.
*   **Context-Aware Analysis (RAG):** Integrates **LangChain Chroma** and **Ollama Embeddings** to retrieve deep semantic context from past historical DAO data and local constitutional documents.
*   **Automated Risk Grading:** Evaluates proposal logic to dynamically output risk scores, analytical logs, and explicit governance recommendations (`YES`, `NO`, `ABSTAIN`).
*   **On-Chain Vote Execution:** Automatically drafts, signs, and pushes modern Solana `VersionedTransaction` payloads to the **Solana Memo Program** to permanently store vote records.
*   **Autonomous Key & Gas Management:** Securely generates and caches local file-system keypairs and auto-requests airdrops from the Solana Devnet faucet if gas balances drop below minimum thresholds.

---

## 📂 Project Architecture

```text
backend/
├── agent/
│   ├── database.py       # Local Chroma DB initialization & document seeding
│   ├── pipeline.py       # RAG logic for semantic document retrieval
│   └── workflows.py      # LangGraph state machine & evaluation workflow runner
├── web3/
│   ├── wallet.py         # Persistent local Solana Keypair storage & gas checks
│   └── governance.py     # Asynchronous versioned transaction signer & broadcaster
├── tests/                # Test suites tracking module behavior
├── config.py             # Pydantic-Settings environment validation setup
├── requirements.txt      # Requirement files for virtual env setup
└── .env                  # Configuration variables (git ignored)
```

---

## 🛠️ Tech Stack

*   **Frameworks:** LangGraph, LangChain (Chroma, Ollama integrations)
*   **Vector Storage:** Chroma DB (Local Persistent Storage)
*   **Inference Engine:** Ollama (Local LLM and Embeddings Engine)
*   **Blockchain Integration:** Solana Python SDK (`solana`, `solders`)
*   **Environment & Validation:** Pydantic Settings, `python-dotenv`

---

## ⚙️ Installation & Setup

### 1. Prerequisites
Ensure you have the following installed locally:
*   [Python 3.10+](https://python.org)
*   [Ollama](https://ollama.com) (Make sure it is running locally with your chosen embedding/LLM models pulled down)

### 2. Clone and Install Dependencies
Navigate to your project root directory and install required Python libraries:
```bash
pip install langgraph langchain-chroma langchain-ollama solana solders pydantic-settings python-dotenv
```

### 3. Environment Configuration
Create a `.env` file inside the `backend/` directory:
```env
OLLAMA_BASE_URL=http://localhost:11434
EMBEDDING_MODEL=nomic-embed-text
CHROMA_COLLECTION_NAME=dao_governance_rules
CHROMA_DB_DIR=chroma_db
SOLANA_RPC_URL=https://api.devnet.solana.com
RPC_URL=https://api.devnet.solana.com
MEMO_PROGRAM_ID=Mem1111111111111111111111111111111111111111
```

---

## 🔌 Usage Guide

### 🧬 Step 1: Initialize & Seed the Vector Database
Populate your local Chroma vector store with constitutional rules, financial histories, and guidelines:
```bash
python -m agent.database
```

### 🤖 Step 2: Execute the AI Agent Workflow
To pass a new test proposal through the multi-node LangGraph analysis pipeline and run the automatic on-chain execution:
```bash
python -m agent.workflows
```

*   **How it works:** The workflow automatically loads the internal project loop, triggers an internal RAG lookup against your Chroma DB vector store, evaluates risk metrics, generates a decision, handles keypair checking, and submits the finalized transaction to Solana Devnet.

---

## 🧪 Testing
To verify functionality across both core modules, run the unit test commands via `pytest` (if configured):
```bash
pytest tests/
```
