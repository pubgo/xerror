package xerror

import (
	slog "log"
	"os"
)

var log = slog.New(os.Stderr, "xerror", slog.LstdFlags|slog.Llongfile)
