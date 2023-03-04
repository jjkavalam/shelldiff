package shelldiff

import (
	"fmt"
	"os"
	"strings"
)

const escape = "\x1b"

// Color chart: https://gist.github.com/JBlond/2fea43a3049b38287e5e9cefc87b2124
const cBold = "1;37"
const cRed = "0;31"
const cGreen = "0;32"

func Bold(str string) string {
	return colorise(str, cBold)
}

func Red(str string) string {
	return colorise(str, cRed)
}

func Green(str string) string {
	return colorise(str, cGreen)
}

// https://no-color.org/
func isNoColor() bool {
	return os.Getenv("NO_COLOR") != ""
}

func colorise(str string, code string) string {
	if isNoColor() {
		return str
	}
	lines := strings.Split(str, "\n")
	outLines := make([]string, 0)
	for _, line := range lines {
		coloredLine := fmt.Sprintf("%s[%sm%s%s[%sm", escape, code, line, escape, "0")
		outLines = append(outLines, coloredLine)
	}
	return strings.Join(outLines, "\n")
}
