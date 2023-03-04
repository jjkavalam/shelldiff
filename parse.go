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

// 	find first REGULAR line; start a new section and Name with last seen COMMENT
//	accumulate subsequent lines under the section until another COMMENT or EOF is reached
func (p *parser) parseSection() *ScriptSection {
	var label string
outer1:
	for p.peek() != nil {
		t := p.peek()
		switch tt := t.(type) {
		case *comment:
			text := strings.TrimSpace(strings.TrimLeft(tt.Label, "#"))
			if text != "" {
				label = text
			}
		case *regular:
			break outer1
		}
		p.pop()
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

	return newScriptSection(label, trimTrailingEmptyLines(sb.String()))

}

func trimTrailingEmptyLines(s string) string {
	// Split the multiline string into lines
	lines := strings.Split(s, "\n")

	// Remove trailing empty lines
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	// Join the lines back into a string
	result := strings.Join(lines, "\n")

	return result
}
