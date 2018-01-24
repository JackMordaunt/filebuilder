package filebuilder

import (
	"bytes"
	"fmt"
	"strings"
)

// Diff is used to record and compute the differences between two sets of file
// paths.
// Generally speaking any two string slices can be used, however this object is
// designed specifically for lists of file paths.
type Diff struct {
	left    map[string]struct{}
	right   map[string]struct{}
	diff    map[string][]string
	changed bool
}

func (d *Diff) appendLeft(path string) {
	if path == "" {
		return
	}
	if strings.HasPrefix(path, "/") {
		path = path[1:len(path)]
	}
	if d.left == nil {
		d.left = map[string]struct{}{}
	}
	d.left[path] = struct{}{}
	d.changed = true
}

func (d *Diff) appendRight(path string) {
	if path == "" {
		return
	}
	if strings.HasPrefix(path, "/") {
		path = path[1:len(path)]
	}
	if d.right == nil {
		d.right = map[string]struct{}{}
	}
	d.right[path] = struct{}{}
	d.changed = true
}

// Diff returns a map containing the contents of each list that is not
// in the other.
func (d *Diff) Diff() map[string][]string {
	defer func() {
		d.changed = false
	}()
	if d.diff == nil || d.changed {
		d.diff = map[string][]string{
			"left":  []string{},
			"right": []string{},
		}
		for l := range d.left {
			if _, ok := d.right[l]; !ok {
				d.diff["left"] = append(d.diff["left"], l)
			}
		}
		for r := range d.right {
			if _, ok := d.left[r]; !ok {
				d.diff["right"] = append(d.diff["right"], r)
			}
		}
	}
	return d.diff
}

// IsEmpty reports whether there is a non-zero amount of diffs.
func (d *Diff) IsEmpty() bool {
	diffs := d.Diff()
	left := diffs["left"]
	right := diffs["right"]
	return len(left) == 0 && len(right) == 0
}

// EqualTo reports whether d and other are considered equal.
func (d *Diff) EqualTo(other *Diff) bool {
	if d == nil && other == nil {
		return true
	}
	if other == nil {
		return false
	}
	d.Diff()
	other.Diff()
	equal := func(left, right []string) bool {
		leftSet := map[string]struct{}{}
		rightSet := map[string]struct{}{}
		for _, path := range left {
			leftSet[path] = struct{}{}
		}
		for _, path := range right {
			rightSet[path] = struct{}{}
		}
		for path := range leftSet {
			if _, ok := rightSet[path]; !ok {
				return false
			}
		}
		for path := range rightSet {
			if _, ok := leftSet[path]; !ok {
				return false
			}
		}
		return true
	}
	for side := range d.diff {
		if !equal(d.diff[side], other.diff[side]) {
			return false
		}
	}
	return true
}

func (d *Diff) String() string {
	var buf = bytes.NewBuffer(nil)
	diffs := d.Diff()
	for side := range diffs {
		buf.WriteString(fmt.Sprintf("%s: [\n", side))
		list := diffs[side]
		for ii := 0; ii < len(list); ii++ {
			path := list[ii]
			buf.WriteString(fmt.Sprintf("\t%s", path))
			if ii != len(list)-1 {
				buf.WriteString(",\n")
			}
		}
		buf.WriteString("\n]\n")
	}
	return buf.String()
}
