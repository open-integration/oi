package utils

import (
	v1 "github.com/open-integration/core/pkg/api/v1"
	"google.golang.org/grpc"
)

type (
	Proto struct{}
)

// New creates new service for given connection
func (_ Proto) New(cc *grpc.ClientConn) v1.ServiceClient {
	return v1.NewServiceClient(cc)
}
