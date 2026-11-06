package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const (
	sepoliaChainID = 11155111
	weiPerETH      = 18
)

type ChainIDProvider interface {
	ChainID(ctx context.Context) (*big.Int, error)
}

func main() {

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		return fmt.Errorf("SEPOLIA_RPC_URL is not set")
	}

	// Validate the number of command-line arguments.
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: go run ./cmd/rpc-prototype <wallet-address>")
	}

	addressInput := os.Args[1]
	// Validate the Ethereum wallet address.
	address, err := validateAddress(addressInput)
	if err != nil {
		return err
	}

	// Set a 10-second timeout for RPC operations.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	if err := validateChainID(ctx, client); err != nil {
		return err
	}

	balanceWei, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		return fmt.Errorf("failed to get wallet balance: %w", err)
	}

	fmt.Println("Network: Sepolia")
	fmt.Println("Address:", address.Hex())
	fmt.Printf("Balance: %s ETH\n", weiToETH(balanceWei))

	return nil
}

func validateAddress(input string) (common.Address, error) {
	if !common.IsHexAddress(input) {
		return common.Address{}, fmt.Errorf("invalid Ethereum wallet address")
	}

	return common.HexToAddress(input), nil
}

func validateChainID(ctx context.Context, client ChainIDProvider) error {
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	expectedChainID := big.NewInt(sepoliaChainID)
	if chainID.Cmp(expectedChainID) != 0 {
		return fmt.Errorf("wrong network: expected Sepolia (%s), got %s",
			expectedChainID, chainID)
	}

	return nil
}

func weiToETH(wei *big.Int) string {
	divisor := new(big.Int).Exp(
		big.NewInt(10),
		big.NewInt(weiPerETH),
		nil,
	)

	// get whole number
	whole := new(big.Int)
	fraction := new(big.Int)
	whole.QuoRem(wei, divisor, fraction)

	if fraction.Sign() == 0 {
		return whole.String()
	}

	// handle decimal place
	fractionText := fraction.String()
	// give zero to left if don't have enough weiPerETH
	fractionText = strings.Repeat(
		"0",
		weiPerETH-len(fractionText),
	) + fractionText

	fractionText = strings.TrimRight(fractionText, "0")

	return whole.String() + "." + fractionText
}
