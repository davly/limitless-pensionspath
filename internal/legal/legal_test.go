package legal

import (
	"strings"
	"testing"
)

func TestLegalLiabilityFooter_NonEmpty(t *testing.T) {
	if LegalLiabilityFooter == "" {
		t.Fatal()
	}
}

func TestLegalLiabilityFooter_ContainsFSMA(t *testing.T) {
	if !strings.Contains(LegalLiabilityFooter, "FSMA 2000") {
		t.Fatal("missing FSMA 2000")
	}
}

func TestLegalLiabilityFooter_ContainsFCA(t *testing.T) {
	if !strings.Contains(LegalLiabilityFooter, "FCA") {
		t.Fatal("missing FCA")
	}
}

func TestLegalLiabilityFooter_ContainsNotFinancialAdvice(t *testing.T) {
	if !strings.Contains(LegalLiabilityFooter, "NOT FINANCIAL ADVICE") {
		t.Fatal("missing escape phrase")
	}
}

func TestReviewedByCounsel_HonestDefaultFalse(t *testing.T) {
	if ReviewedByCounsel {
		t.Fatal("must be false")
	}
}

func TestLibraryRecommendsHostActs_HonestDefaultFalse(t *testing.T) {
	if LibraryRecommendsHostActs {
		t.Fatal("must be false")
	}
}
