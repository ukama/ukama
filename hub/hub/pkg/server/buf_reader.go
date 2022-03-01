package server

import (
	"io"
)

type bufReader struct {
	r          io.Reader
	buff       []byte
	afterReset bool
}

func NewBufReader(r io.Reader) *bufReader {
	return &bufReader{
		r:    r,
		buff: []byte{},
	}
}

func (s *bufReader) Read(p []byte) (n int, err error) {
	if s.afterReset && len(s.buff) > 0 {
		n = copy(p, s.buff)
		s.buff = s.buff[n:]
		return n, nil
	}

	n, err = s.r.Read(p)
	if !s.afterReset {
		s.buff = append(s.buff, p[:n]...)
	}

	return n, err
}

func (s *bufReader) Reset() {
	if s.afterReset {
		panic("bufReader.Reset() called twice. Only one call of Reset() is allowed")
	}
	s.afterReset = true
}
