package utils

import "testing"

func TestBuildEnvVarNames(t *testing.T) {
	cases := []struct {
		prefix string
		env    string
		full   string
		short  string
	}{
		{"N8N_API_KEY", "development", "N8N_API_KEY_DEVELOPMENT", "N8N_API_KEY_DEV"},
		{"N8N_API_KEY", "production", "N8N_API_KEY_PRODUCTION", "N8N_API_KEY_PROD"},
		{"N8N_API_KEY", "staging", "N8N_API_KEY_STAGING", "N8N_API_KEY_STAGING"},
	}

	for _, c := range cases {
		full, short := BuildEnvVarNames(c.prefix, c.env)
		if full != c.full || short != c.short {
			t.Fatalf("%s %s expected (%s,%s) got (%s,%s)", c.prefix, c.env, c.full, c.short, full, short)
		}
	}
}
