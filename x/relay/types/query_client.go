package types

import (
	context "context"

	"google.golang.org/grpc"
)

// QueryClient is the client API for Query service.
type QueryClient interface {
	Endpoint(ctx context.Context, in *QueryEndpointRequest, opts ...grpc.CallOption) (*QueryEndpointResponse, error)
	Endpoints(ctx context.Context, in *QueryEndpointsRequest, opts ...grpc.CallOption) (*QueryEndpointsResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Endpoint(ctx context.Context, in *QueryEndpointRequest, opts ...grpc.CallOption) (*QueryEndpointResponse, error) {
	out := new(QueryEndpointResponse)
	err := c.cc.Invoke(ctx, "/mytc.relay.v1.Query/Endpoint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Endpoints(ctx context.Context, in *QueryEndpointsRequest, opts ...grpc.CallOption) (*QueryEndpointsResponse, error) {
	out := new(QueryEndpointsResponse)
	err := c.cc.Invoke(ctx, "/mytc.relay.v1.Query/Endpoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
