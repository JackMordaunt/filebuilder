package filebuilder

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestBuild(t *testing.T) {
	// TODO(jackmordaunt) test for correct content of files, not just existence.
	tests := []struct {
		desc    string
		fs      afero.Fs
		root    string
		entries []Entry
		want    []string
		wantErr bool
	}{
		{
			desc: "unspecified root",
			fs:   afero.NewMemMapFs(),
			root: "",
			entries: []Entry{
				Directory{Path: "foo", Entries: []Entry{
					File{Path: "bar.exe"},
					File{Path: "baz.exe"},
				}},
				File{Path: "bar.exe"},
				File{Path: "baz.exe"},
			},
			want: []string{
				"/foo/bar.exe",
				"/foo/baz.exe",
				"/bar.exe",
				"/baz.exe",
			},
			wantErr: false,
		},
		{
			desc: "specified root",
			fs:   afero.NewMemMapFs(),
			root: "root/path/here",
			entries: []Entry{
				Directory{Path: "foo", Entries: []Entry{
					File{Path: "bar.exe"},
					File{Path: "baz.exe"},
				}},
				File{Path: "bar.exe"},
				File{Path: "baz.exe"},
			},
			want: []string{
				"root/path/here/foo/bar.exe",
				"root/path/here/foo/baz.exe",
				"root/path/here/bar.exe",
				"root/path/here/baz.exe",
			},
			wantErr: false,
		},
		{
			desc: "flat list",
			fs:   afero.NewMemMapFs(),
			root: "",
			entries: []Entry{
				File{Path: "bar.exe"},
				File{Path: "baz.exe"},
				File{Path: "foo/bar/baz.exe"},
				File{Path: "foo/baz/baz.exe"},
				File{Path: "foo/baz.exe"},
				Directory{Path: "bar"},
			},
			want: []string{
				"/bar.exe",
				"/baz.exe",
				"/foo/bar/baz.exe",
				"/foo/baz/baz.exe",
				"/foo/baz.exe",
				"/foo/bar",
				"/foo/baz",
				"/bar",
			},
			wantErr: false,
		},
		{
			desc: "unspecfied root and filesystem",
			fs:   afero.NewOsFs(),
			root: "",
			entries: []Entry{
				File{Path: "bar.exe"},
				File{Path: "baz.exe"},
				File{Path: "foo/bar/baz.exe"},
				File{Path: "foo/baz/baz.exe"},
				File{Path: "foo/baz.exe"},
				Directory{Path: "bar"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			desc: "zip file",
			fs:   afero.NewMemMapFs(),
			root: "",
			entries: []Entry{
				Zip{
					Path: "bar.zip",
					Files: []File{
						File{Path: "baz.exe"},
						File{Path: "foo/bar/baz.exe"},
						File{Path: "foo/baz/baz.exe"},
						File{Path: "foo/baz.exe"},
					},
				},
			},
			want:    []string{"/bar.zip"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cleanup, err := Build(tt.fs, tt.root, tt.entries...)
		if tt.wantErr {
			if err == nil {
				t.Errorf("[%s] wanted error, got nil", tt.desc)
			}
			continue
		}
		if err != nil {
			t.Errorf("[%s] unexpected error while creating files: %v",
				tt.desc, err)
		}
		for _, path := range tt.want {
			_, err := tt.fs.Stat(path)
			if err != nil {
				t.Errorf("[%s] entry '%s' does not exist",
					tt.desc, path)
			} else if err != nil && !os.IsNotExist(err) {
				t.Errorf("[%s] unexpected error while reading entry: %v",
					tt.desc, path)
			}
		}
		if err := cleanup(); err != nil {
			t.Fatalf("[%s] unexpected error during cleanup: %v", tt.desc, err)
		}
		for _, path := range tt.want {
			_, err := tt.fs.Stat(path)
			if err == nil {
				t.Errorf("[%s] file '%s' was not cleaned up", tt.desc, path)
			}
		}
	}
}
