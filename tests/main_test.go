package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/testutil"
)

var stopServer func()

func TestMain(m *testing.M) {
	var err error
	stopServer, _, err = testutil.SetupMockServer()
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
