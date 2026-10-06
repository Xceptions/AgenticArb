import os
from solana.rpc.async_api import AsyncClient
from solders.keypair import Keypair

def get_or_create_agent_wallet() -> Keypair:
    """Generates a persistent local keypair for the agent so it keeps the same address."""
    key_path = os.path.join(os.path.dirname(__file__), "agent_keypair.json")
    
    if os.path.exists(key_path):
        with open(key_path, "r") as f:
            secret = bytes(eval(f.read()))
            return Keypair.from_bytes(secret)
    else:
        wallet = Keypair()
        with open(key_path, "w") as f:
            f.write(str(list(wallet.to_bytes())))
        print(f"[WEB3] New Agent Wallet created and saved to disk.")
        return wallet

async def ensure_funds(client: AsyncClient, wallet: Keypair):
    """Asynchronously checks the agent's balance and requests a free Devnet airdrop if empty."""
    balance_resp = await client.get_balance(wallet.pubkey())
    balance = balance_resp.value
    
    if balance < 10_000_000:
        print(f"[WEB3] Account empty. Requesting Devnet SOL for {wallet.pubkey()}...")
        try:
            airdrop_resp = await client.request_airdrop(wallet.pubkey(), 10_000_000)
            
            await client.confirm_transaction(airdrop_resp.value)
            print("[WEB3] Airdrop confirmed! Proceeding with vote...")
        except Exception as e:
            print(f"[WEB3 ERROR] Devnet airdrop failed or rate-limited: {e}")
            return "Failed: No funds available"