// Package firewall implements R145.C FIREWALL-TEST-DISCIPLINE pin.
package firewall

import (
	"os"
	"path/filepath"
	"sort"
)

func ExpectedPackages() []string {
	return []string{
		"firewall",
		"honest",
		"legal",
		"manifest",
		"mirrormark",
		"pension-rules",
	}
}

func ExpectedBinaries() []string {
	return []string{
		"limitless-pensionspath",
	}
}

func ScanInternal(repoRoot string) ([]string, error) {
	return scanGoSubtree(filepath.Join(repoRoot, "internal"))
}

func ScanCmd(repoRoot string) ([]string, error) {
	cmdDir := filepath.Join(repoRoot, "cmd")
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(cmdDir, e.Name(), "main.go")); err == nil {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}

func scanGoSubtree(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		hasGo, err := dirHasGoFile(filepath.Join(root, e.Name()))
		if err != nil {
			return nil, err
		}
		if hasGo {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}

func dirHasGoFile(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".go" {
			return true, nil
		}
	}
	return false, nil
}
