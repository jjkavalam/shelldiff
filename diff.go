package shelldiff

import (
	"fmt"
	"io"
)

func DiffScripts(this Script, that Script, w io.StringWriter) error {
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

	for k := range common {

		for ; i < len(this); i++ {
			if compareFn(this[i], common[k]) {
				break
			}
			_, err := w.WriteString("-" + red(shorten(this[i].String())) + "\n")
			if err != nil {
				return err
			}
		}

		for ; j < len(that); j++ {
			if compareFn(that[j], common[k]) {
				break
			}
			_, err := w.WriteString("+" + green(shorten(that[j].String())) + "\n")
			if err != nil {
				return err
			}
		}

		if i < len(this) && j < len(that) {
			// compare
			if *this[i] != *that[j] {
				_, err := w.WriteString(fmt.Sprintf("-%s\n+%s", red(this[i].String()), green(that[j].String())) + "\n")
				if err != nil {
					return err
				}
			} else {
				_, err := w.WriteString(shorten(this[i].String()) + "\n")
				if err != nil {
					return err
				}
			}
		}

		i++
		j++

	}

	for ; i < len(this); i++ {
		_, err := w.WriteString("-" + red(shorten(this[i].String())) + "\n")
		if err != nil {
			return err
		}
	}

	for ; j < len(that); j++ {
		_, err := w.WriteString("+" + green(shorten(that[j].String())) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func shorten(s string) string {
	if len(s) < OptionShortenValueDiffs {
		return s
	}
	return s[:OptionShortenValueDiffs] + "..."
}
