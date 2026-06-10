package config

import (
	"path/filepath"
	"testing"
)

func TestStoreSaveLoad(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "openlucky.json"))
	want := Default()
	want.Theme = "graphite-cyan"
	want.Modules["cron"] = false
	if err := store.Save(want); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got.Theme != want.Theme {
		t.Fatalf("Theme = %q, want %q", got.Theme, want.Theme)
	}
	if got.Modules["cron"] {
		t.Fatal("Modules[cron] = true, want false")
	}
}
