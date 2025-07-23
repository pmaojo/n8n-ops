package client

// User represents a minimal n8n user.
type User struct {
	ID       string                 `json:"id"`
	Email    string                 `json:"email"`
	Settings map[string]interface{} `json:"settings"`
}
