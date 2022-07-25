package pkg

import (
	"fmt"
	"io"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Stdout() io.Writer
	Stderr() io.Writer
}

type logger struct {
	stderr  io.Writer
	stdout  io.Writer
	verbose bool
}

func NewLogger(stdout io.Writer, stderr io.Writer, verbose bool) Logger {
	return &logger{
		stderr:  stderr,
		stdout:  stdout,
		verbose: verbose,
	}
}

func (l logger) Printf(format string, v ...interface{}) {
	if l.verbose {
		fmt.Fprintf(l.stdout, format, v...)
	}
}

func (l logger) Errorf(format string, v ...interface{}) {
	if l.verbose {
		fmt.Fprintf(l.stderr, format, v...)
	}
}

func (l logger) Stdout() io.Writer {
	return l.stdout
}

func (l logger) Stderr() io.Writer {
	return l.stderr
}
