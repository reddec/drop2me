package main

import (
	"path/filepath"
	"testing"
)

func TestCreateUniqueFile(t *testing.T) {
	dir := t.TempDir()

	f, dst, err := createUniqueFile(dir, "file.txt")
	if err != nil {
		t.Fatalf("no collision: %v", err)
	}
	f.Close()
	if dst != filepath.Join(dir, "file.txt") {
		t.Fatalf("no collision: got %q, want %q", dst, filepath.Join(dir, "file.txt"))
	}

	f, dst, err = createUniqueFile(dir, "file.txt")
	if err != nil {
		t.Fatalf("first collision: %v", err)
	}
	f.Close()
	if dst != filepath.Join(dir, "file_1.txt") {
		t.Fatalf("first collision: got %q, want %q", dst, filepath.Join(dir, "file_1.txt"))
	}

	f, dst, err = createUniqueFile(dir, "file.txt")
	if err != nil {
		t.Fatalf("second collision: %v", err)
	}
	f.Close()
	if dst != filepath.Join(dir, "file_2.txt") {
		t.Fatalf("second collision: got %q, want %q", dst, filepath.Join(dir, "file_2.txt"))
	}

	f, dst, err = createUniqueFile(dir, "readme")
	if err != nil {
		t.Fatalf("no-ext: %v", err)
	}
	f.Close()
	if dst != filepath.Join(dir, "readme") {
		t.Fatalf("no-ext: got %q, want %q", dst, filepath.Join(dir, "readme"))
	}

	f, dst, err = createUniqueFile(dir, "readme")
	if err != nil {
		t.Fatalf("no-ext collision: %v", err)
	}
	f.Close()
	if dst != filepath.Join(dir, "readme_1") {
		t.Fatalf("no-ext collision: got %q, want %q", dst, filepath.Join(dir, "readme_1"))
	}

	f, dst, err = createUniqueFile(dir, "other.png")
	if err != nil {
		t.Fatalf("free name: %v", err)
	}
	f.Close()
	if dst != filepath.Join(dir, "other.png") {
		t.Fatalf("free name: got %q, want %q", dst, filepath.Join(dir, "other.png"))
	}

	_, _, err = createUniqueFile("/nonexistent/path", "foo.txt")
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}
}
