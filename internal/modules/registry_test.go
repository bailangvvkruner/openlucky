package modules

import "testing"

func TestRegistryOrdersModules(t *testing.T) {
	registry := DefaultRegistry()
	modules := registry.List()
	if len(modules) == 0 {
		t.Fatal("modules is empty")
	}
	if modules[0].ID != "status" {
		t.Fatalf("first module = %q, want status", modules[0].ID)
	}
	if modules[len(modules)-1].State != StateStub {
		t.Fatalf("last module state = %q, want stub", modules[len(modules)-1].State)
	}
}
