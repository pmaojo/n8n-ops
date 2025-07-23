package providers

import (
        "fmt"
)

// NewGitProvider creates a new Git provider based on the configuration
func NewGitProvider(config *ProviderConfig) (GitProvider, error) {
        if config == nil {
                return nil, fmt.Errorf("provider config is required")
        }

        switch config.Type {
        case ProviderTypeGitLab:
                return NewGitLabProvider(config)
        // case ProviderTypeGitHub:
        //      return NewGitHubProvider(config)
        default:
                return nil, fmt.Errorf("unsupported provider type: %s (only GitLab is currently supported)", config.Type)
        }
}

// GetSupportedProviders returns a list of supported provider types
func GetSupportedProviders() []ProviderType {
        return []ProviderType{
                ProviderTypeGitLab,
                // Future providers can be added here when needed:
                // ProviderTypeGitHub,
                // ProviderTypeBitbucket,
                // ProviderTypeGitea,
        }
}

// IsProviderSupported checks if a provider type is supported
func IsProviderSupported(providerType ProviderType) bool {
        supported := GetSupportedProviders()
        for _, p := range supported {
                if p == providerType {
                        return true
                }
        }
        return false
}

// GetProviderInfo returns information about a provider
func GetProviderInfo(providerType ProviderType) (*ProviderInfo, error) {
        switch providerType {
        case ProviderTypeGitLab:
                return &ProviderInfo{
                        Name:        "GitLab",
                        Type:        ProviderTypeGitLab,
                        APIVersion:  "v4",
                        DefaultURL:  "https://gitlab.com",
                        AuthMethod:  "PRIVATE-TOKEN",
                        Description: "GitLab with Personal Access Token authentication",
                }, nil
        // Future providers can be uncommented when needed:
        // case ProviderTypeGitHub:
        //      return &ProviderInfo{
        //              Name:        "GitHub",
        //              Type:        ProviderTypeGitHub,
        //              APIVersion:  "v3",
        //              DefaultURL:  "https://api.github.com",
        //              AuthMethod:  "Bearer Token",
        //              Description: "GitHub with Personal Access Token authentication",
        //      }, nil
        default:
                return nil, fmt.Errorf("unsupported provider type: %s (only GitLab is currently supported)", providerType)
        }
}

// ProviderInfo contains metadata about a Git provider
type ProviderInfo struct {
        Name        string       `json:"name"`
        Type        ProviderType `json:"type"`
        APIVersion  string       `json:"api_version"`
        DefaultURL  string       `json:"default_url"`
        AuthMethod  string       `json:"auth_method"`
        Description string       `json:"description"`
}