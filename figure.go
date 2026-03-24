// Package figure renders ASCII art text using FIGlet fonts.
// It supports colored output, custom fonts via io.Reader, and terminal animations.
package figure

import (
	"fmt"
	"io"
	"strings"
)

const ascii_offset = 32
const first_ascii = ' '
const last_ascii = '~'

// Figure holds a phrase rendered in a FIGlet font, with optional color.
type Figure struct {
	phrase string
	font
	strict bool
	color  string
}

// NewFigure creates a Figure rendering phrase in the named font.
// Pass an empty fontName to use the default "standard" font.
// If strict is true, non-ASCII characters return an error; otherwise they are
// replaced with '?'.
// Returns an error if the font cannot be found.
func NewFigure(phrase, fontName string, strict bool) (Figure, error) {
	f, err := newFont(fontName)
	if err != nil {
		return Figure{}, err
	}
	if f.reverse {
		phrase = reverse(phrase)
	}
	return Figure{phrase: phrase, font: f, strict: strict}, nil
}

// NewColorFigure creates a colored Figure. color must be one of:
// red, green, yellow, blue, purple, cyan, gray, white.
// Returns an error if the color or font is invalid.
func NewColorFigure(phrase, fontName string, color string, strict bool) (Figure, error) {
	color = strings.ToLower(color)
	if _, found := colors[color]; !found {
		return Figure{}, fmt.Errorf("invalid color %q: must be one of red, green, yellow, blue, purple, cyan, gray, white", color)
	}
	fig, err := NewFigure(phrase, fontName, strict)
	if err != nil {
		return Figure{}, err
	}
	fig.color = color
	return fig, nil
}

// NewFigureWithFont creates a Figure by reading a FIGlet font from reader.
// Returns an error if the font data cannot be parsed.
func NewFigureWithFont(phrase string, reader io.Reader, strict bool) (Figure, error) {
	f, err := newFontFromReader(reader)
	if err != nil {
		return Figure{}, err
	}
	if f.reverse {
		phrase = reverse(phrase)
	}
	return Figure{phrase: phrase, font: f, strict: strict}, nil
}

// Slicify returns the ASCII art as a slice of strings, one per row.
// If strict mode is enabled and phrase contains non-ASCII characters, an error
// is returned. In non-strict mode, non-ASCII characters are replaced with '?'.
func (fig Figure) Slicify() ([]string, error) {
	var rows []string
	for r := 0; r < fig.font.height; r++ {
		var printRow strings.Builder
		for _, char := range fig.phrase {
			if char < first_ascii || char > last_ascii {
				if fig.strict {
					return nil, fmt.Errorf("invalid character %q: outside printable ASCII range", char)
				}
				char = '?'
			}
			fontIndex := char - ascii_offset
			if int(fontIndex) >= len(fig.font.letters) {
				continue
			}
			letter := fig.font.letters[fontIndex]
			if r >= len(letter) {
				continue
			}
			charRowText := scrub(letter[r], fig.font.hardblank)
			printRow.WriteString(charRowText)
		}
		if r < fig.font.baseline || len(strings.TrimSpace(printRow.String())) > 0 {
			rows = append(rows, strings.TrimRight(printRow.String(), " "))
		}
	}
	return rows, nil
}

func scrub(text string, char byte) string {
	return strings.ReplaceAll(text, string(char), " ")
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
