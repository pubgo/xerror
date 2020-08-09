package xerror_color

import "fmt"

// Color represents a text Color.
type Color uint8

const (
	ColorBlack Color = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBold     Color = 1
	ColorDarkGray Color = 90
)

// P adds the coloring to the given string.
func (c Color) P(s string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), fmt.Sprintf(s, args...))
}
