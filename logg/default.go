package logg

import (
	"fmt"
	"os"
)

type fallback struct{}

func (fallback) Infof(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "[fallback logger info] "+format+"\n", args...)
}

func (fallback) Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "[fallback logger error] "+format+"\n", args...)
}
