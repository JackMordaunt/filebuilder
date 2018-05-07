package filebuilder

import (
	"os"
	"strings"

	"github.com/spf13/afero"
)

// CompareDirectories creates a diff of any differences found.
// If the directories are not comparable, ok == false.
// If the directories are comparable, ok == true.
func CompareDirectories(fs afero.Fs, left, right string) (difference *Diff, ok bool, err error) {
	var (
		diffs   = &Diff{}
		current string
	)
	walk := func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		path = strings.Replace(path, current, "", 1)
		switch current {
		case left:
			diffs.appendLeft(path)
		case right:
			diffs.appendRight(path)
		}
		return nil
	}
	current = right
	err = afero.Walk(fs, right, walk)
	if err != nil {
		return nil, false, err
	}
	current = left
	err = afero.Walk(fs, left, walk)
	if err != nil {
		return nil, false, err
	}
	if !diffs.IsEmpty() {
		return diffs, false, nil
	}
	return nil, true, nil
}

// Compare two different filesystem objects.
func Compare(left, right afero.Fs) (*Diff, bool, error) {
	diffs := &Diff{}
	walkLeft := func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		diffs.appendLeft(path)
		return nil
	}
	walkRight := func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		diffs.appendRight(path)
		return nil
	}
	var err error
	err = afero.Walk(left, "/", walkLeft)
	if err != nil {
		return nil, false, err
	}
	err = afero.Walk(right, "/", walkRight)
	if err != nil {
		return nil, false, err
	}
	if diffs.IsEmpty() {
		return diffs, true, nil
	}
	return diffs, false, nil
}
