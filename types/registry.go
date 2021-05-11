package types

import "time"

type DockerOption struct {
	// Auth
	UserName string
	Password string

	// RegistryToken is a bearer token to be sent to a registry
	RegistryToken string

	// GCP
	GcpCredPath string

	InsecureSkipTLSVerify bool
	NonSSL                bool
	SkipPing              bool // this is ignored now
	Timeout               time.Duration
}

