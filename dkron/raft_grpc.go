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
	ln net.Listener

	certFile        string // Path to local X.509 cert.
	certKey         string // Path to corresponding X.509 key.
	remoteEncrypted bool   // Remote nodes use encrypted communication.
	skipVerify      bool   // Skip verification of remote node certs.
}

// NewRaftLayer returns an initialized unecrypted RaftLayer.
func NewRaftLayer() *RaftLayer {
	return &RaftLayer{}
}

// NewTLSRaftLayer returns an initialized TLS-ecrypted RaftLayer.
func NewTLSRaftLayer(certFile, keyPath string, skipVerify bool) *RaftLayer {
	return &RaftLayer{
		certFile:        certFile,
		certKey:         keyPath,
		remoteEncrypted: true,
		skipVerify:      skipVerify,
	}
}

// Open opens the RaftLayer, binding to the supplied address.
func (t *RaftLayer) Open(l net.Listener) error {
	if t.certFile != "" {
		config, err := createTLSConfig(t.certFile, t.certKey)
		if err != nil {
			return err
		}
		l = tls.NewListener(l, config)
	}

	t.ln = l
	return nil
}

// Dial opens a network connection.
func (t *RaftLayer) Dial(addr raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}

	var err error
	var conn net.Conn
	if t.remoteEncrypted {
		conf := &tls.Config{
			InsecureSkipVerify: t.skipVerify,
		}
		fmt.Println("doing a TLS dial")
		conn, err = tls.DialWithDialer(dialer, "tcp", string(addr), conf)
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

// createTLSConfig returns a TLS config from the given cert and key.
func createTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	var err error
	config := &tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return config, nil
}
