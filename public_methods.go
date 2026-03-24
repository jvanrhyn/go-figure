package figure

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

// Print renders the figure to stdout.
func (fig Figure) Print() {
	rows, err := fig.Slicify()
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows {
		if fig.color != "" {
			row = colors[fig.color] + row + colors["reset"]
		}
		fmt.Println(row)
	}
}

// ColorString returns the figure as a color-escaped string suitable for terminal output.
func (fig Figure) ColorString() string {
	rows, err := fig.Slicify()
	if err != nil {
		return ""
	}
	var s strings.Builder
	for _, row := range rows {
		if fig.color != "" {
			row = colors[fig.color] + row + colors["reset"]
		}
		fmt.Fprintf(&s, "%s\n", row)
	}
	return s.String()
}

// String returns the figure as a plain string (implements fmt.Stringer).
func (fig Figure) String() string {
	rows, err := fig.Slicify()
	if err != nil {
		return ""
	}
	var s strings.Builder
	for _, row := range rows {
		fmt.Fprintf(&s, "%s\n", row)
	}
	return s.String()
}

// Scroll animates the phrase scrolling across the terminal for duration milliseconds.
// stillness controls the pause between frames in milliseconds.
// direction: "left" (default) or "right".
func (fig Figure) Scroll(duration, stillness int, direction string) {
	endTime := time.Now().Add(time.Duration(duration) * time.Millisecond)
	fig.phrase = fig.phrase + "   "
	clearScreen()
	for time.Now().Before(endTime) {
		chars := []byte(fig.phrase)
		if strings.HasPrefix(strings.ToLower(direction), "r") {
			fig.phrase = string(append(chars[len(chars)-1:], chars[:len(chars)-1]...))
		} else {
			fig.phrase = string(append(chars[1:], chars[0]))
		}
		fig.Print()
		sleep(stillness)
		clearScreen()
	}
}

// Blink animates the figure blinking for duration milliseconds.
// timeOn and timeOff control the on/off durations in milliseconds.
// If timeOff is negative, it defaults to timeOn.
func (fig Figure) Blink(duration, timeOn, timeOff int) {
	if timeOff < 0 {
		timeOff = timeOn
	}
	endTime := time.Now().Add(time.Duration(duration) * time.Millisecond)
	clearScreen()
	for time.Now().Before(endTime) {
		fig.Print()
		sleep(timeOn)
		clearScreen()
		sleep(timeOff)
	}
}

// Dance animates the figure with an alternating character dance effect for duration milliseconds.
// freeze controls the pause between frames in milliseconds.
func (fig Figure) Dance(duration, freeze int) {
	endTime := time.Now().Add(time.Duration(duration) * time.Millisecond)
	f := fig.font
	f.evenLetters()
	figures := []Figure{{font: f}, {font: f}}
	clearScreen()
	for i, c := range fig.phrase {
		appenders := []string{" ", " "}
		appenders[i%2] = string(c)
		for f := range figures {
			figures[f].phrase = figures[f].phrase + appenders[f]
		}
	}
	for p := 0; time.Now().Before(endTime); p ^= 1 {
		figures[p].Print()
		figures[1-p].Print()
		sleep(freeze)
		clearScreen()
	}
}

// Write renders the figure to w, one row per line.
func Write(w io.Writer, fig Figure) {
	rows, err := fig.Slicify()
	if err != nil {
		return
	}
	for _, row := range rows {
		fmt.Fprintf(w, "%v\n", row)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func sleep(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}
