package ioutils

import (
	"io"
	"os"
)

type mulWriteCloser struct {
	writeClosers []io.WriteCloser
}

func (t *mulWriteCloser) Write(p []byte) (n int, err error) {
	for _, w := range t.writeClosers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

func (t *mulWriteCloser) Close() (err error) {
	for _, c := range t.writeClosers {
		if c == os.Stdin || c == os.Stdout || c == os.Stderr {
			continue
		}

		err = c.Close()

		if err != nil {
			return
		}
	}

	return nil
}

// MultiWriteCloser creates a writeCloser that duplicates its writes or closes to all the provided writers
func MultiWriteCloser(writeClosers ...io.WriteCloser) io.WriteCloser {
	allWriters := make([]io.WriteCloser, 0, len(writeClosers))
	for _, wc := range writeClosers {
		if mw, ok := wc.(*mulWriteCloser); ok {
			allWriters = append(allWriters, mw.writeClosers...)
		} else {
			allWriters = append(allWriters, wc)
		}
	}
	return &mulWriteCloser{allWriters}
}
