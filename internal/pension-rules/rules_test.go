package pensionrules

import (
	"crypto/sha256"
	"testing"
)

func TestAllPins_NonEmpty(t *testing.T) {
	if len(AllPins()) == 0 {
		t.Fatal()
	}
}

func TestAllPins_Count3(t *testing.T) {
	if got := len(AllPins()); got != 3 {
		t.Fatalf("got %d", got)
	}
}

func TestAllPins_UniqueIDs(t *testing.T) {
	seen := map[CorpusID]int{}
	for i, p := range AllPins() {
		if prev, ok := seen[p.ID]; ok {
			t.Errorf("dup %d %d", prev, i)
		}
		seen[p.ID] = i
	}
}

func TestAllPins_AllNonZero(t *testing.T) {
	var zero [sha256.Size]byte
	for _, p := range AllPins() {
		if p.SHA == zero {
			t.Errorf("%q zero", p.ID)
		}
	}
}

func TestPinByID_AllCanonical(t *testing.T) {
	for _, expected := range AllPins() {
		got, ok := PinByID(expected.ID)
		if !ok || got.SHA != expected.SHA {
			t.Errorf("drift on %q", expected.ID)
		}
	}
}

func TestPinByID_UnknownID(t *testing.T) {
	if _, ok := PinByID("x"); ok {
		t.Fatal()
	}
}

func TestCorpusPin_HexSHA_64Chars(t *testing.T) {
	for _, p := range AllPins() {
		if got := len(p.HexSHA()); got != 64 {
			t.Errorf("%d", got)
		}
	}
}

func TestCorpusPin_PrefixHex_16Chars(t *testing.T) {
	for _, p := range AllPins() {
		if got := len(p.PrefixHex()); got != 16 {
			t.Errorf("%d", got)
		}
	}
}
