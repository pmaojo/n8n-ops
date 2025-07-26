package utils

import (
	"fmt"
	"strings"
)

// BuildEnvVarNames returns the full and short environment variable names for the
// given base prefix and environment. The full name is constructed as
// `<prefix>_<ENV>` where `ENV` is the uppercased environment string. The short
// name follows the common convention DEV/PROD for development/production and
// uses the uppercased environment name for all others.
func BuildEnvVarNames(prefix, environment string) (full, short string) {
	envUpper := strings.ToUpper(environment)
	full = fmt.Sprintf("%s_%s", prefix, envUpper)

	switch envUpper {
	case "DEVELOPMENT":
		short = fmt.Sprintf("%s_DEV", prefix)
	case "PRODUCTION":
		short = fmt.Sprintf("%s_PROD", prefix)
	default:
		short = fmt.Sprintf("%s_%s", prefix, envUpper)
	}

	return full, short
}
