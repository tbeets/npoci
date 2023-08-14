package npoci

import (
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// ClientConnect creates a client connection from passed server reference taking optional options
// e.g. user credentials nats.UserCredentials("path/to/creds") or nats.UserInfo("user", "password")
func ClientConnect(t testing.TB, s *server.Server, opts ...nats.Option) *nats.Conn {
	t.Helper()
	nc, err := nats.Connect(s.ClientURL(), opts...)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return nc
}

// ClientConnectURL creates a client connection from passed NATS URL taking optional options
// e.g. user credentials nats.UserCredentials("path/to/creds") or nats.UserInfo("user", "password")
func ClientConnectURL(t testing.TB, url string, opts ...nats.Option) *nats.Conn {
	t.Helper()
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return nc
}

var (
	_ = ClientConnect
	_ = ClientConnectURL
)
