// Package taper computes the HMRC Tapered Annual Allowance — the reduction of a
// high earner's pension annual allowance — as an ESTIMATE, never financial/tax
// advice.
//
// HMRC pension allowances are politically volatile (honest advisory #3; figures
// move at fiscal events), so every Result ships Confidence=Low, the
// CorpusHMRCPensionsTaxManual pin for provenance, the FSMA s.19/s.23 liability
// footer, and a cold-verify-against-the-HMRC-Pensions-Tax-Manual caveat.
//
// Activates the CorpusHMRCPensionsTaxManual pin (defined in pension-rules but
// read by no logic until now) with the first computable rule over it.
package taper

import (
	"github.com/davly/limitless-pensionspath/internal/legal"
	"github.com/davly/limitless-pensionspath/internal/manifest"
	pensionrules "github.com/davly/limitless-pensionspath/internal/pension-rules"
)

// Params are the HMRC tapered-annual-allowance figures, integer pounds.
// ILLUSTRATIVE 2023/24–2024/25 defaults — cold-verify before reliance.
type Params struct {
	StandardAllowance       int // £60,000 (was £40,000 pre-2023-04-06)
	ThresholdIncomeGateway  int // £200,000 — taper applies ONLY if threshold income exceeds this
	AdjustedIncomeThreshold int // £260,000 — reduction starts above this
	Floor                   int // £10,000 — minimum tapered allowance
	ReducePer               int // 2 — reduce £1 of allowance per £2 of adjusted income over the threshold
}

// DefaultParams are the 2023/24–2024/25 figures. Illustrative: a host MUST
// re-source from the corpus + cold-verify against the HMRC Pensions Tax Manual.
func DefaultParams() Params {
	return Params{
		StandardAllowance:       60000,
		ThresholdIncomeGateway:  200000,
		AdjustedIncomeThreshold: 260000,
		Floor:                   10000,
		ReducePer:               2,
	}
}

// Result is the structured, non-authoritative output. Footer + Caveat are always
// populated.
type Result struct {
	ThresholdIncome   int
	AdjustedIncome    int
	StandardAllowance int
	TaperApplies      bool
	Reduction         int // effective reduction (StandardAllowance - TaperedAllowance)
	TaperedAllowance  int
	Confidence        manifest.Confidence
	Jurisdiction      manifest.Jurisdiction
	CorpusPinPrefix   string
	Footer            string
	Caveat            string
}

const caveat = "ESTIMATE ONLY, not financial/tax advice. HMRC pension allowances change at fiscal events; cold-verify threshold income, adjusted income, and every figure against the HMRC Pensions Tax Manual for the relevant tax year before reliance."

func corpusPinPrefix() string {
	if p, ok := pensionrules.PinByID(pensionrules.CorpusHMRCPensionsTaxManual); ok {
		return p.PrefixHex()
	}
	return ""
}

// TaperedAnnualAllowance computes the HMRC tapered annual allowance from threshold
// income and adjusted income (integer pounds). The taper applies ONLY when
// threshold income exceeds the gateway (£200,000) AND adjusted income exceeds the
// threshold (£260,000); then the allowance reduces £1 per £2 of adjusted income
// over the threshold, down to the floor (£10,000).
func TaperedAnnualAllowance(thresholdIncome, adjustedIncome int, p Params) Result {
	r := Result{
		ThresholdIncome:   thresholdIncome,
		AdjustedIncome:    adjustedIncome,
		StandardAllowance: p.StandardAllowance,
		TaperedAllowance:  p.StandardAllowance,
		Confidence:        manifest.ConfidenceLow,
		Jurisdiction:      manifest.JurisdictionUK,
		CorpusPinPrefix:   corpusPinPrefix(),
		Footer:            legal.LegalLiabilityFooter,
		Caveat:            caveat,
	}
	// Two-gate test: the taper applies only when BOTH are exceeded. A high
	// adjusted income with threshold income at/under the gateway -> no taper.
	if thresholdIncome <= p.ThresholdIncomeGateway || adjustedIncome <= p.AdjustedIncomeThreshold {
		return r
	}
	r.TaperApplies = true
	over := adjustedIncome - p.AdjustedIncomeThreshold
	tapered := p.StandardAllowance - over/p.ReducePer
	if tapered < p.Floor {
		tapered = p.Floor
	}
	r.TaperedAllowance = tapered
	r.Reduction = p.StandardAllowance - tapered
	return r
}
