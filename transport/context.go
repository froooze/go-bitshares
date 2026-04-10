package transport

import (
	"context"
	"errors"
)

// ErrNilContext reports that a transport API was called with a nil context.
var ErrNilContext = errors.New("nil context")

func requireContext(ctx context.Context) error {
	if ctx == nil {
		return ErrNilContext
	}
	return nil
}
