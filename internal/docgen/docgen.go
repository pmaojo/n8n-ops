package docgen

import "context"

// DocGenerator generates documentation for commands.
type DocGenerator interface {
	// Generate writes documentation for the provided format ("md" or "html")
	// into the outputDir.
	Generate(ctx context.Context, format, outputDir string) error
}
