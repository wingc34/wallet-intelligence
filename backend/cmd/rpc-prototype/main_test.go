package main

import (
	"math/big"
	"testing"
)

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
				t.Fatalf("valiadateAddress(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}

			if !tt.wantErr && address.Hex() != tt.wantAddress {
				t.Errorf("valiadateAddress(%q) = %s, want %s", tt.input, address.Hex(), tt.wantAddress)
			}
		})
	}
}
