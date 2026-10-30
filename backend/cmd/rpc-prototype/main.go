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

func main() {
	if err := godotenv.Load(); err != nil {
		exitWithError("failed to load .env file: " + err.Error())
	}

	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	if rpcURL == "" {
		exitWithError("SEPOLIA_RPC_URL is not set")
	}

	//check value input number
	if len(os.Args) != 2 {
		exitWithError("usage: go run ./cmd/rpc-prototype <wallet-address>")
	}

	addressInput := os.Args[1]
	// check address valid
	if !common.IsHexAddress(addressInput) {
		exitWithError("invalid Ethereum wallet address")
	}

	address := common.HexToAddress(addressInput)
	// wait 10 sec if no response will cancel
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		exitWithError("failed to connnect to RPC: " + err.Error())
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		exitWithError("failed to get chain ID: " + err.Error())
	}

	expectedChainID := big.NewInt(sepoliaChainID)
	if chainID.Cmp(expectedChainID) != 0 {
		exitWithError(fmt.Sprintf("wrong network: expected Sepolia (%s), got %s",
			expectedChainID, chainID))
	}

	balanceWei, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		exitWithError("failed to get wallet balance: " + err.Error())
	}

	fmt.Println("Network: Sepolia")
	fmt.Println("Address:", address.Hex())
	fmt.Printf("Balance: %s ETH\n", weiToETH(balanceWei))
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

func exitWithError(message string) {
	fmt.Fprintln(os.Stderr, "Error:", message)
	os.Exit(1)
}
