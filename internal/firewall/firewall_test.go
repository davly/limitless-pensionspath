package firewall

import (
	"path/filepath"
	"runtime"
	"sort"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal()
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestExpectedPackages_NonEmpty(t *testing.T) {
	if len(ExpectedPackages()) == 0 {
		t.Fatal()
	}
}

func TestExpectedBinaries_NonEmpty(t *testing.T) {
	if len(ExpectedBinaries()) == 0 {
		t.Fatal()
	}
}

func TestExpectedPackages_Sorted(t *testing.T) {
	if !sort.StringsAreSorted(ExpectedPackages()) {
		t.Fatal()
	}
}

func TestExpectedPackages_Unique(t *testing.T) {
	seen := map[string]int{}
	for i, p := range ExpectedPackages() {
		if prev, ok := seen[p]; ok {
			t.Errorf("dup %q %d %d", p, prev, i)
		}
		seen[p] = i
	}
}

func TestExpectedPackages_PinnedCount(t *testing.T) {
	const expected = 7 // +taper (HMRC tapered annual-allowance calc)
	if got := len(ExpectedPackages()); got != expected {
		t.Fatalf("got %d, want %d", got, expected)
	}
}

func TestFirewall_EveryExpectedPackageExistsOnDisk(t *testing.T) {
	root := repoRoot(t)
	onDisk, err := ScanInternal(root)
	if err != nil {
		t.Fatalf("%v", err)
	}
	set := map[string]bool{}
	for _, p := range onDisk {
		set[p] = true
	}
	for _, expected := range ExpectedPackages() {
		if !set[expected] {
			t.Errorf("missing %q", expected)
		}
	}
}

func TestFirewall_EveryOnDiskPackageInExpectedList(t *testing.T) {
	root := repoRoot(t)
	onDisk, err := ScanInternal(root)
	if err != nil {
		t.Fatalf("%v", err)
	}
	set := map[string]bool{}
	for _, p := range ExpectedPackages() {
		set[p] = true
	}
	for _, found := range onDisk {
		if !set[found] {
			t.Errorf("unexpected %q", found)
		}
	}
}

func TestFirewall_EveryExpectedBinaryExistsOnDisk(t *testing.T) {
	root := repoRoot(t)
	onDisk, err := ScanCmd(root)
	if err != nil {
		t.Fatalf("%v", err)
	}
	set := map[string]bool{}
	for _, b := range onDisk {
		set[b] = true
	}
	for _, expected := range ExpectedBinaries() {
		if !set[expected] {
			t.Errorf("missing binary %q", expected)
		}
	}
}
