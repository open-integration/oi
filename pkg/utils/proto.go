package utils

import (
	v1 "github.com/open-integration/oi/pkg/api/v1"
	"google.golang.org/grpc"
)

type (
	// Proto expose abilities to create new service client
	Proto struct{}
)

// New creates new service for given connection
func (_p Proto) New(cc *grpc.ClientConn) v1.ServiceClient {
	return v1.NewServiceClient(cc)
}
