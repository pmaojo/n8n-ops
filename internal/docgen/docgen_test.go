package docgen

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestCobraGenerator_Generate(t *testing.T) {
	tmp := t.TempDir()

	root := &cobra.Command{Use: "root", Short: "root cmd"}
	root.AddCommand(&cobra.Command{Use: "sub", Short: "sub cmd"})

	g := NewCobraGenerator(root)
	if err := g.Generate(context.Background(), "md", tmp); err != nil {
		t.Fatalf("generate markdown: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "root.md")); err != nil {
		t.Fatalf("markdown not generated: %v", err)
	}

	if err := g.Generate(context.Background(), "html", tmp); err != nil {
		t.Fatalf("generate html: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "root.html")); err != nil {
		t.Fatalf("html not generated: %v", err)
	}
}
