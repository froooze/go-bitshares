package transport

import (
	"encoding/json"
	"errors"
)

var (
	ErrNotConnected = errors.New("transport is not connected")
	ErrUnsupported  = errors.New("operation is not supported by this transport")
	ErrShutdown     = errors.New("transport is shut down")
)

// RPCRequest is a generic JSON-RPC request envelope.
type RPCRequest struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
	ID     uint64 `json:"id"`
}

// RPCError is a transport-level JSON-RPC error.
type RPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// RPCResponse is a generic JSON-RPC response envelope.
type RPCResponse struct {
	ID     uint64           `json:"id"`
	Result *json.RawMessage `json:"result,omitempty"`
	Error  *RPCError        `json:"error,omitempty"`
}
