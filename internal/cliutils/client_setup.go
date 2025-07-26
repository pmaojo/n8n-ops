package cliutils

import (
	"errors"
	"fmt"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/credentials"
	"github.com/pmaojo/n8n-ops/internal/utils"
)

// MissingCredentialError indicates required credentials are not configured.
// It reports which credential fields are missing.
type MissingCredentialError struct {
	Missing []string
}

func (e MissingCredentialError) Error() string {
	return fmt.Sprintf("missing credentials: %v", e.Missing)
}

// ErrMissingCredentials is returned when n8n URL or API key is absent.
var ErrMissingCredentials = errors.New("missing n8n credentials")

// SetupClient creates a n8n client using the credential manager for the given
// environment. Demo mode returns an in-memory client. The credential manager is
// returned for further use by callers.
func SetupClient(env string, demo bool) (client.Client, *credentials.CredentialManager, error) {
	cm := credentials.NewCredentialManager(utils.OSProvider{}, env)

	if demo {
		return client.NewDemoN8nClient(), cm, nil
	}

	url, key, err := cm.GetN8nCredentials()
	if err != nil {
		return nil, cm, fmt.Errorf("failed to load credentials: %w", err)
	}

	missing := []string{}
	if url == "" {
		missing = append(missing, "n8n_url")
	}
	if key == "" {
		missing = append(missing, "n8n_api_key")
	}
	if len(missing) > 0 {
		return nil, cm, MissingCredentialError{Missing: missing}
	}

	c, err := client.New(url, key, nil)
	if err != nil {
		return nil, cm, fmt.Errorf("failed to create n8n client: %w", err)
	}
	return c, cm, nil
}
