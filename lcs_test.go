package shelldiff

import (
	"testing"
)

func TestLongestCommonSubsequence(t *testing.T) {
	a := []*ScriptSection{
		newScriptSection("a", ""),
		newScriptSection("b", ""),
		newScriptSection("c", ""),
		newScriptSection("d", ""),
	}

	b := []*ScriptSection{
		newScriptSection("a", ""),
		newScriptSection("c", ""),
		newScriptSection("e", ""),
	}
	r := longestCommonSubsequence[*ScriptSection](a, b, func(a, b *ScriptSection) bool {
		return a.Name == b.Name
	})
	if len(r) != 2 {
		t.Fatal("wrong length")
	}
	if r[0].Name != "a" || r[1].Name != "c" {
		t.Fatal("wrong subsequence")
	}
}
