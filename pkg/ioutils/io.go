package ioutils

import "io"

type nopWriteCloser struct {
	io.Writer
}

func (*nopWriteCloser) Close() error { return nil }

// NopWriteCloser returns a WriteCloser with a no-op Close method wrapping the provided Writer w.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{w}
}
