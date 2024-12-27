package internal

import (
	"fmt"
	"io"
)

var DefaultLogger *Logger

type Logger struct {
	verbose bool
	out     io.Writer
}

func NewLogger(out io.Writer, verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
		out:     out,
	}
}

func (l *Logger) Printf(format string, args ...any) {
	if !l.verbose {
		return
	}
	fmt.Fprintf(l.out, format+"\n", args...)
}

func (l *Logger) Notice(format string, args ...any) {
	fmt.Fprintf(l.out, format+"\n", args...)
}
