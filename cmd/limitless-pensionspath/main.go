// Command limitless-pensionspath — UK FCA + TPR + HMRC pension rules
// compliance forge CLI.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/davly/limitless-pensionspath/internal/honest"
	"github.com/davly/limitless-pensionspath/internal/manifest"
	pensionrules "github.com/davly/limitless-pensionspath/internal/pension-rules"
	"github.com/davly/limitless-pensionspath/internal/taper"
)

const version = "0.1.0-i51-scaffold"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "corpus":
		corpusCmd(os.Args[2:])
	case "advisories":
		advisoriesCmd(os.Args[2:])
	case "manifest":
		manifestCmd(os.Args[2:])
	case "taper":
		taperCmd(os.Args[2:])
	case "version":
		fmt.Println("limitless-pensionspath", version)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: limitless-pensionspath <command>

Commands:
  corpus list     -- list pinned pension-rules corpus SHAs
  advisories list -- list R143 advisories
  manifest list   -- list R150 schematised-knowledge entries
  taper <threshold-income> <adjusted-income> -- HMRC tapered annual-allowance ESTIMATE (not advice)
  version         -- print version`)
}

func taperCmd(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: taper <threshold-income-pounds> <adjusted-income-pounds>")
		os.Exit(2)
	}
	ti, err1 := strconv.Atoi(args[0])
	ai, err2 := strconv.Atoi(args[1])
	if err1 != nil || err2 != nil {
		fmt.Fprintln(os.Stderr, "incomes must be integer pounds, e.g. taper 210000 280000")
		os.Exit(2)
	}
	r := taper.TaperedAnnualAllowance(ti, ai, taper.DefaultParams())
	fmt.Println("Tapered Annual Allowance ESTIMATE")
	fmt.Printf("  threshold income: GBP %d | adjusted income: GBP %d | standard allowance: GBP %d\n",
		r.ThresholdIncome, r.AdjustedIncome, r.StandardAllowance)
	if r.TaperApplies {
		fmt.Printf("  taper APPLIES: reduction GBP %d -> tapered allowance GBP %d\n", r.Reduction, r.TaperedAllowance)
	} else {
		fmt.Printf("  taper does not apply -> full allowance GBP %d\n", r.TaperedAllowance)
	}
	fmt.Printf("  confidence: Low (illustrative defaults) | jurisdiction: %s | corpus pin: %s\n",
		r.Jurisdiction, r.CorpusPinPrefix)
	fmt.Printf("  %s\n", r.Caveat)
	fmt.Printf("  %s\n", r.Footer)
}

func corpusCmd(args []string) {
	if len(args) < 1 || args[0] != "list" {
		fmt.Fprintln(os.Stderr, "Usage: corpus list")
		os.Exit(2)
	}
	for _, p := range pensionrules.AllPins() {
		fmt.Printf("%s\n  sha256: %s\n  prefix: %s\n", p.ID, p.HexSHA(), p.PrefixHex())
	}
}

func advisoriesCmd(args []string) {
	if len(args) < 1 || args[0] != "list" {
		fmt.Fprintln(os.Stderr, "Usage: advisories list")
		os.Exit(2)
	}
	for _, a := range honest.CanonicalAdvisories() {
		fmt.Printf("[%s] %s\n  %s\n", a.Severity, a.Code, a.Message)
	}
}

func manifestCmd(args []string) {
	if len(args) < 1 || args[0] != "list" {
		fmt.Fprintln(os.Stderr, "Usage: manifest list")
		os.Exit(2)
	}
	for _, e := range manifest.Seed() {
		fmt.Printf("%s\n  desc: %s\n  source: %s\n  jurisdiction: %s\n  version: %s\n",
			e.Key, e.Description, e.Source, e.Jurisdiction, e.Version)
	}
}
