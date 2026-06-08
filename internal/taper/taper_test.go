package taper

import (
	"testing"

	"github.com/davly/limitless-pensionspath/internal/legal"
	"github.com/davly/limitless-pensionspath/internal/manifest"
)

func TestTaperBoundaries(t *testing.T) {
	p := DefaultParams()
	cases := []struct {
		name             string
		threshold, adj   int
		wantApplies      bool
		wantTapered      int
	}{
		{"below gateway (high adj, low threshold)", 150_000, 300_000, false, 60_000},
		{"gateway met, adj at threshold", 210_000, 260_000, false, 60_000},
		{"standard taper", 210_000, 280_000, true, 50_000},  // over 20k /2 = 10k reduction
		{"at floor", 210_000, 360_000, true, 10_000},        // over 100k /2 = 50k -> 60k-50k=10k
		{"beyond floor (capped)", 210_000, 400_000, true, 10_000}, // over 140k /2 = 70k -> floored
	}
	for _, c := range cases {
		r := TaperedAnnualAllowance(c.threshold, c.adj, p)
		if r.TaperApplies != c.wantApplies || r.TaperedAllowance != c.wantTapered {
			t.Errorf("%s: got (applies=%v, tapered=%d), want (applies=%v, tapered=%d)",
				c.name, r.TaperApplies, r.TaperedAllowance, c.wantApplies, c.wantTapered)
		}
		// reduction must reconcile with the allowance
		if r.Reduction != r.StandardAllowance-r.TaperedAllowance {
			t.Errorf("%s: reduction %d does not reconcile", c.name, r.Reduction)
		}
	}
}

func TestGatewayProtectsLowThresholdIncome(t *testing.T) {
	// The defining subtlety: a £300k adjusted income does NOT taper if threshold
	// income is at/under £200k (e.g. large employer/salary-sacrifice contributions).
	r := TaperedAnnualAllowance(200_000, 300_000, DefaultParams())
	if r.TaperApplies || r.TaperedAllowance != 60_000 {
		t.Fatalf("gateway not met must keep full allowance: applies=%v tapered=%d", r.TaperApplies, r.TaperedAllowance)
	}
}

func TestHonestEnvelopeAlways(t *testing.T) {
	for _, r := range []Result{
		TaperedAnnualAllowance(150_000, 100_000, DefaultParams()), // no taper
		TaperedAnnualAllowance(300_000, 350_000, DefaultParams()), // tapered
	} {
		if r.Footer != legal.LegalLiabilityFooter {
			t.Error("FSMA liability footer must always be present")
		}
		if r.Confidence != manifest.ConfidenceLow {
			t.Errorf("estimate must be Confidence=Low, got %v", r.Confidence)
		}
		if r.Jurisdiction != manifest.JurisdictionUK {
			t.Errorf("jurisdiction must be UK, got %q", r.Jurisdiction)
		}
		if r.CorpusPinPrefix == "" {
			t.Error("result must carry the HMRC Pensions Tax Manual pin prefix (the dead pin is now load-bearing)")
		}
		if r.Caveat == "" {
			t.Error("cold-verify caveat must always be present")
		}
	}
}

func TestMonotonicNonIncreasing(t *testing.T) {
	// Above the gateway, more adjusted income never increases the allowance.
	p := DefaultParams()
	prev := p.StandardAllowance + 1
	for adj := 260_000; adj <= 400_000; adj += 5_000 {
		got := TaperedAnnualAllowance(250_000, adj, p).TaperedAllowance
		if got > prev {
			t.Fatalf("allowance increased with income at adj=%d: %d > %d", adj, got, prev)
		}
		if got < p.Floor {
			t.Fatalf("allowance fell below floor at adj=%d: %d", adj, got)
		}
		prev = got
	}
}
