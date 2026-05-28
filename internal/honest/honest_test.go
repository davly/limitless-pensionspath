package honest

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoudOncePrefix(t *testing.T) {
	if LoudOncePrefix != "[LOUD-ONCE-WARNING]" {
		t.Fatalf("drift: %q", LoudOncePrefix)
	}
}

func TestLoudOnce_EmitsOnFirstCall(t *testing.T) {
	Reset()
	var buf bytes.Buffer
	LoudOnce(Advisory{Code: "T_A", Severity: SeverityInfo, Message: "m", DocLink: "d"}, &buf)
	if !strings.Contains(buf.String(), "T_A") {
		t.Errorf("missing: %q", buf.String())
	}
}

func TestLoudOnce_SilentOnSubsequent(t *testing.T) {
	Reset()
	var buf bytes.Buffer
	adv := Advisory{Code: "ONCE", Severity: SeverityInfo, Message: "m", DocLink: "d"}
	LoudOnce(adv, &buf)
	buf.Reset()
	LoudOnce(adv, &buf)
	LoudOnce(adv, &buf)
	if buf.Len() != 0 {
		t.Fatalf("leaked: %q", buf.String())
	}
}

func TestCanonicalAdvisories_Count5(t *testing.T) {
	if len(CanonicalAdvisories()) != 5 {
		t.Fatalf("got %d", len(CanonicalAdvisories()))
	}
}

func TestCanonicalAdvisories_AllFieldsNonEmpty(t *testing.T) {
	for i, a := range CanonicalAdvisories() {
		if a.Code == "" || a.Severity == "" || a.Message == "" || a.DocLink == "" {
			t.Errorf("%d empty", i)
		}
	}
}

func TestCanonicalAdvisories_UniqueCodes(t *testing.T) {
	seen := map[string]int{}
	for i, a := range CanonicalAdvisories() {
		if prev, ok := seen[a.Code]; ok {
			t.Errorf("dup %q at %d and %d", a.Code, prev, i)
		}
		seen[a.Code] = i
	}
}

func TestCanonicalAdvisories_AllStartWithPrefix(t *testing.T) {
	for _, a := range CanonicalAdvisories() {
		if !strings.HasPrefix(a.Code, "LIMITLESS_PENSIONSPATH_") {
			t.Errorf("missing prefix: %q", a.Code)
		}
	}
}

func TestCanonicalAdvisories_R143A_Ladder_4Error_1Warn(t *testing.T) {
	var nE, nW int
	for _, a := range CanonicalAdvisories() {
		switch a.Severity {
		case SeverityError:
			nE++
		case SeverityWarn:
			nW++
		}
	}
	// Pensions Phase-1 ladder: 4 Error (FCA + TPR + HMRC + not-advice are all liability-bearing) + 1 Warn (counsel-review-false)
	if nE != 4 || nW != 1 {
		t.Fatalf("ladder drift: %dE %dW (want 4E 1W for pensions Phase-1)", nE, nW)
	}
}

func TestFindAdvisory_ByCanonicalCode(t *testing.T) {
	for _, expected := range CanonicalAdvisories() {
		got, ok := FindAdvisory(expected.Code)
		if !ok || got.Code != expected.Code {
			t.Errorf("drift on %q", expected.Code)
		}
	}
}

func TestFindAdvisory_UnknownCode(t *testing.T) {
	if _, ok := FindAdvisory("X"); ok {
		t.Fatal("ok=true")
	}
}

func TestReset_ClearsRegistry(t *testing.T) {
	Reset()
	var buf bytes.Buffer
	adv := Advisory{Code: "RST", Severity: SeverityInfo, Message: "m", DocLink: "d"}
	LoudOnce(adv, &buf)
	first := buf.String()
	buf.Reset()
	LoudOnce(adv, &buf)
	if buf.Len() != 0 {
		t.Fatal("expected silent")
	}
	Reset()
	LoudOnce(adv, &buf)
	if buf.String() != first {
		t.Error("drift")
	}
}
