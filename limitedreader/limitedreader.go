package limitedreader

import (
	"errors"
	"io"
)

var ErrReadLimitExceeded = errors.New("read limit exceeded")

type LimitedReader struct {
	R io.Reader
	N int64
}

func (l *LimitedReader) Read(b []byte) (int, error) {
	if l.N <= 0 {
		if _, err := l.R.Read(make([]byte, 1)); err == nil {
			return 0, ErrReadLimitExceeded
		}
		return 0, io.EOF
	}
	if int64(len(b)) > l.N {
		b = b[0:l.N]
	}
	n, err := l.R.Read(b)
	l.N -= int64(n)
	return n, err

}

func New(r io.Reader, n int64) *LimitedReader {
	return &LimitedReader{r, n}
}
