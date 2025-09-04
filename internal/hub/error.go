package hub

import "errors"

// ErrClientSendBufferFull Custom errors
var (
	ErrClientSendBufferFull = errors.New("client send buffer is full")
)
