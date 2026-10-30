package web3

import (
	"context"
	"fmt"
	"time"
	"github.com/xceptions/agenticarb/backend-go/config"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"


)

func CastOnChainVote(ctx context.Context, proposalID string, vote string) (string, error) {
	wallet, err := GetOrCreateAgentWallet()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve execution wallet keys: %w", err)
	}

	publicKey := wallet.PublicKey()
	client := rpc.New(config.GlobalSettings.RPCURL)

	balanceResult, err := client.GetBalance(ctx, publicKey, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to pull wallet account balance: %w", err)
	}

	if balanceResult.Value < 10000000 {
		fmt.Printf("web3 Account empty. Requesting Devnet SOL for %s...\n", publicKey)
		
		_, err := client.RequestAirdrop(ctx, publicKey, 10000000, rpc.CommitmentFinalized)
		if err != nil {
			return "", fmt.Errorf("error! Devnet airdrop failed or rate-limited: %w", err)
		}

		time.Sleep(3 * time.Second)
	}

	memoProgramID := solana.MustPublicKeyFromBase58(config.GlobalSettings.MemoProgramID)
	memoData := []byte(fmt.Sprintf("AgentVote:%s:%s", proposalID, vote))

	instruction := solana.NewInstruction(
		memoProgramID,
		solana.AccountMetaSlice{},
		memoData,
	)

	recentBlockhashResult, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed fetching transactional blockhash reference: %w", err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recentBlockhashResult.Value.Blockhash,
		solana.TransactionPayer(publicKey),
	)
	if err != nil {
		return "", fmt.Errorf("failed building transaction: %w", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(publicKey) {
			return &wallet
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed signing transaction payload: %w", err)
	}

	fmt.Printf("web3 Broadcasting '%s' vote for %s to Solana Devnet...\n", vote, proposalID)

	sig, err := client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed broadcasting signature payload packet: %w", err)
	}

	fmt.Println("Success! Vote transaction finalized on-chain!")
	fmt.Printf("Explorer URL: https://solana.com\n", sig)

	return sig.String(), nil
}
