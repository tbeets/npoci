package npoci

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/nats-io/nats-server/v2/server"
)

func httpMonitorGet(t *testing.T, httpHost string, httpPort int, path string) []byte {
	t.Helper()
	if httpHost == "" {
		httpHost = "127.0.0.1"
	}
	if httpPort < 80 {
		httpPort = 8222
	}
	if path == "" {
		path = "/"
	}
	url := fmt.Sprintf("http://%s:%d%s", httpHost, httpPort, path)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Expected no error: Got %v\n", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected a 200 response, got %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Got an error reading the body: %v\n", err)
	}
	return body
}

// HttpMonitorGetVarz returns the Varz struct from the NATS Server HTTP Monitor endpoint
func HttpMonitorGetVarz(t *testing.T, httpHost string, httpPort int, opts string) *server.Varz {
	t.Helper()
	page := "/varz"
	if opts != "" {
		page = fmt.Sprintf("%s?%s", page, opts)
	}
	body := httpMonitorGet(t, httpHost, httpPort, page)
	v := server.Varz{}
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatalf("Got an error unmarshalling the body: %v\n", err)
	}
	return &v
}

// connz
// subsz
// routez
// gatewayz
// leafz
