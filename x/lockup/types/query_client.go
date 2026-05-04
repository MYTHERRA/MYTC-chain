package types

import (
	context "context"

	"google.golang.org/grpc"
)

// QueryClient is the client API for Query service.
type QueryClient interface {
	Lock(ctx context.Context, in *QueryLockRequest, opts ...grpc.CallOption) (*QueryLockResponse, error)
	LocksByOwner(ctx context.Context, in *QueryLocksByOwnerRequest, opts ...grpc.CallOption) (*QueryLocksByOwnerResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Lock(ctx context.Context, in *QueryLockRequest, opts ...grpc.CallOption) (*QueryLockResponse, error) {
	out := new(QueryLockResponse)
	err := c.cc.Invoke(ctx, "/mytc.lockup.v1.Query/Lock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) LocksByOwner(ctx context.Context, in *QueryLocksByOwnerRequest, opts ...grpc.CallOption) (*QueryLocksByOwnerResponse, error) {
	out := new(QueryLocksByOwnerResponse)
	err := c.cc.Invoke(ctx, "/mytc.lockup.v1.Query/LocksByOwner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
