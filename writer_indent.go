package shelldiff

import (
	"io"
	"strings"
)

// newWriterWithIndent returns a StringWriter that adds a space of indentN size before writing
func newWriterWithIndent(w io.StringWriter, indentN int) *writerWithIndent {
	indent := strings.Repeat(" ", indentN)
	return &writerWithIndent{
		wrappedWriter: w,
		indent:        indent,
	}
}

type writerWithIndent struct {
	wrappedWriter io.StringWriter
	indent        string
}

func (w *writerWithIndent) WriteString(s string) (n int, err error) {
	must(w.wrappedWriter.WriteString(w.indent))
	must(w.wrappedWriter.WriteString(s))
	return len(s) + len(w.indent), nil
}
