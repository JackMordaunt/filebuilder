package filebuilder

import (
	"testing"

	"github.com/spf13/afero"
)

func TestCompareDirectories(t *testing.T) {
	tests := []struct {
		desc         string
		left         []Entry
		right        []Entry
		shouldEqual  bool
		expectedDiff *Diff
		wantErr      bool
	}{
		{
			desc: "should equal",
			left: []Entry{
				File{Path: "foo/bar.exe"},
				File{Path: "foo/baz.exe"},
			},
			right: []Entry{
				File{Path: "foo/bar.exe"},
				File{Path: "foo/baz.exe"},
			},
			shouldEqual:  true,
			expectedDiff: nil,
			wantErr:      false,
		},
		{
			desc: "should not equal",
			left: []Entry{
				File{Path: "foo/foo.exe"},
				File{Path: "foo/foobar.exe"},
			},
			right: []Entry{
				File{Path: "foo/bar.exe"},
				File{Path: "foo/baz.exe"},
			},
			shouldEqual: false,
			expectedDiff: &Diff{
				left: map[string]struct{}{
					"foo/foo.exe":    struct{}{},
					"foo/foobar.exe": struct{}{},
					"foo":            struct{}{},
				},
				right: map[string]struct{}{
					"foo/bar.exe": struct{}{},
					"foo/baz.exe": struct{}{},
					"foo":         struct{}{},
				},
			},
			wantErr: false,
		},
		{
			desc:         "directory does not exist",
			left:         []Entry{},
			right:        []Entry{},
			shouldEqual:  false,
			expectedDiff: nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		var err error
		fs := afero.NewMemMapFs()
		_, err = Build(fs, "right", tt.right...)
		if err != nil {
			t.Errorf("[%s] unexpected error creating files: %v",
				tt.desc, err)
		}
		_, err = Build(fs, "left", tt.left...)
		if err != nil {
			t.Errorf("[%s] unexpected error creating files: %v",
				tt.desc, err)
		}
		diff, ok, err := CompareDirectories(fs, "left", "right")
		if tt.wantErr {
			if err == nil {
				t.Errorf("[%s] want error but got nil", tt.desc)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("[%s] unexpected error while comparing: %v",
					tt.desc, err)
			}
		}
		if tt.shouldEqual && !ok {
			t.Errorf("[%s] directories should equal but do not, \ndiffs,\n%+v",
				tt.desc, diff)
		}
		if !tt.shouldEqual && ok {
			t.Errorf("[%s] directories should not equal but do",
				tt.desc)
		}
		if !tt.shouldEqual && !ok {
			if !diff.EqualTo(tt.expectedDiff) {
				t.Errorf("[%s] diffs do not match: \nwant\n%+v\ngot\n%+v",
					tt.desc, tt.expectedDiff, diff)
			}
		}
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		desc         string
		left         Entries
		right        Entries
		shouldEqual  bool
		expectedDiff *Diff
		wantErr      bool
	}{
		{
			desc: "should equal",
			left: []Entry{
				File{Path: "foo/bar.exe"},
				File{Path: "foo/baz.exe"},
			},
			right: []Entry{
				File{Path: "foo/bar.exe"},
				File{Path: "foo/baz.exe"},
			},
			shouldEqual:  true,
			expectedDiff: nil,
			wantErr:      false,
		},
		{
			desc: "should not equal",
			left: []Entry{
				File{Path: "foo/foo.exe"},
				File{Path: "foo/foobar.exe"},
			},
			right: []Entry{
				File{Path: "foo/bar.exe"},
				File{Path: "foo/baz.exe"},
			},
			shouldEqual: false,
			expectedDiff: &Diff{
				left: map[string]struct{}{
					"foo/foo.exe":    struct{}{},
					"foo/foobar.exe": struct{}{},
					"foo":            struct{}{},
				},
				right: map[string]struct{}{
					"foo/bar.exe": struct{}{},
					"foo/baz.exe": struct{}{},
					"foo":         struct{}{},
				},
			},
			wantErr: false,
		},
		{
			desc:         "directory does not exist",
			left:         []Entry{},
			right:        []Entry{},
			shouldEqual:  true,
			expectedDiff: nil,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		var err error
		leftFs := buildFs(t, tt.desc, tt.left)
		rightFs := buildFs(t, tt.desc, tt.right)
		diff, ok, err := Compare(leftFs, rightFs)
		if tt.wantErr && err == nil {
			t.Fatalf("[%s] want error but got nil", tt.desc)
			continue
		}
		if !tt.wantErr && err != nil {
			t.Fatalf("[%s] unexpected error while comparing: %v",
				tt.desc, err)
		}
		if tt.shouldEqual && !ok {
			t.Fatalf("[%s] directories should equal but do not, \ndiffs,\n%+v",
				tt.desc, diff)
		}
		if !tt.shouldEqual && ok {
			t.Fatalf("[%s] directories should not equal but do",
				tt.desc)
		}
		if !tt.shouldEqual && !ok {
			if !diff.EqualTo(tt.expectedDiff) {
				t.Fatalf("[%s] diffs do not match: \nwant\n%+v\ngot\n%+v",
					tt.desc, tt.expectedDiff, diff)
			}
		}
	}
}

func buildFs(t *testing.T, desc string, entry Entry) afero.Fs {
	fs := afero.NewMemMapFs()
	_, err := Build(fs, "", entry)
	if err != nil {
		t.Fatalf("[%s] unexpected error creating files: %v", desc, err)
	}
	return fs
}
