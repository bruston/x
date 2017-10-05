package squish

import (
	"bufio"
	"bytes"
	"io"
)

func Encode(dst io.Writer, src io.Reader) error {
	br := bufio.NewReader(src)
	bw := bufio.NewWriter(dst)
	previous, err := br.ReadByte()
	if err != nil {
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		}
		return err
	}
	buf := make([]byte, 2)
	var count byte = 1
	for {
		current, err := br.ReadByte()
		if err != nil {
			if err == io.EOF {
				write(bw, count, previous, buf)
				break
			}
			return err
		}
		if current == previous {
			if count == 255 {
				write(bw, count, current, buf)
				count = 1
				continue
			}
			count++
			continue
		}
		write(bw, count, previous, buf)
		previous = current
		count = 1
	}
	return bw.Flush()
}

func write(w *bufio.Writer, count, ch byte, b []byte) {
	b[0], b[1] = count, ch
	w.Write(b)
}

func Decode(src io.Reader, dst io.Writer) error {
	br := bufio.NewReader(src)
	bw := bufio.NewWriter(dst)
	buf := make([]byte, 2)
	for {
		if n, err := br.Read(buf); err != nil {
			if err == io.EOF {
				if n > 0 {
					return io.ErrUnexpectedEOF
				}
				break
			}
			return err
		}
		bw.Write(bytes.Repeat(buf[1:], int(buf[0])))
	}
	return bw.Flush()
}
