package creator_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/efarrer/test-files/creator"
)

func isFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return false
	}

	return !fi.IsDir()
}

func isDirectory(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func TestNew_CreatesDirectories(t *testing.T) {
	dir := "foo/"
	root, err := creator.New(creator.DirectoryContents{
		dir: nil,
	})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	defer root.Destroy()

	if !isDirectory(path.Join(root.Directory, dir)) {
		t.Errorf("Expected %s to be created, but it wasn't", path.Join(root.Directory, dir))
	}
}

func TestNew_FailsIfDiretoryIsPassedFileContent(t *testing.T) {
	dir := "foo/"
	_, err := creator.New(creator.DirectoryContents{
		dir: []byte("dirs don't have content"),
	})

	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestNew_FailsIfFileNotPassedContents(t *testing.T) {
	file := "foo"
	_, err := creator.New(creator.DirectoryContents{
		file: nil,
	})

	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestNew_CreatesFiles(t *testing.T) {
	file := "foo/bar.txt"
	expectedContents := []byte("baz")
	root, err := creator.New(creator.DirectoryContents{
		file: expectedContents,
	})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	defer root.Destroy()

	if !isFile(path.Join(root.Directory, file)) {
		t.Errorf("Expected %s to be created, but it wasn't", path.Join(root.Directory, file))
	}

	contents, err := ioutil.ReadFile(path.Join(root.Directory, file))
	if err != nil {
		t.Errorf("Expected to be able to read %s", path.Join(root.Directory, file))
	}

	if string(contents) != string(expectedContents) {
		t.Errorf("Expected %s to contain %s", path.Join(root.Directory, file), expectedContents)
	}
}

func TestDestroy_CleansUpFilesAndDirectories(t *testing.T) {
	file := "foo/bar.txt"
	dir := "bar/baz/biz/"
	expectedContents := []byte("baz")
	root, err := creator.New(creator.DirectoryContents{
		file: expectedContents,
		dir:  nil,
	})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !isFile(path.Join(root.Directory, file)) {
		t.Errorf("Expected %s to be created, but it wasn't", path.Join(root.Directory, file))
	}
	if !isDirectory(path.Join(root.Directory, dir)) {
		t.Errorf("Expected %s to be created, but it wasn't", path.Join(root.Directory, dir))
	}
	root.Destroy()
	if isFile(path.Join(root.Directory, file)) {
		t.Errorf("Expected %s to be created, but it wasn't", path.Join(root.Directory, file))
	}
	if isDirectory(path.Join(root.Directory, dir)) {
		t.Errorf("Expected %s to be created, but it wasn't", path.Join(root.Directory, dir))
	}
}
