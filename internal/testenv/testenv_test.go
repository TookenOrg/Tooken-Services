package testenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The name is spliced into DROP DATABASE: only a tooken_it_ scratch name may
// ever reach it, never the database the URL points at nor a server one.
func TestScratchDatabaseNames(t *testing.T) {
	for name, accepted := range map[string]bool{
		"tooken_it_e2e":               true,
		"tooken_it_delivery_database": true,
		"tooken_it":                   false,
		"tooken_it_":                  false,
		"tooken":                      false,
		"postgres":                    false,
		"template1":                   false,
		"Tooken_it_e2e":               false,
		"tooken_it_a-b":               false,
		"tooken_it_e2e; DROP TABLE x": false,
		// PostgreSQL truncates identifiers beyond 63 bytes: a longer name
		// would silently drop another database.
		"tooken_it_" + strings.Repeat("a", 53): true,
		"tooken_it_" + strings.Repeat("a", 54): false,
	} {
		if got := databaseName.MatchString(name); got != accepted {
			t.Errorf("%q accepted = %v, want %v", name, got, accepted)
		}
	}
}

func TestBaselinePathExists(t *testing.T) {
	path, err := BaselinePath()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("the baseline must exist: %v", err)
	}
}

func TestModuleRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "internal", "pkg")
	if err := os.MkdirAll(nested, 0o700); err != nil {
		t.Fatal(err)
	}

	for _, dir := range []string{root, nested} {
		got, err := moduleRoot(dir)
		if err != nil || got != root {
			t.Errorf("moduleRoot(%s) = %q, %v; want %q", dir, got, err, root)
		}
	}
}

func TestModuleRootOutsideAModule(t *testing.T) {
	dir := t.TempDir()
	if root, err := moduleRoot(dir); err == nil {
		t.Skipf("the temporary directory sits inside a module (%s)", root)
	}
	_, err := moduleRoot(dir)
	if err == nil || !strings.Contains(err.Error(), dir) {
		t.Fatalf("the error must name the directory searched from, got %v", err)
	}
}
