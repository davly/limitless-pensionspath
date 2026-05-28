// Package manifest implements R150 PARALLEL-MAP review-metadata envelope
// for limitless-pensionspath, extended with the R150 Class-3 jurisdiction-
// version anchor for UK FCA COBS + TPR Funding Code + HMRC Pensions Tax
// Manual corpus pinning.
package manifest

import (
	"sort"
	"time"
)

const SchemaVersion = 1

var FreshAtUnknown = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

const (
	SourceFCAHandbookCOBS19    = "FCA Handbook, Conduct of Business Sourcebook (COBS) Chapter 19 — Pensions and Retirement (fca.org.uk/handbook/COBS/19)"
	SourceTPRFundingCode       = "TPR DB Funding Code of Practice (effective September 2024) — thepensionsregulator.gov.uk + Pension Schemes Act 2021"
	SourceHMRCPensionsTaxManual = "HMRC Pensions Tax Manual (PTM) — gov.uk/hmrc-internal-manuals/pensions-tax-manual"
	SourceFSMA2000              = "Financial Services and Markets Act 2000 (s.19 General Prohibition + s.23 criminal offence)"
	SourcePensionsAct2008       = "Pensions Act 2008 (auto-enrolment + employer duties)"
	SourceMethodologyCorpusPkg  = "limitless-pensionspath internal/pension-rules package"
	SourceContextDoc            = "limitless-pensionspath CONTEXT.md"
	SourceR85ParityMarker       = "limitless-pensionspath R85 CLEAN-PARITY"
)

type Confidence int

const (
	ConfidenceHigh   Confidence = 3
	ConfidenceMedium Confidence = 2
	ConfidenceLow    Confidence = 1
)

type Jurisdiction string

const (
	JurisdictionUK   Jurisdiction = "UK"
	JurisdictionNone Jurisdiction = ""
)

type Entry struct {
	Key           string
	Description   string
	FreshAt       time.Time
	Source        string
	SchemaVersion int
	Confidence    Confidence
	Jurisdiction  Jurisdiction
	Version       string
}

func (e Entry) IsStale(now time.Time, maxAge time.Duration) bool {
	if e.FreshAt.Equal(FreshAtUnknown) {
		return true
	}
	return now.Sub(e.FreshAt) > maxAge
}

type Manifest []Entry

func (m Manifest) SortedKeys() []string {
	keys := make([]string, 0, len(m))
	for _, e := range m {
		keys = append(keys, e.Key)
	}
	sort.Strings(keys)
	return keys
}

func (m Manifest) StaleEntries(now time.Time, maxAge time.Duration) []Entry {
	var out []Entry
	for _, e := range m {
		if e.IsStale(now, maxAge) {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

func AllSources() []string {
	return []string{
		SourceContextDoc,
		SourceFCAHandbookCOBS19,
		SourceFSMA2000,
		SourceHMRCPensionsTaxManual,
		SourceMethodologyCorpusPkg,
		SourcePensionsAct2008,
		SourceR85ParityMarker,
		SourceTPRFundingCode,
	}
}

func Seed() Manifest {
	scaffold := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC)
	tbdPhase2 := FreshAtUnknown
	return Manifest{
		{
			Key:           "corpus.uk.fca_cobs_chapter_19",
			Description:   "FCA Handbook COBS Chapter 19 (Pensions and Retirement) corpus pin. Amended ~monthly via Handbook Notices. Local SHA pin awaiting Phase-2 cold-verify.",
			FreshAt:       tbdPhase2,
			Source:        SourceFCAHandbookCOBS19,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceLow,
			Jurisdiction:  JurisdictionUK,
			Version:       "as-amended-2024-handbook-notice",
		},
		{
			Key:           "corpus.uk.tpr_db_funding_code_2024",
			Description:   "TPR Defined-Benefit Funding Code of Practice (effective September 2024). Pension Schemes Act 2021 + Occupational Pension Schemes (Funding and Investment Strategy and Amendment) Regulations 2024 (SI 2024/462).",
			FreshAt:       tbdPhase2,
			Source:        SourceTPRFundingCode,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceLow,
			Jurisdiction:  JurisdictionUK,
			Version:       "2024-09",
		},
		{
			Key:           "corpus.uk.hmrc_pensions_tax_manual",
			Description:   "HMRC Pensions Tax Manual (PTM) corpus pin. LTA abolished 2024-04-06; LSA £268,275 + LSDBA £1,073,100; AA £60,000; MPAA £10,000.",
			FreshAt:       tbdPhase2,
			Source:        SourceHMRCPensionsTaxManual,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceLow,
			Jurisdiction:  JurisdictionUK,
			Version:       "2024-25-tax-year",
		},
		{
			Key:           "regulation.uk.fsma_2000_s19",
			Description:   "FSMA 2000 s.19 General Prohibition — carrying out a regulated activity without authorisation is a criminal offence under s.23.",
			FreshAt:       scaffold,
			Source:        SourceFSMA2000,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionUK,
			Version:       "2000-c.8",
		},
		{
			Key:           "regulation.uk.pensions_act_2008",
			Description:   "Pensions Act 2008 — auto-enrolment + employer duties (Chapters 1-3 of Part 1).",
			FreshAt:       scaffold,
			Source:        SourcePensionsAct2008,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionUK,
			Version:       "2008-c.30",
		},
		{
			Key:           "regulation.uk.pension_schemes_act_2021",
			Description:   "Pension Schemes Act 2021 — DB funding regulation framework, master trust authorisation, GMP equalisation.",
			FreshAt:       scaffold,
			Source:        SourceTPRFundingCode,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionUK,
			Version:       "2021-c.1",
		},
		{
			Key:           "cohort.l43.mirrormark_v1",
			Description:   "L43 Mirror-Mark v1 receipt algorithm byte-identical to foundation/pkg/mirrormark.",
			FreshAt:       scaffold,
			Source:        SourceMethodologyCorpusPkg,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionNone,
			Version:       "v1",
		},
		{
			Key:           "cohort.r151.kat1_canonical_hex",
			Description:   "R151 KAT-1 cross-substrate hex anchor: 239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca.",
			FreshAt:       scaffold,
			Source:        SourceMethodologyCorpusPkg,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionNone,
			Version:       "v1",
		},
		{
			Key:           "placeholder.counsel_review_status",
			Description:   "R166 LIABILITY-FOOTER-CONST honest-default: ReviewedByCounsel = false. Phase-1 templates have NOT been reviewed by qualified financial-services counsel or FCA-authorised advisers.",
			FreshAt:       scaffold,
			Source:        SourceContextDoc,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceLow,
			Jurisdiction:  JurisdictionNone,
			Version:       "phase-1",
		},
		{
			Key:           "r85.parity.code_vs_context",
			Description:   "R85 CLEAN-PARITY anchor — CONTEXT.md status row vs runtime ground truth.",
			FreshAt:       scaffold,
			Source:        SourceR85ParityMarker,
			SchemaVersion: SchemaVersion,
			Confidence:    ConfidenceHigh,
			Jurisdiction:  JurisdictionNone,
			Version:       "v1",
		},
	}
}
