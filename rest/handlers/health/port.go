package health

import "context"

// Pinger defines the capability to check connectivity with a downstream service.
type Pinger interface {
	PingContext(ctx context.Context) error
}
