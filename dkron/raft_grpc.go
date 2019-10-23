package dkron

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/raft"
)

// RaftLayer is the network layer for internode communications.
type RaftLayer struct {
	TLSConfig *tls.Config

	ln net.Listener
}

// NewRaftLayer returns an initialized unecrypted RaftLayer.
func NewRaftLayer() *RaftLayer {
	return &RaftLayer{}
}

// NewTLSRaftLayer returns an initialized TLS-ecrypted RaftLayer.
func NewTLSRaftLayer(tlsConfig *tls.Config) *RaftLayer {
	return &RaftLayer{TLSConfig: tlsConfig}
}

// Open opens the RaftLayer, binding to the supplied address.
func (t *RaftLayer) Open(l net.Listener) error {
	t.ln = l
	return nil
}

// Dial opens a network connection.
func (t *RaftLayer) Dial(addr raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}

	var err error
	var conn net.Conn
	if t.TLSConfig != nil {
		log.Debug("doing a TLS dial")
		conn, err = tls.DialWithDialer(dialer, "tcp", string(addr), t.TLSConfig)
	} else {
		conn, err = dialer.Dial("tcp", string(addr))
	}

	return conn, err
}

// Accept waits for the next connection.
func (t *RaftLayer) Accept() (net.Conn, error) {
	c, err := t.ln.Accept()
	if err != nil {
		fmt.Println("error accepting: ", err.Error())
	}
	return c, err
}

// Close closes the RaftLayer
func (t *RaftLayer) Close() error {
	return t.ln.Close()
}

// Addr returns the binding address of the RaftLayer.
func (t *RaftLayer) Addr() net.Addr {
	return t.ln.Addr()
}
