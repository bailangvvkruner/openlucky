package logs

import "testing"

func TestQueryReturnsNewestEntries(t *testing.T) {
	store := NewStore(2)
	store.Add("info", "test", "one")
	store.Add("info", "test", "two")
	store.Add("info", "test", "three")
	entries := store.Query(10)
	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(entries))
	}
	if entries[0].Message != "two" || entries[1].Message != "three" {
		t.Fatalf("entries = %#v", entries)
	}
}
