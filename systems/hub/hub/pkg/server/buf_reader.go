/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import "io"

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
