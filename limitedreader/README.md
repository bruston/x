Exactly the same as [io.LimitedReader](http://golang.org/pkg/io/#LimitedReader) except it returns an 
ErrReadLimitExceeded error if there's too much stuff to read instead of io.EOF, so you can distinguish between 
actual EOF.

ErrReadLimitExceeded is defined as:
```Go
var ErrReadLimitExceeded = errors.New("read limit exceeded")
```
