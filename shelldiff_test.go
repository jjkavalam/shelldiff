package shelldiff_test

import (
	"bytes"
	"github.com/jjkavalam/shelldiff"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.ReadFile("testdata/test1.sh")
	if err != nil {
		t.Fatal(err)
	}
	script, err := shelldiff.Parse(string(f))
	if err != nil {
		t.Fatal(err)
	}
	expected := [][]string{
		{"Section 1", "a\nb\n\nc"},
		{"Section 2", "d"},
		{"Section 3", "e\nf"},
	}
	if len(expected) != len(script) {
		t.Fatalf("expected %d sections; got %d sections", len(expected), len(script))
	}
	for i := 0; i < len(expected); i++ {
		if expected[i][0] != script[i].Name {
			t.Errorf("section %d; expected Name '%s', got '%s'", i, expected[i][0], script[i].Name)
		}
		if expected[i][1] != script[i].Contents {
			t.Errorf("section %d; expected Contents '%s', got '%s'", i, expected[i][1], script[i].Contents)
		}
	}
}

func TestSample(t *testing.T) {
	err := shelldiff.Diff(`# get config
C=...

# compute
R=C+D

# print result
echo $R`, `# compute
R=C-D`, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDiff(t *testing.T) {
	t.Setenv("NO_COLOR", "true")
	f1, err := os.ReadFile("testdata/test1.sh")
	if err != nil {
		t.Fatal(err)
	}
	f2, err := os.ReadFile("testdata/test2.sh")
	if err != nil {
		t.Fatal(err)
	}

	var outBuf bytes.Buffer

	err = shelldiff.Diff(string(f1), string(f2), &outBuf)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(outBuf.String())

	expected := `+[Section 0] x
-[Section 1] a
b

c
+[Section 1] d
-[Section 2] d
[Section 3] e...
`

	if expected != outBuf.String() {
		t.Fatalf("want '%s', got '%s'", expected, outBuf.String())
	}
}
