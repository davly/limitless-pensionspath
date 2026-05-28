// Command limitless-pensionspath — UK FCA + TPR + HMRC pension rules
// compliance forge CLI.
package main

import (
	"fmt"
	"os"

	"github.com/davly/limitless-pensionspath/internal/honest"
	"github.com/davly/limitless-pensionspath/internal/manifest"
	pensionrules "github.com/davly/limitless-pensionspath/internal/pension-rules"
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
  version         -- print version`)
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
