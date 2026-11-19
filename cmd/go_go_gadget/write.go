package go_go_gadget

import (
	"fmt"
	"io"
)

func write(w io.Writer, a ...any) {
	_, _ = fmt.Fprint(w, a...)
}

func writeln(w io.Writer, a ...any) {
	_, _ = fmt.Fprintln(w, a...)
}

func writef(w io.Writer, format string, a ...any) {
	_, _ = fmt.Fprintf(w, format, a...)
}
