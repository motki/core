// +build !windows

package text

import (
	"os"

	"github.com/mattn/go-isatty"
)

func init() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		stdOutIsColorTerm = true
	}
}
