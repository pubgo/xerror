package xerror_color

import "fmt"

// Color represents a text Color.
type Color uint8

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Bold     Color = 1
	DarkGray Color = 90
)

// P adds the coloring to the given string.
func (c Color) P(s string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), fmt.Sprintf(s, args...))
}
