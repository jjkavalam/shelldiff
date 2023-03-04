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

func Diff(s1 string, s2 string, sw io.StringWriter) error {
	sc1, err := Parse(s1)
	if err != nil {
		return err
	}
	sc2, err := Parse(s2)
	if err != nil {
		return err
	}
	return diff(sc1, sc2, sw)
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