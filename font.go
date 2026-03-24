package figure

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"strings"
)

//go:embed fonts
var fontsFS embed.FS

const defaultFont = "standard"

var colors = map[string]string{
	"reset":  "\033[0m",
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"purple": "\033[35m",
	"cyan":   "\033[36m",
	"gray":   "\033[37m",
	"white":  "\033[97m",
}

type font struct {
	name      string
	height    int
	baseline  int
	hardblank byte
	reverse   bool
	letters   [][]string
}

func newFont(name string) (font, error) {
	var f font
	f.setName(name)
	fontBytes, err := fontsFS.ReadFile("fonts/" + f.name + ".flf")
	if err != nil {
		return font{}, fmt.Errorf("font %q not found", f.name)
	}
	fontBytesReader := bytes.NewReader(fontBytes)
	scanner := bufio.NewScanner(fontBytesReader)
	if err := f.setAttributes(scanner); err != nil {
		return font{}, fmt.Errorf("font %q: %w", f.name, err)
	}
	f.setLetters(scanner)
	return f, nil
}

func newFontFromReader(reader io.Reader) (font, error) {
	var f font
	scanner := bufio.NewScanner(reader)
	if err := f.setAttributes(scanner); err != nil {
		return font{}, err
	}
	f.setLetters(scanner)
	return f, nil
}

func (font *font) setName(name string) {
	font.name = name
	if len(name) < 1 {
		font.name = defaultFont
	}
}

func (font *font) setAttributes(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, signature) {
			fields := strings.Fields(text)
			h, err := getHeight(fields)
			if err != nil {
				return err
			}
			b, err := getBaseline(fields)
			if err != nil {
				return err
			}
			font.height = h
			font.baseline = b
			font.hardblank = getHardblank(fields)
			font.reverse = getReverse(fields)
			return nil
		}
	}
	return fmt.Errorf("FIGlet header not found: missing %q signature", signature)
}

func (font *font) setLetters(scanner *bufio.Scanner) {
	font.letters = append(font.letters, make([]string, font.height))
	for i := range font.letters[0] {
		font.letters[0][i] = "  "
	}
	letterIndex := 0
	for scanner.Scan() {
		text, cutLength, letterIndexInc := scanner.Text(), 1, 0
		if lastCharLine(text, font.height) {
			font.letters = append(font.letters, []string{})
			letterIndexInc = 1
			if font.height > 1 {
				cutLength = 2
			}
		}
		if letterIndex > 0 {
			appendText := ""
			if len(text) > 1 {
				appendText = text[:len(text)-cutLength]
			}
			font.letters[letterIndex] = append(font.letters[letterIndex], appendText)
		}
		letterIndex += letterIndexInc
	}
}

func (font *font) evenLetters() {
	var longest int
	for _, letter := range font.letters {
		if len(letter) > 0 {
			longest = max(longest, len(letter[0]))
		}
	}
	for _, letter := range font.letters {
		for i, row := range letter {
			letter[i] = row + strings.Repeat(" ", longest-len(row))
		}
	}
}
