package status

import (
	"testing"
	"time"
)

func TestHostOverview(t *testing.T) {
	overview := Host(time.Now().Add(-time.Second))
	if overview.Service != "OpenLucky" {
		t.Fatalf("Service = %q", overview.Service)
	}
	if overview.Runtime != "go+hertz" {
		t.Fatalf("Runtime = %q", overview.Runtime)
	}
	if overview.GoVersion == "" {
		t.Fatal("GoVersion is empty")
	}
}
