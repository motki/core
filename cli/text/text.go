package text

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/shopspring/decimal"
)

var stdOutIsColorTerm = false

// StandardTerminalWidthInChars describes the width of a standard terminal window.
const StandardTerminalWidthInChars = 80

// CenterText returns the given text centered according to the specified width in characters.
func CenterText(text string, width int) string {
	text = WrapText(text, width)
	lines := bytes.Split([]byte(text), []byte("\n"))
	replacements := make([][]byte, len(lines))

	for i, line := range lines {
		var padding int = (width - len(line)) / 2
		if padding < 0 {
			padding = 0
		}
		replacements[i] = append(append(bytes.Repeat([]byte(" "), padding), line...))
	}

	return string(bytes.Join(replacements, []byte("\n")))
}

// WrapText returns the given text wrapped at the specified width in characters.
func WrapText(text string, width int) string {
	b := bytes.Replace([]byte(text), []byte("\r"), []byte(""), -1)
	b = bytes.Replace(b, []byte("-"), []byte(" "), -1)

	paragraphs := bytes.Split(b, []byte("\n\n"))
	replacements := make([][]byte, len(paragraphs))

	for i, par := range paragraphs {
		par = bytes.Replace(par, []byte("\n"), []byte(" "), -1)
		words := bytes.Split(par, []byte(" "))
		var lines [][]byte
		var currLine []byte
		var currLen int
		for _, word := range words {
			ln := len(word) + 1 // A space on the right side
			if currLen+ln > width {
				lines = append(lines, currLine)
				currLine = []byte{}
				currLen = 0
			}
			currLine = append(append(currLine, word...), []byte(" ")...)
			currLen += ln
		}
		// Append the final line of the paragraph.
		lines = append(lines, currLine)
		replacements[i] = bytes.Join(lines, []byte("\n"))
	}
	return string(bytes.Join(replacements, []byte("\n\n")))
}

// PadIntegerLeft formats the integer and pads the left side to the specified width.
func PadIntegerLeft(i int, width int) string {
	s := strconv.Itoa(i)
	var val string
	for i := len(s); i > 0; i-- {
		val = s[i-1:i] + val
		if i > 1 && (len(s)-i-2)%3 == 0 {
			val = "," + val
		}
	}
	return PadTextLeft(val, width)
}

// PadCurrencyLeft formats the decimal value and pads the left side to the specified width.
func PadCurrencyLeft(d decimal.Decimal, width int) string {
	s := strings.TrimPrefix(d.StringFixed(2), "-")
	var val = s[len(s)-3:]
	for i := len(s) - 3; i > 0; i-- {
		val = s[i-1:i] + val
		if i > 1 && (len(s)-i-2)%3 == 0 {
			val = "," + val
		}
	}
	if d.Sign() < 0 {
		val = "-" + val
	}
	return PadTextLeft(val, width)
}

// PadTextLeft pads the left side of the given text to ensure it is the specified width.
func PadTextLeft(text string, width int) string {
	if width == 0 {
		return ""
	}
	ln := len(text)
	if ln >= width {
		return " " + text[0:width-2] + " "
	}
	return strings.Repeat(" ", width-ln-1) + text + " "
}

// PadTextRight pads the right side of the given text to ensure it is the specified width.
func PadTextRight(text string, width int) string {
	if width == 0 {
		return ""
	}
	ln := len(text)
	if ln >= width {
		return text[0:width-1] + " "
	}
	return text + strings.Repeat(" ", width-ln)
}

func Boldf(format string, args ...interface{}) string {
	if !stdOutIsColorTerm {
		return fmt.Sprintf(format, args...)
	}
	return color.New(color.FgHiWhite).Sprintf(format, args...)
}
