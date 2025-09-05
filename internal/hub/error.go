package hub

import "errors"

// ErrClientSendBufferFull Custom errors
var (
	ErrClientSendBufferFull = errors.New("chat_client send buffer is full")
)
