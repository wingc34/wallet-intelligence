package main

import (
	"context"
	"fmt"
	"math/big"
	"testing"
)

type MockChainIDProvider struct {
	chainID *big.Int
	err     error
}

func (m *MockChainIDProvider) ChainID(ctx context.Context) (*big.Int, error) {
	return m.chainID, m.err
}

func TestWeiToETH(t *testing.T) {
	tests := []struct {
		name string
		wei  string
		want string
	}{
		{name: "0", wei: "0", want: "0"}, {
			name: "One ETH",
			wei:  "1000000000000000000", want: "1",
		},
		{
			name: "Decimal",
			wei:  "1250000000000000000", want: "1.25",
		},
		{
			name: "Leading zeros",
			wei:  "100000000000000000", want: "0.1",
		},
		{
			name: "Trailing zeros",
			wei:  "1234000000000000000", want: "1.234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wei, _ := new(big.Int).SetString(tt.wei, 10)

			got := weiToETH(wei)

			if got != tt.want {
				t.Errorf("weiToETH(%s) = %s, want %s", tt.wei, got, tt.want)
			}
		})
	}
}

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantAddress string
		wantErr     bool
	}{
		{
			name:        "valid address",
			input:       "0x000000000000000000000000000000000000dEaD",
			wantAddress: "0x000000000000000000000000000000000000dEaD",
			wantErr:     false,
		},
		{
			name:    "invalid address",
			input:   "invalid-address",
			wantErr: true,
		},
		{
			name:    "empty address",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, err := validateAddress(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("validateAddress(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}

			if !tt.wantErr && address.Hex() != tt.wantAddress {
				t.Errorf("validateAddress(%q) = %s, want %s", tt.input, address.Hex(), tt.wantAddress)
			}
		})
	}
}

func TestValidateChainID(t *testing.T) {
	tests := []struct {
		name    string
		chainID *big.Int
		err     error
		wantErr bool
	}{
		{
			name:    "valid Sepolia network",
			chainID: big.NewInt(11155111),
			err:     nil,
			wantErr: false,
		},
		{
			name:    "wrong network",
			chainID: big.NewInt(1),
			err:     nil,
			wantErr: true,
		},
		{
			name:    "RPC error",
			chainID: nil,
			err:     fmt.Errorf("connection timeout"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockChainIDProvider{
				chainID: tt.chainID,
				err:     tt.err,
			}

			err := validateChainID(context.Background(), mock)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateChainID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
