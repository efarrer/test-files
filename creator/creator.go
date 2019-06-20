package creator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// A Root is the root directory of a collection of files and directories used for testing
type Root struct {
	Directory string
}

func (r *Root) Destroy() error {
	return os.RemoveAll(r.Directory)
}

func wantDirectory(path string) bool {
	return strings.HasSuffix(path, "/")
}

// A DirectoryContents is a map of file and directory paths to either nil (for directories) or the file contents (for files)
type DirectoryContents map[string][]byte

// New creates file and directories.
func New(contents DirectoryContents) (*Root, error) {
	tempDir, err := ioutil.TempDir("", "test-files")
	if err != nil {
		return nil, err
	}

	root := &Root{
		Directory: tempDir,
	}
	failed := true
	defer func() {
		if failed {
			root.Destroy()
		}
	}()

	for filepath, contents := range contents {
		if wantDirectory(filepath) {
			if contents != nil {
				return root, fmt.Errorf("Directory %s should have nil for it's contents", filepath)
			}

			err = os.MkdirAll(path.Join(root.Directory, filepath), 0700)
			if err != nil {
				return root, err
			}
		} else {
			if contents == nil {
				return root, fmt.Errorf("File %s should have non-nil contents", filepath)
			}
			dir, _ := path.Split(filepath)

			// Make sure the containing directory exists
			err = os.MkdirAll(path.Join(root.Directory, dir), 0700)
			if err != nil {
				return root, err
			}

			err = ioutil.WriteFile(path.Join(root.Directory, filepath), []byte(contents), 0600)
			if err != nil {
				return root, err
			}
		}
	}

	// We succeeded so we can let the caller destroy the root
	failed = false
	return root, nil
}
