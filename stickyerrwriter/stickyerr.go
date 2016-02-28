// Package stickyerrwriter provides a Writer based on the
// sticky error writer written about here:
// http://blog.golang.org/error-handling-and-go
// Implements the Writer interface. Stashes the first error
// encountered which can be later retrieved via the Err method.
// Subsequent writes following an error are essentially NOOPs.
package stickyerrwriter

import "io"

type StickyErrorWriter struct {
	io.Writer
	err error
}

func (s *StickyErrorWriter) Write(p []byte) (int, error) {
	if s.err != nil {
		return len(p), nil
	}
	n, err := s.Writer.Write(p)
	if err != nil {
		s.err = err
		return len(p), nil
	}
	return n, nil
}

func (s *StickyErrorWriter) Err() error { return s.err }
