package integration_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/testutil"
)

var stopServer func()
var cleanup func()

func TestMain(m *testing.M) {
	var err error
	cleanup, err = testutil.BuildMockServer()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	stopServer, err = testutil.StartMockServer()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if cleanup != nil {
			cleanup()
		}
		os.Exit(1)
	}
	code := m.Run()
	if stopServer != nil {
		stopServer()
	}
	if cleanup != nil {
		cleanup()
	}
	os.Exit(code)
}
