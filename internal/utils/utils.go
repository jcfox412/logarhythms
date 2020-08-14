package utils

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	clearLine          = "\033[K"
	cursorAbsoluteLeft = "\033[G"
	setBold            = "\033[1m"
	setUnbold          = "\033[0m"
)

// BeatCount returns a string representation of whole-beat increments at the top
// of the track printout.
func BeatCount(beatCount, divisionsPerBeat int) (string, error) {
	if divisionsPerBeat <= 0 {
		return "", errors.New("divisionsPerBeat must be greater than 0")
	}

	if beatCount%divisionsPerBeat != 0 {
		return fmt.Sprint("  "), nil
	}

	beatCount = beatCount/divisionsPerBeat + 1

	beatStr := "%d"
	if beatCount < 10 {
		beatStr = "%d "
	}

	return fmt.Sprintf(beatStr, beatCount), nil
}

// Bold returns an ANSI-supported bolding of the input text.
func Bold(text string) string {
	return fmt.Sprintf("%s%s%s", setBold, text, setUnbold)
}

// BeatTracker returns a string representation of an asterisk used to track the
// current position of the beat in a track.
func BeatTracker() string {
	out := ""
	out += fmt.Sprint(cursorLeft(3))
	out += fmt.Sprint("   *")

	return out
}

// CursorToNextColumn returns an ANSI-enabled string for moving the track player's
// cursor to the top of its next column.
func CursorToNextColumn(columnHeight int) string {
	out := ""
	out += fmt.Sprint(cursorUp(columnHeight))
	out += fmt.Sprint(cursorRight(1))

	return out
}

// CursorToNextRow returns an ANSI-enabled string for moving the track player's
// cursor to the row below the current row.
func CursorToNextRow() string {
	out := ""
	out += fmt.Sprint(cursorDown(1))
	out += fmt.Sprint(cursorLeft(2))

	return out
}

// ClearLine returns an ANSI-enabled string for moving the track player's
// cursor to the very left, clearing the beat line, and moving the header's
// width to the right.
func ClearLine(headerWidth int) string {
	out := ""
	out += fmt.Sprint(cursorAbsoluteLeft)
	out += fmt.Sprint(clearLine)
	out += fmt.Sprint(cursorRight(headerWidth))

	return out
}

func cursorUp(spaces int) string {
	return moveCursor(spaces, "A")
}

func cursorDown(spaces int) string {
	return moveCursor(spaces, "B")
}

func cursorRight(spaces int) string {
	return moveCursor(spaces, "C")
}

func cursorLeft(spaces int) string {
	return moveCursor(spaces, "D")
}

func moveCursor(spaces int, direction string) string {
	return fmt.Sprintf("\033[%d%s", spaces, direction)
}
