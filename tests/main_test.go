package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/testutil"
)

var stopServer func()

func TestMain(m *testing.M) {
	cleanup, err := testutil.BuildMockServer()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	stopServer, err = testutil.StartMockServer(0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		cleanup()
		os.Exit(1)
	}
	defer cleanup()
	code := m.Run()
	if stopServer != nil {
		stopServer()
	}
	os.Exit(code)
}
