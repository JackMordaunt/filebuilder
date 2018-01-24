package filebuilder

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Build creates a file structure based on the entries given, resolved against
// root.
// fs defaults to package "os" (afero.OsFs).
// If package "os" is used without a specified root an error is returned to
// avoid operations on the host OS's root.
func Build(fs afero.Fs, root string, entries ...Entry) (CleanFunc, error) {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	if root == "" {
		root = "/"
	}
	if _, ok := fs.(*afero.OsFs); ok && root == "/" {
		return nil, fmt.Errorf("avoided attempt to operate on host OS root folder")
	}
	cleanup := func() error {
		return fs.RemoveAll(root)
	}
	for _, e := range entries {
		if err := e.Create(fs, root); err != nil {
			return cleanup, err
		}
	}
	return cleanup, nil
}

// CleanFunc removes the files created.
type CleanFunc func() error

// Entry represents a file system entry, typically a file or a directory.
type Entry interface {
	// Create implementations should default to `os` when fs is nil.
	// If base is empty, the entry's path is interpreted as absolute.
	Create(fs afero.Fs, base string) error
}

// File represents a file.
type File struct {
	Path    string
	Content []byte
}

// Create the file at the given path with the given content.
func (f File) Create(fs afero.Fs, base string) error {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	path := filepath.Join(base, f.Path)
	if err := fs.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return errors.Wrap(err, "creating parent directory")
	}
	handle, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer handle.Close()
	if len(f.Content) == 0 {
		return nil
	}
	if _, err := io.Copy(handle, bytes.NewBuffer(f.Content)); err != nil {
		return err
	}
	return nil
}

// Dir represents a folder.
type Dir = Directory

// Directory represents a folder.
type Directory struct {
	Path    string
	Entries []Entry
}

// Create the directory at the given path with the given entries.
func (d Directory) Create(fs afero.Fs, base string) error {
	if fs == nil {
		fs = afero.NewOsFs()
	}
	path := filepath.Join(base, d.Path)
	if err := fs.MkdirAll(path, 0755); err != nil {
		return err
	}
	for _, e := range d.Entries {
		if err := e.Create(fs, path); err != nil {
			return err
		}
	}
	return nil
}
