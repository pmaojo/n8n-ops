package bubbleui

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/glamour"
)

// DocRenderer renders markdown to a viewable format.
type DocRenderer interface {
	Render(source []byte) (string, error)
}

// GlamourRenderer uses the glamour library for rendering.
type GlamourRenderer struct {
	style string
}

// NewGlamourRenderer creates a GlamourRenderer with the given style.
func NewGlamourRenderer(style string) *GlamourRenderer {
	return &GlamourRenderer{style: style}
}

func (g *GlamourRenderer) Render(source []byte) (string, error) {
	r, err := glamour.NewTermRenderer(glamour.WithStandardStyle(g.style))
	if err != nil {
		return "", err
	}
	out, err := r.RenderBytes(source)
	return string(out), err
}

func listDocs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var docs []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext == ".md" || ext == ".html" {
			docs = append(docs, e.Name())
		}
	}
	sort.Strings(docs)
	return docs, nil
}
