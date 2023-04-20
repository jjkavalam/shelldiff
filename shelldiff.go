package shelldiff

import (
	"fmt"
	"io"
)

type Script []*ScriptSection

// ScriptSection is a whitespace trimmed block of script that contains no comment lines
// sections cannot be empty
// each section could have a Name, which is the text of the comment line just preceding it
// to diff two shell scripts, the sections are matched by Name while respecting order
// e.g. [a b c d] vs [x a c e]
// should give a diff [+x -b -d +e]
type ScriptSection struct {
	Name     string
	Contents string
}

func newScriptSection(name string, contents string) *ScriptSection {
	return &ScriptSection{Name: name, Contents: contents}
}

func (ss *ScriptSection) String() string {
	return fmt.Sprintf("[%s] %s", ss.Name, ss.Contents)
}

func (ss *ScriptSection) Equals(other *ScriptSection) bool {
	return *ss == *other
}

func (ss *ScriptSection) WriteDiff(other *ScriptSection, w io.StringWriter) {
	must(w.WriteString(fmt.Sprintf("[%s] -/+\n", ss.Name)))
	iw := newWriterWithIndent(w, 4)
	diffLines(ss.Contents, other.Contents, iw)
}

// Diff compares the given strings and returns true if there are no differences between them
func Diff(s1 string, s2 string, sw io.StringWriter) (bool, error) {
	sc1, err := Parse(s1)
	if err != nil {
		return false, err
	}
	sc2, err := Parse(s2)
	if err != nil {
		return false, err
	}
	d := DiffScripts(sc1, sc2, sw)
	return d, nil
}

// DiffScripts compares scripts and returns true if they have no difference between them
func DiffScripts(this Script, that Script, w io.StringWriter) bool {
	// two sections that have the same name identify points in the script that needs to match
	compareFn := func(a, b *ScriptSection) bool {
		return a.Name == b.Name
	}
	return diff(this, that, w, compareFn)
}

func Parse(s string) (Script, error) {
	tokens, err := tokenize(s)
	if err != nil {
		return nil, err
	}
	p := newParser(tokens)
	sections := p.parse()
	return sections, nil
}
