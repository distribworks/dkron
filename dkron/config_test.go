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

// TestDockerCustomAddressPool tests the scenario from the issue where
// Docker uses custom address pools and assigns IPs like 172.80.0.2
func TestDockerCustomAddressPool(t *testing.T) {
	// Simulate Docker scenario with custom address pool
	// The bind address would be resolved from template "{{ GetPrivateIP }}:8946"
	// to something like "172.80.0.2:8946"
	
	cfg := Config{
		BindAddr: "172.80.0.2:8946",  // Simulates Docker custom pool IP
	}
	
	err := cfg.normalizeAddrs()
	
	// Should not error - this was the bug
	require.NoError(t, err, "normalizeAddrs should not fail with IP:port in bind address")
	
	// Should have set advertise address
	assert.NotEmpty(t, cfg.AdvertiseAddr, "AdvertiseAddr should be set")
	
	// The advertise address should contain the IP from bind
	assert.Contains(t, cfg.AdvertiseAddr, "172.80.0.2", "AdvertiseAddr should use the bind IP")
}
