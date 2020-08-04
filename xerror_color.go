package xerror

import "fmt"

// color represents a text color.
type color uint8

const (
	colorBlack color = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
	colorBold     color = 1
	colorDarkGray color = 90
)

// P adds the coloring to the given string.
func (c color) P(s string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), fmt.Sprintf(s, args...))
}
