package filebuilder

import "testing"

func TestDiff_EqualTo(t *testing.T) {
	newDiff := func() *Diff {
		return &Diff{
			left: map[string]struct{}{
				"foo/bar.exe": struct{}{},
				"foo":         struct{}{},
			},
			right: map[string]struct{}{
				"foo/baz.exe": struct{}{},
				"bar/baz.exe": struct{}{},
				"foo":         struct{}{},
				"bar":         struct{}{},
			},
		}
	}
	left := newDiff()
	right := newDiff()
	if !left.EqualTo(right) || !right.EqualTo(left) {
		t.Fatalf("[should equal] \nleft %s \nright %v", left, right)
	}
}
