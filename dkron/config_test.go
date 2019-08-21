package dkron

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_normalizeAddrs(t *testing.T) {
	type test struct {
		config Config
		want   Config
	}

	tests := []test{
		{
			config: Config{BindAddr: "192.168.1.1:8946"},
			want:   Config{BindAddr: "192.168.1.1:8946"},
		},
		{
			config: Config{BindAddr: ":8946"},
			want:   Config{BindAddr: ":8946"},
		},
	}

	for _, tc := range tests {
		err := tc.config.normalizeAddrs()
		require.Error(t, err)
		assert.EqualValues(t, tc.want, tc.config)
	}
}

func Test_normalizeAdvertise(t *testing.T) {
	type test struct {
		addr string
		bind string
		dev  bool
		want string
	}

	tests := []test{
		{addr: "192.168.1.1", bind: ":8946", want: "192.168.1.1:8946", dev: false},
		{addr: "", bind: "127.0.0.1", want: "127.0.0.1:8946", dev: true},
	}

	for _, tc := range tests {
		addr, err := normalizeAdvertise(tc.addr, tc.bind, DefaultBindPort, tc.dev)
		require.NoError(t, err)
		assert.Equal(t, tc.want, addr)
	}
}
