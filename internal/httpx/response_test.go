package httpx

import "testing"

func TestOKEnvelope(t *testing.T) {
	got := OK(map[string]string{"service": "openlucky"})
	if got.Code != 0 {
		t.Fatalf("Code = %d, want 0", got.Code)
	}
	if got.Message != "ok" {
		t.Fatalf("Message = %q, want ok", got.Message)
	}
	if got.Data == nil {
		t.Fatal("Data is nil")
	}
}

func TestErrorEnvelope(t *testing.T) {
	got := Error("unauthorized", "login required")
	if got.Code == 0 {
		t.Fatal("Code = 0, want non-zero")
	}
	if got.Error == nil {
		t.Fatal("Error is nil")
	}
	if got.Error.Code != "unauthorized" {
		t.Fatalf("Error.Code = %q", got.Error.Code)
	}
}
