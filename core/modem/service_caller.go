package modem

import "context"

type (
	// ServiceCaller describes the interface to call the external service
	ServiceCaller interface {
		Call(ctx context.Context, service string, endpoint string, arguments map[string]interface{}, fd string) ([]byte, error)
	}
)
