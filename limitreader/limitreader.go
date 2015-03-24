package limitreader

import (
	"errors"
	"io"
)

type LimitedReader struct {
	R io.Reader
	N int64
}

func New(r io.Reader, n int64) io.Reader { return &LimitedReader{r, n} }

func (l *LimitedReader) Read(b []byte) (int, error) {
	if l.N <= 0 {
		return 0, errors.New("exceeded read limit")
	}
	if int64(len(b)) > l.N {
		b = b[0:l.N]
	}
	n, err := l.R.Read(b)
	l.N -= int64(n)
	return n, err
}
