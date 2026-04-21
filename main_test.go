package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUniquePath(t *testing.T) {
	dir := t.TempDir()

	// No collision — returns the original path.
	got := uniquePath(dir, "file.txt")
	want := filepath.Join(dir, "file.txt")
	if got != want {
		t.Fatalf("no collision: got %q, want %q", got, want)
	}

	// Create the file so the next call must pick a different name.
	touch(t, got)

	got = uniquePath(dir, "file.txt")
	want = filepath.Join(dir, "file_1.txt")
	if got != want {
		t.Fatalf("first collision: got %q, want %q", got, want)
	}

	// Create _1 too — expect _2.
	touch(t, got)

	got = uniquePath(dir, "file.txt")
	want = filepath.Join(dir, "file_2.txt")
	if got != want {
		t.Fatalf("second collision: got %q, want %q", got, want)
	}

	// File without extension.
	touch(t, filepath.Join(dir, "readme"))
	got = uniquePath(dir, "readme")
	want = filepath.Join(dir, "readme_1")
	if got != want {
		t.Fatalf("no-ext collision: got %q, want %q", got, want)
	}

	// Completely free name in a busy directory — no suffix added.
	got = uniquePath(dir, "other.png")
	want = filepath.Join(dir, "other.png")
	if got != want {
		t.Fatalf("free name: got %q, want %q", got, want)
	}
}

func touch(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("touch %q: %v", path, err)
	}
	f.Close()
}
