package integration_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/testutil"
)

var stopServer func()

func TestMain(m *testing.M) {
	var err error
	stopServer, err = testutil.StartMockServer()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	code := m.Run()
	if stopServer != nil {
		stopServer()
	}
	os.Exit(code)
}
