package dkron

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_normalizeAddrs(t *testing.T) {
	type test struct {
		config Config
		// After normalization, these fields should be populated
		expectAdvertiseSet bool
	}

	tests := []test{
		{
			config: Config{BindAddr: "192.168.1.1:8946"},
			expectAdvertiseSet: true, // Should successfully set advertise addr
		},
		{
			config: Config{BindAddr: ":8946"},
			expectAdvertiseSet: false, // `:8946` cannot be resolved, will error
		},
	}

	for _, tc := range tests {
		err := tc.config.normalizeAddrs()
		if tc.expectAdvertiseSet {
			// Should succeed and set advertise address
			require.NoError(t, err)
			assert.NotEmpty(t, tc.config.AdvertiseAddr, "AdvertiseAddr should be set")
		} else {
			// Should fail to resolve
			require.Error(t, err)
		}
	}
}

func Test_normalizeAdvertise(t *testing.T) {
	type test struct {
		addr         string
		bind         string
		dev          bool
		want         string
		wantAltIPv6  string // Alternative expected value for IPv6
	}

	tests := []test{
		{
			addr: "192.168.1.1",
			bind: ":8946",
			dev:  false,
			want: "192.168.1.1:8946",
		},
		{
			addr: "",
			bind: "127.0.0.1",
			dev:  true,
			want: "127.0.0.1:8946",
		},
		{
			addr: "",
			bind: "127.0.0.1:8946",
			dev:  true,
			want: "127.0.0.1:8946",
		},
		{
			addr:         "",
			bind:         "localhost:8946",
			dev:          true,
			want:         "127.0.0.1:8946",
			wantAltIPv6:  "[::1]:8946",
		},
	}

	for _, tc := range tests {
		addr, err := normalizeAdvertise(tc.addr, tc.bind, DefaultBindPort, tc.dev)
		require.NoError(t, err)
		if tc.wantAltIPv6 != "" && addr == tc.wantAltIPv6 {
			// Accept IPv6 loopback as well
			assert.Equal(t, tc.wantAltIPv6, addr)
		} else {
			assert.Equal(t, tc.want, addr)
		}
	}
}
