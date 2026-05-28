package manifest

import (
	"sort"
	"testing"
	"time"
)

func TestSeed_NonEmpty(t *testing.T) {
	if len(Seed()) == 0 {
		t.Fatal("empty")
	}
}

func TestSeed_EntryCount(t *testing.T) {
	const expected = 10
	if got := len(Seed()); got != expected {
		t.Fatalf("got %d, want %d", got, expected)
	}
}

func TestSeed_AllKeysNonEmpty(t *testing.T) {
	for i, e := range Seed() {
		if e.Key == "" {
			t.Errorf("%d empty Key", i)
		}
	}
}

func TestSeed_AllKeysUnique(t *testing.T) {
	seen := map[string]int{}
	for i, e := range Seed() {
		if prev, ok := seen[e.Key]; ok {
			t.Errorf("dup at %d and %d", prev, i)
		}
		seen[e.Key] = i
	}
}

func TestSeed_AllSourceValuesCanonical(t *testing.T) {
	allowed := map[string]bool{}
	for _, s := range AllSources() {
		allowed[s] = true
	}
	for _, e := range Seed() {
		if !allowed[e.Source] {
			t.Errorf("%q: Source %q not in AllSources", e.Key, e.Source)
		}
	}
}

func TestSeed_AllSchemaVersionsCurrent(t *testing.T) {
	for _, e := range Seed() {
		if e.SchemaVersion != SchemaVersion {
			t.Errorf("%q: %d", e.Key, e.SchemaVersion)
		}
	}
}

func TestSeed_CorpusPinsAreStale(t *testing.T) {
	now := time.Now()
	for _, e := range Seed() {
		if len(e.Key) >= 6 && e.Key[:6] == "corpus" && !e.IsStale(now, 30*24*time.Hour) {
			t.Errorf("%q not stale", e.Key)
		}
	}
}

func TestSeed_AllPinnedEntriesHaveVersion(t *testing.T) {
	for _, e := range Seed() {
		if e.Version == "" {
			t.Errorf("%q empty Version", e.Key)
		}
	}
}

func TestSchemaVersion_PinnedAtV1(t *testing.T) {
	if SchemaVersion != 1 {
		t.Fatalf("got %d", SchemaVersion)
	}
}

func TestIsStale_SentinelAlwaysTrue(t *testing.T) {
	if !(Entry{FreshAt: FreshAtUnknown}).IsStale(time.Now(), time.Hour) {
		t.Error()
	}
}

func TestIsStale_FreshNotStale(t *testing.T) {
	now := time.Now()
	if (Entry{FreshAt: now.Add(-1 * time.Hour)}).IsStale(now, 24*time.Hour) {
		t.Error()
	}
}

func TestSortedKeys_Deterministic(t *testing.T) {
	if !sort.StringsAreSorted(Seed().SortedKeys()) {
		t.Fatal()
	}
}

func TestAllSources_NonEmpty(t *testing.T) {
	if len(AllSources()) == 0 {
		t.Fatal()
	}
}

func TestStaleEntries_IncludesAllSentinels(t *testing.T) {
	if got := len(Seed().StaleEntries(time.Now(), 30*24*time.Hour)); got < 3 {
		t.Errorf("got %d", got)
	}
}
