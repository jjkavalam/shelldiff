package shelldiff

import (
	"fmt"
	"io"
	"strings"
)

// DiffScripts compares scripts and returns true if they have no difference between them
func DiffScripts(this Script, that Script, w io.StringWriter) bool {
	// two sections that have the same name identify points in the script that needs to match
	compareFn := func(a, b *ScriptSection) bool {
		return a.Name == b.Name
	}

	// find the longest common subsequence of matching sections
	common := longestCommonSubsequence[*ScriptSection](this, that, compareFn)

	// keep a cursor for each side; advance each to the next match in the common sequence;
	// anything skipped on this was deleted; while anything skipped on that was added
	i := 0
	j := 0

	foundSomeDifference := false

	for k := range common {

		for ; i < len(this); i++ {
			if compareFn(this[i], common[k]) {
				break
			}
			must(w.WriteString("-" + red(shorten(this[i].String())) + "\n"))
			foundSomeDifference = true
		}

		for ; j < len(that); j++ {
			if compareFn(that[j], common[k]) {
				break
			}
			must(w.WriteString("+" + green(shorten(that[j].String())) + "\n"))
			foundSomeDifference = true
		}

		if i < len(this) && j < len(that) {
			// compare
			if *this[i] != *that[j] {
				must(w.WriteString(fmt.Sprintf("-%s\n+%s", red(this[i].String()), green(that[j].String())) + "\n"))
				foundSomeDifference = true
			} else {
				must(w.WriteString(shortenFirstLine(this[i].String()) + "\n"))
			}
		}

		i++
		j++

	}

	for ; i < len(this); i++ {
		must(w.WriteString("-" + red(shorten(this[i].String())) + "\n"))
		foundSomeDifference = true
	}

	for ; j < len(that); j++ {
		must(w.WriteString("+" + green(shorten(that[j].String())) + "\n"))
		foundSomeDifference = true
	}

	if !foundSomeDifference {
		must(w.WriteString("There are no differences !\n"))
	}

	return !foundSomeDifference
}

// shortenFirstLine only ever prints the first line and also trims the string to a certain length
func shortenFirstLine(s string) string {
	nlPos := strings.Index(s, "\n")
	if nlPos == -1 {
		return shortenHelper(s, true)
	}
	return shortenHelper(s[:nlPos], false) + "..."
}

// shorten trims the string to a certain length
func shorten(s string) string {
	return shortenHelper(s, true)
}

func shortenHelper(s string, ellipsis bool) string {
	if len(s) < OptionShortenValueDiffs {
		return s
	}
	s = s[:OptionShortenValueDiffs]
	if ellipsis {
		return s + "..."
	} else {
		return s
	}
}
