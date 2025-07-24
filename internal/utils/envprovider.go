package utils

import "os"

// EnvProvider defines the interface for retrieving and setting environment variables.
type EnvProvider interface {
	Getenv(key string) string
	Setenv(key, value string) error
}

// OSProvider implements EnvProvider using the operating system environment.
type OSProvider struct{}

// Getenv retrieves the environment variable named by the key.
func (OSProvider) Getenv(key string) string { return os.Getenv(key) }

// Setenv sets the value of the environment variable named by the key.
func (OSProvider) Setenv(key, value string) error { return os.Setenv(key, value) }
