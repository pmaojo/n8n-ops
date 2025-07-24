package bubbleui

import "testing"

func TestGlamourRenderer_Render(t *testing.T) {
	r := NewGlamourRenderer("dark")
	out, err := r.Render([]byte("# Title"))
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("no output")
	}
}
