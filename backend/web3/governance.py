from solana.rpc.async_api import AsyncClient
from solders.transaction import VersionedTransaction  # 1. Use VersionedTransaction instead of Transaction
from solders.message import MessageV0
from solders.instruction import Instruction
from solders.pubkey import Pubkey
from solders.keypair import Keypair
from web3.wallet import get_or_create_agent_wallet, ensure_funds
from config import settings

RPC_URL = settings.RPC_URL
MEMO_PROGRAM_ID = Pubkey.from_string(settings.MEMO_PROGRAM_ID)

async def cast_on_chain_vote(proposal_id: str, vote: str) -> str:
    """Asynchronously signs and broadcasts a vote transaction to the Solana Devnet ledger."""
    wallet = get_or_create_agent_wallet()
    
    async with AsyncClient(RPC_URL) as client:
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

        memo_data = f"AgentVote:{proposal_id}:{vote}".encode("utf-8")
        
        instruction = Instruction(
            program_id=MEMO_PROGRAM_ID,
            accounts=[],
            data=memo_data
        )
        
        blockhash_response = await client.get_latest_blockhash()
        latest_blockhash = blockhash_response.value.blockhash
        
        compiled_message = MessageV0.try_compile(
            payer=wallet.pubkey(),
            instructions=[instruction],
            address_lookup_table_accounts=[],
            recent_blockhash=latest_blockhash
        )
        
        tx = VersionedTransaction(compiled_message, [wallet])
        
        print(f"[WEB3] Broadcasting '{vote}' vote for {proposal_id} to Solana Devnet...")
        
        response = await client.send_transaction(tx)
        tx_signature = response.value
        
        await client.confirm_transaction(tx_signature)
        
        print("[SUCCESS] Vote transaction finalized on-chain!")
        explorer_url = f"https://solana.com{tx_signature}?cluster=devnet"
        print(f"Explorer URL: {explorer_url}")
        
        return str(tx_signature)
