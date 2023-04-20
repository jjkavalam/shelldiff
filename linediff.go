package shelldiff

import (
	"io"
	"strings"
)

type line struct {
	content string
}

func newLine(content string) *line {
	return &line{content: content}
}

func (l *line) String() string {
	return l.content
}

func (l *line) Equals(other *line) bool {
	return l.content == other.content
}

func (l *line) WriteDiff(other *line, w io.StringWriter) {
	// this function will never be called because, two "common" lines
	// will always be also "equal" (see that the conditions of compareFn and Equals are identical);
	// hence WriteDiff will never be called
	panic("never")
}

type lines []*line

func diffLines(this string, that string, w io.StringWriter) {
	thisLines := getLines(this)
	thatLines := getLines(that)
	diff(thisLines, thatLines, w, func(a, b *line) bool {
		return a.content == b.content
	})
}

func getLines(s string) lines {
	ls := strings.Split(s, "\n")
	var lns []*line
	for _, l := range ls {
		lns = append(lns, newLine(l))
	}
	return lns
}
