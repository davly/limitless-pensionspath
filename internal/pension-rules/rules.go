// Package pensionrules holds the THE-MOAT corpus pin for UK pension
// rules — FCA COBS Chapter 19 (pensions + retirement), TPR DB Funding
// Code of Practice, HMRC Pensions Tax Manual.
//
// PHASE-1 SCAFFOLD NOTE — SHA values are PLACEHOLDER. Phase-2 binds
// them to canonical regulator-published artefacts.
package pensionrules

import (
	"crypto/sha256"
	"encoding/hex"
)

type CorpusID string

const (
	CorpusFCACOBS19            CorpusID = "uk_fca_cobs_chapter_19_2024_handbook_notice"
	CorpusTPRDBFundingCode2024 CorpusID = "uk_tpr_db_funding_code_2024_09"
	CorpusHMRCPensionsTaxManual CorpusID = "uk_hmrc_pensions_tax_manual_2024_25"
)

var FCACOBS19SHA = sha256.Sum256([]byte(string(CorpusFCACOBS19)))
var TPRDBFundingCode2024SHA = sha256.Sum256([]byte(string(CorpusTPRDBFundingCode2024)))
var HMRCPensionsTaxManualSHA = sha256.Sum256([]byte(string(CorpusHMRCPensionsTaxManual)))

type CorpusPin struct {
	ID  CorpusID
	SHA [sha256.Size]byte
}

func (p CorpusPin) HexSHA() string {
	return hex.EncodeToString(p.SHA[:])
}

func (p CorpusPin) PrefixHex() string {
	return hex.EncodeToString(p.SHA[:8])
}

func AllPins() []CorpusPin {
	return []CorpusPin{
		{ID: CorpusFCACOBS19, SHA: FCACOBS19SHA},
		{ID: CorpusHMRCPensionsTaxManual, SHA: HMRCPensionsTaxManualSHA},
		{ID: CorpusTPRDBFundingCode2024, SHA: TPRDBFundingCode2024SHA},
	}
}

func PinByID(id CorpusID) (CorpusPin, bool) {
	for _, p := range AllPins() {
		if p.ID == id {
			return p, true
		}
	}
	return CorpusPin{}, false
}
