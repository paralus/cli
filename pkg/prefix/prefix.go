package prefix

import (
	"fmt"
	"io"
	"os"
)

// Each level has 2 spaces for PrefixWriter
const (
	LEVEL_0 = iota
	LEVEL_1
	LEVEL_2
	LEVEL_3
)

// PrefixWriter can write text at various indentation levels.
type PrefixWriter interface {
	// Write writes text with the specified indentation level.
	Write(level int, format string, a ...interface{})
	// WriteLine writes an entire line with no indentation level.
	WriteLine(a ...interface{})
	// Flush forces indentation to be reset.
	Flush()
}

// prefixWriter implements PrefixWriter
type prefixWriter struct {
	out io.Writer
}

var _ PrefixWriter = &prefixWriter{}

// NewPrefixWriter creates a new PrefixWriter.
func NewPrefixWriter() PrefixWriter {
	return &prefixWriter{out: os.Stdout}
}

func (pw *prefixWriter) LevelSpace() string {
	return "  "
}

func (pw *prefixWriter) LevelToPrefix(level int) string {
	levelSpace := pw.LevelSpace()
	prefix := ""
	for i := 0; i < level; i++ {
		prefix += levelSpace
	}
	return prefix
}

func (pw *prefixWriter) Write(level int, format string, a ...interface{}) {
	prefix := pw.LevelToPrefix(level)
	fmt.Fprintf(pw.out, prefix+format, a...)
}

func (pw *prefixWriter) WriteLine(a ...interface{}) {
	fmt.Fprintln(pw.out, a...)
}

func (pw *prefixWriter) Flush() {
	// os.Stdout.Flush()
}
