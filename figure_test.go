package figure

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func Example_alphabet() {
	myFigure, _ := NewFigure("Hello world", "alphabet", true)
	fmt.Println(myFigure)
	// Output:
	// H  H     l l                     l    d
	// H  H     l l                     l    d
	// HHHH eee l l ooo   w   w ooo rrr l  ddd
	// H  H e e l l o o   w w w o o r   l d  d
	// H  H ee  l l ooo    w w  ooo r   l  ddd
}

func Example_standard() {
	myFigure, _ := NewFigure("Hello world", "", true)
	fmt.Println(myFigure)
	// Output:
	//   _   _          _   _                                       _       _
	//  | | | |   ___  | | | |   ___     __      __   ___    _ __  | |   __| |
	//  | |_| |  / _ \ | | | |  / _ \    \ \ /\ / /  / _ \  | '__| | |  / _` |
	//  |  _  | |  __/ | | | | | (_) |    \ V  V /  | (_) | | |    | | | (_| |
	//  |_| |_|  \___| |_| |_|  \___/      \_/\_/    \___/  |_|    |_|  \__,_|
}

// --- Constructor error cases ---

func TestNewFigure_InvalidFont(t *testing.T) {
	_, err := NewFigure("Hi", "no-such-font", true)
	if err == nil {
		t.Fatal("expected error for unknown font, got nil")
	}
}

func TestNewFigure_EmptyString(t *testing.T) {
	fig, err := NewFigure("", "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Empty phrase should produce no error; rows may be empty strings up to baseline.
	_, err = fig.Slicify()
	if err != nil {
		t.Fatalf("Slicify error on empty phrase: %v", err)
	}
}

func TestNewColorFigure_InvalidColor(t *testing.T) {
	_, err := NewColorFigure("Hi", "", "ultraviolet", false)
	if err == nil {
		t.Fatal("expected error for unknown color, got nil")
	}
}

func TestNewColorFigure_ValidColor(t *testing.T) {
	_, err := NewColorFigure("Hi", "", "red", false)
	if err != nil {
		t.Fatalf("unexpected error for valid color: %v", err)
	}
}

// --- Slicify edge cases ---

func TestSlicify_NonAsciiNonStrict(t *testing.T) {
	fig, err := NewFigure("A£B", "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rows, err := fig.Slicify()
	if err != nil {
		t.Fatalf("unexpected Slicify error in non-strict mode: %v", err)
	}
	// Should produce output (non-ASCII replaced with '?')
	if len(rows) == 0 {
		t.Error("expected non-empty rows for A?B")
	}
}

func TestSlicify_NonAsciiStrict(t *testing.T) {
	fig, err := NewFigure("A£B", "", true)
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}
	_, err = fig.Slicify()
	if err == nil {
		t.Fatal("expected error for non-ASCII in strict mode, got nil")
	}
}

// --- Write ---

func TestWrite(t *testing.T) {
	fig, err := NewFigure("Hi", "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	Write(&buf, fig)
	if buf.Len() == 0 {
		t.Error("Write produced no output")
	}
	if !strings.Contains(buf.String(), "\n") {
		t.Error("Write output has no newlines")
	}
}

// --- NewFigureWithFont ---

// minimalFLF is a tiny valid FIGlet font (height=3) containing space and '!'.
// Lines end with @ (intermediate) or @@ (last line of character).
const minimalFLF = "flf2a$ 3 2 15 0 0\nComment\n @\n @\n @@\n  @\n  @\n  @@\n"

func TestNewFigureWithFont(t *testing.T) {
	fig, err := NewFigureWithFont("!", strings.NewReader(minimalFLF), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rows, err := fig.Slicify()
	if err != nil {
		t.Fatalf("Slicify error: %v", err)
	}
	_ = rows // font may produce empty rows; just verify no crash
}

// --- Parser unit tests ---

func TestNewFigureWithFont_MalformedHeader(t *testing.T) {
	_, err := NewFigureWithFont("x", strings.NewReader("flf2\n"), false)
	if err == nil {
		t.Fatal("expected error for malformed FIGlet header, got nil")
	}
}

func TestGetHeight(t *testing.T) {
	fields := strings.Fields("flf2a$ 6 5 15 15 8")
	if got, err := getHeight(fields); err != nil || got != 6 {
		t.Errorf("getHeight = %d, %v; want 6, nil", got, err)
	}
}

func TestGetBaseline(t *testing.T) {
	fields := strings.Fields("flf2a$ 6 5 15 15 8")
	if got, err := getBaseline(fields); err != nil || got != 5 {
		t.Errorf("getBaseline = %d, %v; want 5, nil", got, err)
	}
}

func TestGetHardblank(t *testing.T) {
	fields := strings.Fields("flf2a$ 6 5 15 15 8")
	if got := getHardblank(fields); got != '$' {
		t.Errorf("getHardblank = %q, want '$'", got)
	}
}

func TestGetHardblank_Blacklisted(t *testing.T) {
	// 'a' and '2' are blacklisted hardblanks; should return space
	if got := getHardblank(strings.Fields("flf2aa 6 5 15 15 8")); got != ' ' {
		t.Errorf("getHardblank (blacklisted 'a') = %q, want ' '", got)
	}
	if got := getHardblank(strings.Fields("flf2a2 6 5 15 15 8")); got != ' ' {
		t.Errorf("getHardblank (blacklisted '2') = %q, want ' '", got)
	}
}

func TestGetReverse_True(t *testing.T) {
	fields := strings.Fields("flf2a$ 6 5 15 15 8 1")
	if !getReverse(fields) {
		t.Error("getReverse should be true when field 6 is '1'")
	}
}

func TestGetReverse_False(t *testing.T) {
	fields := strings.Fields("flf2a$ 6 5 15 15 8")
	if getReverse(fields) {
		t.Error("getReverse should be false when field 6 is absent")
	}
}

func TestLastCharLine(t *testing.T) {
	tests := []struct {
		text   string
		height int
		want   bool
	}{
		{"text@@", 3, true},
		{"text@", 1, true},
		{"text##", 3, true},
		{"text$$", 3, true},
		{"text", 3, false},
	}
	for _, tt := range tests {
		got := lastCharLine(tt.text, tt.height)
		if got != tt.want {
			t.Errorf("lastCharLine(%q, %d) = %v, want %v", tt.text, tt.height, got, tt.want)
		}
	}
}
