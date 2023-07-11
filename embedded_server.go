package npoci

import (
	"fmt"
	"testing"
	"time"

	"github.com/tbeets/poci"

	"github.com/nats-io/nats-server/v2/server"
)

// DefaultOptions returns a set of default server options
func DefaultOptions() *server.Options {
	return &server.Options{
		Host:     "127.0.0.1",
		Port:     -1,
		HTTPPort: -1,
		Cluster:  server.ClusterOpts{Port: -1, Name: "abc"},
		NoLog:    true,
		NoSigs:   true,
		Debug:    true,
		Trace:    true,
	}
}

// RunServer starts a new Go Routine based server
func RunServer(opts *server.Options) *server.Server {
	if opts == nil {
		opts = DefaultOptions()
	}
	s, err := server.NewServer(opts)
	if err != nil || s == nil {
		panic(fmt.Sprintf("No NATS Server object returned: %v", err))
	}

	if !opts.NoLog {
		s.ConfigureLogger()
	}

	// Run server in Go routine.
	go s.Start()

	if !s.ReadyForConnections(10 * time.Second) {
		panic(err)
	}
	return s
}

// LoadConfig loads a configuration from a filename
func LoadConfig(configFile string) (opts *server.Options) {
	opts, err := server.ProcessConfigFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("Error processing configuration file: %v", err))
	}
	return
}

// MergeOptions merges two sets of server options, one considered overrides
func MergeOptions(base *server.Options, override *server.Options) (opts *server.Options) {
	if base == nil && override == nil {
		return nil
	}
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}
	return server.MergeOptions(base, override)
}

// UpByOpts run an embedded server from passed server options
func UpByOpts(opts *server.Options) (*server.Server, *server.Options) {
	s := RunServer(opts)
	return s, opts
}

// Up run an embedded server from a configuration file
func Up(configFile string) (s *server.Server, opts *server.Options) {
	opts = LoadConfig(configFile)
	s = RunServer(opts)
	return
}

// CheckLeafNodeConnectedCount checks that the server has the expected number of leaf node connections
func CheckLeafNodeConnectedCount(t *testing.T, s *server.Server, lnCons int) {
	t.Helper()
	poci.CheckFor(t, 5*time.Second, 15*time.Millisecond, func() error {
		if nln := s.NumLeafNodes(); nln != lnCons {
			return fmt.Errorf("expected %d connected leafnode(s) for server %v, got %d",
				lnCons, s, nln)
		}
		return nil
	})
}

var (
	_ = MergeOptions
	_ = UpByOpts
	_ = Up
	_ = CheckLeafNodeConnectedCount
)
