package tui

import (
	"strings"
	"testing"
)

func TestSyntheticDemoRendersSnapshotAndSpread(t *testing.T) {
	ui := newDemo()
	if count := ui.liveMarket.GetRowCount(); count != 3 {
		t.Fatalf("want two synthetic market rows and one header, got %d", count)
	}
	if text := ui.liveOpts.GetText(false); !strings.Contains(text, "74500") {
		t.Fatalf("synthetic option strike missing from %q", text)
	}
	if text := strings.TrimSpace(ui.liveFlow.GetCell(1, 1).Text); text != "180" {
		t.Fatalf("synthetic spread: want 180, got %q", text)
	}
}
