// Package honest implements the cohort R143 LOUD-ONCE-WARNING-FLAG
// discipline for limitless-pensionspath.
//
// 5 honest-defaults surfaces:
//
//  1. LIMITLESS_PENSIONSPATH_FCA_COBS_VERSION_PIN_REQUIRED — FCA COBS
//     (Conduct of Business Sourcebook) is amended ~monthly via Handbook
//     Notices; the local corpus SHA must be cold-verified.
//  2. LIMITLESS_PENSIONSPATH_TPR_FUNDING_CODE_VERSION_PIN_REQUIRED — TPR
//     Defined-Benefit Funding Code of Practice (effective Sept 2024 per
//     Pension Schemes Act 2021 + Occupational Pension Schemes (Funding
//     and Investment Strategy and Amendment) Regs 2024) is a moving
//     target.
//  3. LIMITLESS_PENSIONSPATH_HMRC_ALLOWANCES_FROZEN — HMRC pension
//     allowances are politically volatile: LTA abolished 2024-04-06
//     (replaced by Lump Sum Allowance + Lump Sum and Death Benefit
//     Allowance), Annual Allowance currently £60,000 (was £40,000
//     pre-2023). Hard-coded figures are wrong-confidence-risk.
//  4. LIMITLESS_PENSIONSPATH_NOT_FINANCIAL_ADVICE — every emitted
//     recommendation MUST carry R166 LIABILITY-FOOTER-CONST: regulated
//     financial advice requires FSMA 2000 s.19 authorisation by the FCA.
//  5. LIMITLESS_PENSIONSPATH_REVIEWED_BY_COUNSEL_FALSE — R166 honest-
//     default. Phase-1 templates have NOT been reviewed by qualified
//     financial-services counsel or FCA-authorised advisers.
package honest

import (
	"fmt"
	"io"
	"sync"
)

const LoudOncePrefix = "[LOUD-ONCE-WARNING]"

type Severity string

const (
	SeverityInfo  Severity = "INFO"
	SeverityWarn  Severity = "WARN"
	SeverityError Severity = "ERROR"
)

type Advisory struct {
	Code     string
	Severity Severity
	Message  string
	DocLink  string
}

var canonicalAdvisories = []Advisory{
	{
		Code:     "LIMITLESS_PENSIONSPATH_FCA_COBS_VERSION_PIN_REQUIRED",
		Severity: SeverityError,
		Message:  "FCA COBS (Conduct of Business Sourcebook), particularly COBS 19 (Pensions and Retirement) + COBS 19.1 (DB-to-DC transfer advice), is amended ~monthly via FCA Handbook Notices. The local corpus SHA pinned in internal/pension-rules/ MUST be cold-verified against fca.org.uk-published canonical before any live recommendation. Stale pin = silent enforcement of superseded rules = FCA s.19 violation risk.",
		DocLink:  "SECURITY.md",
	},
	{
		Code:     "LIMITLESS_PENSIONSPATH_TPR_FUNDING_CODE_VERSION_PIN_REQUIRED",
		Severity: SeverityError,
		Message:  "TPR Defined-Benefit Funding Code of Practice (effective September 2024 under Pension Schemes Act 2021 + Occupational Pension Schemes (Funding and Investment Strategy and Amendment) Regulations 2024 SI 2024 No. 462) is the load-bearing scheme-funding instrument. Local corpus pin MUST be cold-verified against thepensionsregulator.gov.uk canonical before any DB scheme valuation determination.",
		DocLink:  "SECURITY.md",
	},
	{
		Code:     "LIMITLESS_PENSIONSPATH_HMRC_ALLOWANCES_FROZEN",
		Severity: SeverityError,
		Message:  "HMRC pension allowances are politically volatile. The Lifetime Allowance (LTA) was abolished from 2024-04-06 and replaced by the Lump Sum Allowance (LSA, currently £268,275) + Lump Sum and Death Benefit Allowance (LSDBA, currently £1,073,100). The Annual Allowance is currently £60,000 (was £40,000 pre-2023-04-06). Money Purchase Annual Allowance (MPAA) is £10,000. Hard-coded figures in code are wrong-confidence-risk — any quoted figure MUST be sourced from internal/pension-rules/ corpus + cold-verified against current HMRC Pensions Tax Manual.",
		DocLink:  "SECURITY.md",
	},
	{
		Code:     "LIMITLESS_PENSIONSPATH_NOT_FINANCIAL_ADVICE",
		Severity: SeverityError,
		Message:  "Every pension recommendation emitted by this software MUST carry the R166 LIABILITY-FOOTER-CONST escape: 'NOT FINANCIAL ADVICE. Regulated financial advice in the UK is restricted by FSMA 2000 s.19 (the General Prohibition) to FCA-authorised firms. Carrying out a regulated activity without authorisation is a criminal offence under FSMA 2000 s.23. Consult an FCA-authorised independent financial adviser before relying on any recommendation.'",
		DocLink:  "SECURITY.md",
	},
	{
		Code:     "LIMITLESS_PENSIONSPATH_REVIEWED_BY_COUNSEL_FALSE",
		Severity: SeverityWarn,
		Message:  "R166 LIABILITY-FOOTER-CONST honest-default. Phase-1 scaffold ships ReviewedByCounsel = false. Placeholder narrative templates + suitability-rule scaffolds have NOT been reviewed by qualified financial-services counsel or FCA-authorised advisers. Operator MUST commission counsel review + flip ReviewedByCounsel to true on its own R145.B sibling branch before any live deployment.",
		DocLink:  "SECURITY.md",
	},
}

var (
	registryMu sync.RWMutex
	registry   = map[string]*sync.Once{}
)

func LoudOnce(adv Advisory, w io.Writer) {
	registryMu.RLock()
	once, ok := registry[adv.Code]
	registryMu.RUnlock()
	if !ok {
		registryMu.Lock()
		once, ok = registry[adv.Code]
		if !ok {
			once = &sync.Once{}
			registry[adv.Code] = once
		}
		registryMu.Unlock()
	}
	once.Do(func() {
		_, _ = fmt.Fprintf(w, "%s %s %s: %s (see %s)\n",
			LoudOncePrefix, adv.Severity, adv.Code, adv.Message, adv.DocLink)
	})
}

func Reset() {
	registryMu.Lock()
	registry = map[string]*sync.Once{}
	registryMu.Unlock()
}

func CanonicalAdvisories() []Advisory {
	out := make([]Advisory, len(canonicalAdvisories))
	copy(out, canonicalAdvisories)
	return out
}

func FindAdvisory(code string) (Advisory, bool) {
	for _, a := range canonicalAdvisories {
		if a.Code == code {
			return a, true
		}
	}
	return Advisory{}, false
}
