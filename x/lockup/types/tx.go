package types

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	gogogrpc "github.com/gogo/protobuf/grpc"
	gogotypes "github.com/gogo/protobuf/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc"
)

const (
	TypeMsgLockTokens = "lock_tokens"
	TypeMsgUnlock     = "unlock"
)

var _ sdk.Msg = &MsgLockTokens{}
var _ sdk.Msg = &MsgUnlock{}

// MsgLockTokens defines a message to lock tokens
type MsgLockTokens struct {
	Owner    string            `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	Amount   sdk.Coin          `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
	Duration *gogotypes.Duration `protobuf:"bytes,3,opt,name=duration,proto3" json:"duration,omitempty"`
}

func NewMsgLockTokens(owner string, amount sdk.Coin, duration time.Duration) *MsgLockTokens {
	return &MsgLockTokens{
		Owner:    owner,
		Amount:   amount,
		Duration: gogotypes.DurationProto(duration),
	}
}

func (msg *MsgLockTokens) Route() string {
	return RouterKey
}

func (msg *MsgLockTokens) Type() string {
	return TypeMsgLockTokens
}

func (msg *MsgLockTokens) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (msg *MsgLockTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgLockTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return err
	}
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.ErrInvalidCoins
	}
	if msg.Duration == nil {
		return ErrLockNotMatured // Reusing error or define ErrInvalidDuration
	}
	dur, err := gogotypes.DurationFromProto(msg.Duration)
	if err != nil {
		return err
	}
	if dur <= 0 {
		return ErrLockNotMatured // Reusing error or define ErrInvalidDuration
	}
	return nil
}

func (msg *MsgLockTokens) ProtoMessage() {}
func (msg *MsgLockTokens) Reset()        { *msg = MsgLockTokens{} }
func (msg *MsgLockTokens) String() string {
	return "MsgLockTokens"
}

// MsgUnlock defines a message to unlock tokens
type MsgUnlock struct {
	Owner  string `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	LockId uint64 `protobuf:"varint,2,opt,name=lock_id,json=lockId,proto3" json:"lock_id"`
}

func NewMsgUnlock(owner string, lockId uint64) *MsgUnlock {
	return &MsgUnlock{
		Owner:  owner,
		LockId: lockId,
	}
}

func (msg *MsgUnlock) Route() string {
	return RouterKey
}

func (msg *MsgUnlock) Type() string {
	return TypeMsgUnlock
}

func (msg *MsgUnlock) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (msg *MsgUnlock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnlock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return err
	}
	return nil
}

func (msg *MsgUnlock) ProtoMessage() {}
func (msg *MsgUnlock) Reset()        { *msg = MsgUnlock{} }
func (msg *MsgUnlock) String() string {
	return "MsgUnlock"
}

// Response types
type MsgLockTokensResponse struct {
	LockId uint64 `protobuf:"varint,1,opt,name=lock_id,json=lockId,proto3" json:"lock_id"`
}

func (m *MsgLockTokensResponse) Reset()         { *m = MsgLockTokensResponse{} }
func (m *MsgLockTokensResponse) String() string { return "MsgLockTokensResponse" }
func (m *MsgLockTokensResponse) ProtoMessage()  {}

type MsgUnlockResponse struct{}

func (m *MsgUnlockResponse) Reset()         { *m = MsgUnlockResponse{} }
func (m *MsgUnlockResponse) String() string { return "MsgUnlockResponse" }
func (m *MsgUnlockResponse) ProtoMessage()  {}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	LockTokens(context.Context, *MsgLockTokens) (*MsgLockTokensResponse, error)
	Unlock(context.Context, *MsgUnlock) (*MsgUnlockResponse, error)
}

func RegisterMsgServer(s gogogrpc.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_LockTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgLockTokens)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).LockTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mytc.lockup.v1.Msg/LockTokens",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).LockTokens(ctx, req.(*MsgLockTokens))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Unlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUnlock)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Unlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mytc.lockup.v1.Msg/Unlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Unlock(ctx, req.(*MsgUnlock))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mytc.lockup.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LockTokens",
			Handler:    _Msg_LockTokens_Handler,
		},
		{
			MethodName: "Unlock",
			Handler:    _Msg_Unlock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mytc/lockup/v1/tx.proto",
}

func init() {
	proto.RegisterType((*MsgLockTokens)(nil), "mytc.lockup.v1.MsgLockTokens")
	proto.RegisterType((*MsgUnlock)(nil), "mytc.lockup.v1.MsgUnlock")
	proto.RegisterType((*MsgLockTokensResponse)(nil), "mytc.lockup.v1.MsgLockTokensResponse")
	proto.RegisterType((*MsgUnlockResponse)(nil), "mytc.lockup.v1.MsgUnlockResponse")
}
