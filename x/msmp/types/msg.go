package types

import (
	"context"

	"github.com/gogo/protobuf/proto"
	gogogrpc "github.com/gogo/protobuf/grpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc"
)

const (
	TypeMsgCollectFee          = "collect_fee"
	TypeMsgDistributeRewards   = "distribute_rewards"
	TypeMsgClaimActivityPoints = "claim_activity_points"
)

var _ sdk.Msg = &MsgCollectFee{}

type MsgCollectFee struct {
	Sender string    `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Amount sdk.Coins `protobuf:"bytes,2,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
}

func NewMsgCollectFee(sender string, amount sdk.Coins) *MsgCollectFee {
	return &MsgCollectFee{Sender: sender, Amount: amount}
}

func (msg *MsgCollectFee) Route() string { return RouterKey }
func (msg *MsgCollectFee) Type() string  { return TypeMsgCollectFee }
func (msg *MsgCollectFee) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
func (msg *MsgCollectFee) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
func (msg *MsgCollectFee) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}
	return nil
}
func (msg *MsgCollectFee) ProtoMessage() {}
func (msg *MsgCollectFee) Reset()        { *msg = MsgCollectFee{} }
func (msg *MsgCollectFee) String() string { return "MsgCollectFee" }

var _ sdk.Msg = &MsgDistributeRewards{}

type MsgDistributeRewards struct {
	Distributor string `protobuf:"bytes,1,opt,name=distributor,proto3" json:"distributor,omitempty"`
}

func NewMsgDistributeRewards(distributor string) *MsgDistributeRewards {
	return &MsgDistributeRewards{Distributor: distributor}
}

func (msg *MsgDistributeRewards) Route() string { return RouterKey }
func (msg *MsgDistributeRewards) Type() string  { return TypeMsgDistributeRewards }
func (msg *MsgDistributeRewards) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Distributor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgDistributeRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
func (msg *MsgDistributeRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Distributor)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid distributor address (%s)", err)
	}
	return nil
}
func (msg *MsgDistributeRewards) ProtoMessage() {}
func (msg *MsgDistributeRewards) Reset()        { *msg = MsgDistributeRewards{} }
func (msg *MsgDistributeRewards) String() string { return "MsgDistributeRewards" }

var _ sdk.Msg = &MsgClaimActivityPoints{}

type MsgClaimActivityPoints struct {
	Sender string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Points uint64 `protobuf:"varint,2,opt,name=points,proto3" json:"points,omitempty"`
}

func NewMsgClaimActivityPoints(sender string, points uint64) *MsgClaimActivityPoints {
	return &MsgClaimActivityPoints{Sender: sender, Points: points}
}

func (msg *MsgClaimActivityPoints) Route() string { return RouterKey }
func (msg *MsgClaimActivityPoints) Type() string  { return TypeMsgClaimActivityPoints }
func (msg *MsgClaimActivityPoints) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgClaimActivityPoints) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
func (msg *MsgClaimActivityPoints) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Points == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "points must be > 0")
	}
	return nil
}
func (msg *MsgClaimActivityPoints) ProtoMessage() {}
func (msg *MsgClaimActivityPoints) Reset()        { *msg = MsgClaimActivityPoints{} }
func (msg *MsgClaimActivityPoints) String() string { return "MsgClaimActivityPoints" }

// Response types
type MsgCollectFeeResponse struct{}
func (m *MsgCollectFeeResponse) Reset()        { *m = MsgCollectFeeResponse{} }
func (m *MsgCollectFeeResponse) String() string { return "MsgCollectFeeResponse" }
func (m *MsgCollectFeeResponse) ProtoMessage() {}

type MsgDistributeRewardsResponse struct{}
func (m *MsgDistributeRewardsResponse) Reset()        { *m = MsgDistributeRewardsResponse{} }
func (m *MsgDistributeRewardsResponse) String() string { return "MsgDistributeRewardsResponse" }
func (m *MsgDistributeRewardsResponse) ProtoMessage() {}

type MsgClaimActivityPointsResponse struct{}
func (m *MsgClaimActivityPointsResponse) Reset()        { *m = MsgClaimActivityPointsResponse{} }
func (m *MsgClaimActivityPointsResponse) String() string { return "MsgClaimActivityPointsResponse" }
func (m *MsgClaimActivityPointsResponse) ProtoMessage() {}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	CollectFee(context.Context, *MsgCollectFee) (*MsgCollectFeeResponse, error)
	DistributeRewards(context.Context, *MsgDistributeRewards) (*MsgDistributeRewardsResponse, error)
	ClaimActivityPoints(context.Context, *MsgClaimActivityPoints) (*MsgClaimActivityPointsResponse, error)
}

func RegisterMsgServer(s gogogrpc.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_CollectFee_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCollectFee)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CollectFee(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mytc.msmp.v1.Msg/CollectFee"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CollectFee(ctx, req.(*MsgCollectFee))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DistributeRewards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDistributeRewards)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DistributeRewards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mytc.msmp.v1.Msg/DistributeRewards"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DistributeRewards(ctx, req.(*MsgDistributeRewards))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClaimActivityPoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimActivityPoints)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimActivityPoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mytc.msmp.v1.Msg/ClaimActivityPoints"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimActivityPoints(ctx, req.(*MsgClaimActivityPoints))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mytc.msmp.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CollectFee", Handler: _Msg_CollectFee_Handler},
		{MethodName: "DistributeRewards", Handler: _Msg_DistributeRewards_Handler},
		{MethodName: "ClaimActivityPoints", Handler: _Msg_ClaimActivityPoints_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mytc/msmp/v1/tx.proto",
}

func init() {
	proto.RegisterType((*MsgCollectFee)(nil), "mytc.msmp.v1.MsgCollectFee")
	proto.RegisterType((*MsgDistributeRewards)(nil), "mytc.msmp.v1.MsgDistributeRewards")
	proto.RegisterType((*MsgClaimActivityPoints)(nil), "mytc.msmp.v1.MsgClaimActivityPoints")
	proto.RegisterType((*MsgCollectFeeResponse)(nil), "mytc.msmp.v1.MsgCollectFeeResponse")
	proto.RegisterType((*MsgDistributeRewardsResponse)(nil), "mytc.msmp.v1.MsgDistributeRewardsResponse")
	proto.RegisterType((*MsgClaimActivityPointsResponse)(nil), "mytc.msmp.v1.MsgClaimActivityPointsResponse")
}
