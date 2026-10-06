import os
import pytest
from unittest.mock import MagicMock, AsyncMock, patch
from solders.keypair import Keypair
from web3.wallet import get_or_create_agent_wallet, ensure_funds


def test_wallet_creation_and_persistence(tmp_path):
    """Verifies that a new wallet is generated, saved to disk, and reloaded correctly."""
    
    mock_dir = str(tmp_path)
    mock_key_path = os.path.join(mock_dir, "agent_keypair.json")
    
    def mock_join(path, filename):
        if filename == "agent_keypair.json":
            return mock_key_path
        return os.path.join(path, filename)

    with patch("os.path.join", side_effect=mock_join):
        
        assert not os.path.exists(mock_key_path)
        wallet_1 = get_or_create_agent_wallet()
        
        assert isinstance(wallet_1, Keypair)
        assert os.path.exists(mock_key_path)
        
        wallet_2 = get_or_create_agent_wallet()
        assert wallet_1.pubkey() == wallet_2.pubkey()


@pytest.mark.asyncio
async def test_ensure_funds_sufficient_balance():
    """Verifies no airdrop is requested if the wallet already has enough SOL."""

    mock_balance_resp = MagicMock()
    mock_balance_resp.value = 50_000_000
    
    mock_client = AsyncMock()
    mock_client.get_balance.return_value = mock_balance_resp
    
    wallet = Keypair()
    await ensure_funds(mock_client, wallet)
    
    mock_client.get_balance.assert_called_once_with(wallet.pubkey())
    mock_client.request_airdrop.assert_not_called()


@pytest.mark.asyncio
async def test_ensure_funds_requests_airdrop_when_low():
    """Verifies that an airdrop is requested and confirmed if the balance is low."""
    mock_balance_resp = MagicMock()
    mock_balance_resp.value = 5_000_000
    mock_airdrop_resp = MagicMock()
    mock_airdrop_resp.value = "mock_tx_signature"
    
    mock_client = AsyncMock()
    mock_client.get_balance.return_value = mock_balance_resp
    mock_client.request_airdrop.return_value = mock_airdrop_resp
    mock_client.confirm_transaction = AsyncMock()
    
    wallet = Keypair()
    await ensure_funds(mock_client, wallet)
    
    mock_client.get_balance.assert_called_once_with(wallet.pubkey())
    mock_client.request_airdrop.assert_called_once_with(wallet.pubkey(), 10_000_000)
    mock_client.confirm_transaction.assert_called_once_with("mock_tx_signature")


@pytest.mark.asyncio
async def test_ensure_funds_handles_airdrop_failure():
    """Verifies that network exceptions or rate limits during airdrops are gracefully handled."""
    mock_balance_resp = MagicMock()
    mock_balance_resp.value = 0
    
    mock_client = AsyncMock()
    mock_client.get_balance.return_value = mock_balance_resp
    mock_client.request_airdrop.side_effect = Exception("Rate limit exceeded")
    
    wallet = Keypair()
    result = await ensure_funds(mock_client, wallet)
    
    assert result == "Failed: No funds available"
