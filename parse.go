package shelldiff

import (
	"bufio"
	"strings"
)

type token interface {
}

type comment struct {
	Label string
}

type regular struct {
	Line string
}

type empty struct {
}

// tokenize: each line is a token (it is one of COMMENT, EMPTY, REGULAR)
func tokenize(s string) ([]token, error) {
	tokens := make([]token, 0)
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") {
			tokens = append(tokens, &comment{line[1:]})
		} else if strings.TrimSpace(line) == "" {
			tokens = append(tokens, &empty{})
		} else {
			tokens = append(tokens, &regular{line})
		}
	}
	return tokens, sc.Err()
}

type parser struct {
	next   int
	tokens []token
}

func newParser(tokens []token) *parser {
	return &parser{next: 0, tokens: tokens}
}

func (p *parser) peek() token {
	if p.next < len(p.tokens) {
		return p.tokens[p.next]
	} else {
		return nil
	}
}

func (p *parser) pop() token {
	if p.next < len(p.tokens) {
		item := p.tokens[p.next]
		p.next++
		return item
	} else {
		return nil
	}
}

func (p *parser) parse() []*ScriptSection {
	sections := make([]*ScriptSection, 0)

	for p.peek() != nil {
		section := p.parseSection()
		// skip empty sections
		if strings.TrimSpace(section.Contents) != "" {
			sections = append(sections, section)
		}
	}

	return sections
}

// find first REGULAR line; start a new section and Name with last seen "comment heading"
// accumulate subsequent lines under the section until another COMMENT or EOF is reached
//
// a "comment heading" is the first non-empty line in a block of comments
func (p *parser) parseSection() *ScriptSection {
	var commentBlock []string

outer1:
	for p.peek() != nil {
		t := p.peek()
		switch tt := t.(type) {
		case *comment:
			text := strings.TrimSpace(strings.TrimLeft(tt.Label, "#"))
			if text != "" {
				commentBlock = append(commentBlock, text)
			}
		case *regular, *empty:
			break outer1
		}
		p.pop()
	}

	var label string
	if len(commentBlock) > 0 {
		label = commentBlock[0]
	}

	var sb strings.Builder

outer2:
	for p.peek() != nil {
		t := p.peek()
		switch tt := t.(type) {
		case *comment:
			break outer2
		case *empty:
			sb.WriteString("\n")
		case *regular:
			sb.WriteString(tt.Line)
			sb.WriteString("\n")
		}
		p.pop()
	}

	return newScriptSection(label, trimEmptyLines(sb.String()))

}

func trimEmptyLines(s string) string {
	lines := strings.Split(s, "\n")
	for len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}
