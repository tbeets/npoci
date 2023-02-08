package npoci

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

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

func UpByOpts(opts *server.Options) (*server.Server, *server.Options) {
	s := RunServer(opts)
	return s, opts
}

func Up(configFile string) (s *server.Server, opts *server.Options) {
	opts = LoadConfig(configFile)
	s = RunServer(opts)
	return
}
