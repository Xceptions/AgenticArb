# Autonomous AI Governance Agent

An autonomous, stateful **AI Agent for decentralized governance**, engineered to intelligently evaluate DAO proposals and execute automated, on-chain voting. 

The agent utilizes **Retrieval-Augmented Generation (RAG)** to parse incoming proposals against historical DAO constitutions, rules, and past decisions, calculates a dynamic risk assessment, and securely signs and broadcasts voting transactions to the **Solana Devnet ledger**.

---

## 🚀 Key Features

*   **Dual-Engine Backend Architecture:** Implemented in both **Python** (for fast data prototyping and graph design) and **Go** (for enterprise-grade performance, strict typing, and high-throughput execution).
*   **Stateful AI Workflows:** Orchestrated via sequential multi-node layouts (`Researcher` -> `Analyst`) to safely evaluate complex multi-faceted cognitive reasoning loops.
*   **Decoupled Context-Aware RAG:** Connects cleanly to raw REST collection matrices to perform vector searches on **Chroma DB** and **Ollama Embeddings** without heavy client library dependencies.
*   **Automated Risk Grading:** Evaluates proposal logic dynamically to output explicit governance decisions (`YES`, `NO`, `ABSTAIN`) with associated risk indexes and analytical tracking logs.
*   **Native Crypto Ledger Execution:** Leverages native blockchain bindings to build instructions, sign with secure private keys, and broadcast memo tracking metrics directly onto the **Solana Devnet ledger**.
*   **Autonomous Key & Gas Management:** Automatically initializes and manages persistent local file-system keypairs and auto-requests airdrops from the Solana faucet if gas balances drop below execution minimums.

---

## 📂 Project Architecture

```text
.
├── backend-py/           # Python Prototyping Stack
│   ├── agent/
│   │   ├── database.py   # Local Chroma DB initialization & document seeding
│   │   ├── pipeline.py   # RAG logic for semantic document retrieval
│   │   └── workflows.py  # LangGraph state machine & evaluation workflow runner
│   ├── web3/
│   │   ├── wallet.py     # Persistent local Solana Keypair storage & gas checks
│   │   └── governance.py # Asynchronous versioned transaction signer & broadcaster
│   ├── config.py         # Pydantic-Settings environment validation setup
│   └── requirements.txt  # Requirement files for virtual env setup
│
├── backend-go/           # Production Go Stack
│   ├── cmd/
│   │   ├── seed/
│   │   │   └── main.go   # Standalone administrative DB seeding CLI tool
│   │   └── server/
│   │       └── main.go   # Main engine runner executing agent HTTP server & workflows
│   ├── config/
│   │   └── config.go     # Centralized zero-dependency configuration parser
│   ├── internal/
│   │   ├── agent/
│   │   │   ├── database.go      # Native JSON REST client utilities for Chroma DB
│   │   │   ├── pipeline.go      # Vector-search query execution pipeline
│   │   │   ├── pipeline_test.go # Mock HTTP server pipeline vector tests
│   │   │   ├── workflow.go      # Strong-typed sequential state orchestration engine
│   │   │   └── workflow_test.go # Table-driven test suite for analytical grading
│   │   ├── api/
│   │   │   ├── handlers.go      # HTTP router endpoint logic evaluating payload streams
│   │   │   └── handlers_test.go # In-memory API request validation tests
│   │   └── web3/
│   │       ├── wallet.go        # Keypair storage, loading, and local caching
│   │       ├── wallet_test.go   # Key parsing & file-system serialization tests
│   │       ├── governance.go    # Modern gagliardetto/solana-go ledger transaction manager
│   │       └── governance_test.go # Crypto serialization structure assertions
│   └── go.mod            # Go module dependency management graph
│
└── webUI/                # Frontend User Interface
```

---

## 🛠️ Tech Stack

*   **Languages:** Go (Golang 1.24+), Python (3.10+)
*   **Frameworks:** LangGraph, LangChain Ecosystem
*   **Vector Storage:** Chroma DB (Run containerized via Docker on port `8000`)
*   **Inference Engine:** Ollama (Local LLM and Embeddings Engine on port `11434`)
*   **Blockchain Integration:** `gagliardetto/solana-go` (Go), `solana`/`solders` (Python)
*   **Environment & Validation:** Go standard library (`os`), Pydantic Settings (Python)

---

## ⚙️ Installation & Setup

### 1. Prerequisites
Ensure you have the following installed locally:
*   [Go 1.24+](https://go.dev)
*   [Python 3.10+](https://python.org)
*   [Docker Desktop](https://docker.com)
*   [Ollama](https://ollama.com) (Ensure it is running locally with `nomic-embed-text` pulled down)

### 2. Infrastructure Setup
Spin up your local containerized Chroma DB instance and check your local inference engines:
```bash
# Start Chroma DB vector server in the background
docker run -d -p 8000:8000 chromadb/chroma

# Verify your Ollama inference engine node is awake
curl http://localhost:11434
```

### 3. Environment Configuration
Create a `.env` file at the root level of your projects:
```env
CHROMA_BASE_URL=http://localhost:8000
OLLAMA_BASE_URL=http://localhost:11434
EMBEDDING_MODEL=nomic-embed-text
CHROMA_COLLECTION_NAME=dao_governance
RPC_URL=https://solana.com
MEMO_PROGRAM_ID=Memomet1vbglC67g1A9wA6saChn6N5gB4Yw9B9
```

---

## 🔌 Usage Guide

### 🐹 Running the Go Production Stack

#### Step 1: Seed the Vector Database
Populate your containerized Chroma vector engine with constitutional guidelines and historical parameters using the dedicated toolchain script:
```bash
cd backend-go
go run ./cmd/seed
```

#### Step 2: Run the Production Agent Server
Execute the main orchestration workflow loop to launch the persistent backend web server, ready to receive incoming payloads from the UI on port `:8080`:
```bash
go run ./cmd/server
```

---

### 🐍 Running the Python Prototyping Stack

#### Step 1: Install Python Libraries
```bash
cd backend-py
pip install langgraph langchain-chroma langchain-ollama solana solders pydantic-settings python-dotenv
```

#### Step 2: Initialize & Run Workflows
```bash
# Seed records
python -m agent.database

# Execute LangGraph machine pipelines
python -m agent.workflows
```

---

## 🧪 Testing

The Go production backend includes an exhaustive testing suite incorporating **table-driven verification matrices**, **in-memory HTTP recorders**, and **mock network servers** to guarantee robust code execution without testing side-effects.

To run the complete Go verification suite, navigate to the `backend-go` folder and use the native testing toolchain:
```bash
cd backend-go

# Run all packages and print verbose evaluation outputs
go test -v ./...
```

*   `internal/agent/pipeline_test.go`: Captures outbound Chroma context fetches using an isolated mock server hook.
*   `internal/agent/workflow_test.go`: Validates multiple risk-grading edge-cases simultaneously via matrix inputs.
*   `internal/api/handlers_test.go`: Verifies CORS handling and routing assertions using in-memory request recorders (`httptest`).
*   `internal/web3/wallet_test.go`: Asserts disk-caching properties to ensure local key generation loops do not break or mutate wallet storage on reuse.
