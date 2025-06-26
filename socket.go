package main

import (
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

type timeoutReader struct {
	net.Conn
	once sync.Once
}

func (r *timeoutReader) Read(b []byte) (int, error) {
	n, err := r.Conn.Read(b)

	// Set a read deadline only after the first Read completes
	r.once.Do(func() {
		r.Conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	})

	// If we got a timeout, treat it as an io.EOF so the bufio.Scanner handles
	// the error as if it was the normal end of the stream.
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return n, io.EOF
	}
	return n, err
}
