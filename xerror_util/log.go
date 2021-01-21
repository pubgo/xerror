package xerror_util

import (
	slog "log"
	"os"
)

var log = slog.New(os.Stderr, "xerror_util", slog.LstdFlags|slog.Llongfile)
