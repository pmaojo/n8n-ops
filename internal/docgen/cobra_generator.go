package docgen

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/yuin/goldmark"
)

// CobraGenerator implements DocGenerator using a Cobra command tree.
type CobraGenerator struct {
	root *cobra.Command
}

// NewCobraGenerator creates a new generator for the given root command.
func NewCobraGenerator(root *cobra.Command) *CobraGenerator {
	return &CobraGenerator{root: root}
}

// Generate writes the command reference in the specified format to outputDir.
// Supported formats: "md" and "html".
func (g *CobraGenerator) Generate(ctx context.Context, format, outputDir string) error {
	if format != "md" && format != "html" {
		return fmt.Errorf("unsupported format: %s", format)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}
	if err := doc.GenMarkdownTree(g.root, outputDir); err != nil {
		return err
	}
	if format == "html" {
		return convertMarkdownDir(outputDir)
	}
	return nil
}

func convertMarkdownDir(dir string) error {
	md := goldmark.New()
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err := md.Convert(data, &buf); err != nil {
			return err
		}
		htmlPath := filepath.Join(dir, filepath.Base(path[:len(path)-3]+".html"))
		if err := os.WriteFile(htmlPath, buf.Bytes(), 0o644); err != nil {
			return err
		}
		return os.Remove(path)
	})
}
