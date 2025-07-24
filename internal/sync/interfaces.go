package sync

import (
	"io"
	"os"
)

// FileSystem abstracts filesystem operations used by Service.
type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (io.WriteCloser, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
}

// OSFileSystem is the default FileSystem implementation using the os package.
type OSFileSystem struct{}

func (OSFileSystem) MkdirAll(path string, perm os.FileMode) error { return os.MkdirAll(path, perm) }
func (OSFileSystem) Create(name string) (io.WriteCloser, error)   { return os.Create(name) }
func (OSFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// EnvProvider abstracts environment variable access.
type EnvProvider interface{ Get(key string) string }

// OSEnvProvider is the default EnvProvider implementation using os.Getenv.
type OSEnvProvider struct{}

func (OSEnvProvider) Get(key string) string { return os.Getenv(key) }
