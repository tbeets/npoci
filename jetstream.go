package npoci

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tbeets/poci"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// JsClientConnect creates a client connection and JetStream context from passed server reference taking optional options
// e.g. user credentials nats.UserCredentials("path/to/creds") or nats.UserInfo("user", "password")
func JsClientConnect(t *testing.T, s *server.Server, opts ...nats.Option) (*nats.Conn, nats.JetStreamContext) {
	t.Helper()
	nc, err := nats.Connect(s.ClientURL(), opts...)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	js, err := nc.JetStream(nats.MaxWait(10 * time.Second))
	if err != nil {
		t.Fatalf("Unexpected error getting JetStream context: %v", err)
	}
	return nc, js
}

func GetJsStreamInfo(t *testing.T, nc *nats.Conn, stream string, apiPrefix string) *server.StreamInfo {
	t.Helper()
	apiDest := fmt.Sprintf(server.JSApiStreamInfoT, stream)
	if apiPrefix != "" {
		apiDest = strings.Replace(apiDest, "$JS.API", apiPrefix, 1)
	}
	rmsg, err := nc.Request(apiDest, []byte("{}"), 2*time.Second)
	poci.RequireNoError(t, err)
	respSi := server.JSApiStreamInfoResponse{}
	err = json.Unmarshal(rmsg.Data, &respSi)
	poci.RequireNoError(t, err)
	if respSi.Error != nil {
		t.Logf("Response Error: %v", respSi.Error)
	}
	poci.RequireTrue(t, respSi.Error == nil)
	return respSi.StreamInfo
}

func JsCreateStreamFromFile(t *testing.T, nc *nats.Conn, stream string, apiPrefix string, jsReqFile string) {
	t.Helper()
	jsReq, err := os.ReadFile(jsReqFile)
	poci.RequireNoError(t, err)
	JsCreateStream(t, nc, stream, apiPrefix, jsReq)
}

func JsCreateStream(t *testing.T, nc *nats.Conn, stream string, apiPrefix string, jsReq []byte) {
	t.Helper()
	apiDest := fmt.Sprintf(server.JSApiStreamCreateT, stream)
	if apiPrefix != "" {
		apiDest = strings.Replace(apiDest, "$JS.API", apiPrefix, 1)
	}
	rmsg, err := nc.Request(apiDest, jsReq, 2*time.Second)
	poci.RequireNoError(t, err)
	respC := server.JSApiStreamCreateResponse{}
	err = json.Unmarshal(rmsg.Data, &respC)
	poci.RequireNoError(t, err)
	if respC.Error != nil {
		t.Logf("Response error: %v", respC.Error)
	}
	poci.RequireTrue(t, respC.Error == nil)
}

func GetJsConsumerInfo(t *testing.T, nc *nats.Conn, stream string, consumer string, apiPrefix string) *server.ConsumerInfo {
	t.Helper()
	apiDest := fmt.Sprintf(server.JSApiConsumerInfoT, stream, consumer)
	if apiPrefix != "" {
		apiDest = strings.Replace(apiDest, "$JS.API", apiPrefix, 1)
	}
	rmsg, err := nc.Request(apiDest, []byte("{}"), 2*time.Second)
	poci.RequireNoError(t, err)
	respCi := server.JSApiConsumerInfoResponse{}
	err = json.Unmarshal(rmsg.Data, &respCi)
	poci.RequireNoError(t, err)
	if respCi.Error != nil {
		t.Logf("Response Error: %v", respCi.Error)
	}
	poci.RequireTrue(t, respCi.Error == nil)
	return respCi.ConsumerInfo
}

var (
	_ = JsClientConnect
	_ = GetJsStreamInfo
	_ = GetJsConsumerInfo
	_ = JsCreateStream
	_ = JsCreateStreamFromFile
)
