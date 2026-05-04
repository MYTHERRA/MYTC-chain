package types

import (
	context "context"

	"google.golang.org/grpc"
)

// MsgClient is the client API for Msg service.
type MsgClient interface {
	LockTokens(ctx context.Context, in *MsgLockTokens, opts ...grpc.CallOption) (*MsgLockTokensResponse, error)
	Unlock(ctx context.Context, in *MsgUnlock, opts ...grpc.CallOption) (*MsgUnlockResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) LockTokens(ctx context.Context, in *MsgLockTokens, opts ...grpc.CallOption) (*MsgLockTokensResponse, error) {
	out := new(MsgLockTokensResponse)
	err := c.cc.Invoke(ctx, "/mytc.lockup.v1.Msg/LockTokens", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Unlock(ctx context.Context, in *MsgUnlock, opts ...grpc.CallOption) (*MsgUnlockResponse, error) {
	out := new(MsgUnlockResponse)
	err := c.cc.Invoke(ctx, "/mytc.lockup.v1.Msg/Unlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
