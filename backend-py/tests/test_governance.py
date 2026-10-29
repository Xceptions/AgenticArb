import pytest
from unittest.mock import AsyncMock, MagicMock, patch
from solders.hash import Hash
from solders.keypair import Keypair

from web3.governance import cast_on_chain_vote


@pytest.mark.asyncio
@patch("web3.governance.get_or_create_agent_wallet")
@patch("web3.governance.AsyncClient")
async def test_cast_on_chain_vote_successful(mock_async_client_class, mock_get_wallet):
    """Verifies a full successful vote compiling, signing, and broadcasting sequence."""
    
    mock_wallet = Keypair()
    mock_get_wallet.return_value = mock_wallet

    mock_client_instance = AsyncMock()
    
    # Balance mock (sufficient funds: 20M lamports)
    mock_balance_resp = MagicMock()
    mock_balance_resp.value = 20_000_000
    mock_client_instance.get_balance.return_value = mock_balance_resp
    
    # Blockhash mock - using Hash instead of Pubkey
    mock_blockhash_resp = MagicMock()
    mock_blockhash_resp.value.blockhash = Hash.default()  # <-- Fix here
    mock_client_instance.get_latest_blockhash.return_value = mock_blockhash_resp
    
    # Transaction Broadcast mock
    mock_send_resp = MagicMock()
    mock_send_resp.value = "mock_tx_signature_abc123"
    mock_client_instance.send_transaction.return_value = mock_send_resp
    
    mock_async_client_class.return_value.__aenter__.return_value = mock_client_instance

    tx_signature = await cast_on_chain_vote(proposal_id="Prop-44", vote="NO")

    assert tx_signature == "mock_tx_signature_abc123"
    
    mock_client_instance.get_balance.assert_called_once_with(mock_wallet.pubkey())
    mock_client_instance.get_latest_blockhash.assert_called_once()
    mock_client_instance.send_transaction.assert_called_once()
    mock_client_instance.confirm_transaction.assert_called_once_with("mock_tx_signature_abc123")


@pytest.mark.asyncio
@patch("web3.governance.get_or_create_agent_wallet")
@patch("web3.governance.AsyncClient")
async def test_cast_on_chain_vote_triggers_airdrop_when_low_funds(mock_async_client_class, mock_get_wallet):
    """Verifies that a low balance wallet dynamically pulls down an airdrop before assembling the tx."""
    
    mock_wallet = Keypair()
    mock_get_wallet.return_value = mock_wallet
    mock_client_instance = AsyncMock()
    
    # Balance mock (0 Lamports - triggers airdrop workflow)
    mock_balance_resp = MagicMock()
    mock_balance_resp.value = 0
    mock_client_instance.get_balance.return_value = mock_balance_resp
    
    # Airdrop response mock
    mock_airdrop_resp = MagicMock()
    mock_airdrop_resp.value = "mock_airdrop_sig"
    mock_client_instance.request_airdrop.return_value = mock_airdrop_resp
    
    # Blockhash & Send mocks
    mock_blockhash_resp = MagicMock()
    mock_blockhash_resp.value.blockhash = Hash.default()  # <-- Fix here
    mock_client_instance.get_latest_blockhash.return_value = mock_blockhash_resp
    
    # Transaction Broadcast mock
    mock_send_resp = MagicMock()
    mock_send_resp.value = "final_tx_sig"
    mock_client_instance.send_transaction.return_value = mock_send_resp
    
    mock_async_client_class.return_value.__aenter__.return_value = mock_client_instance

    tx_signature = await cast_on_chain_vote(proposal_id="Prop-44", vote="YES")

    assert tx_signature == "final_tx_sig"
    mock_client_instance.request_airdrop.assert_called_once_with(mock_wallet.pubkey(), 10_000_000)
    
    calls = mock_client_instance.confirm_transaction.call_args_list
    assert calls[0][0][0] == "mock_airdrop_sig"
    assert calls[1][0][0] == "final_tx_sig"
