package types

import (
	"context"
	"net/url"
	"strings"

	"github.com/gogo/protobuf/proto"
	gogogrpc "github.com/gogo/protobuf/grpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc"
)

const (
	TypeMsgRegisterRelay   = "register_relay"
	TypeMsgUnregisterRelay = "unregister_relay"
	TypeMsgHeartbeat       = "heartbeat"

	// MaxWssUrlLength caps the stored URL to avoid abuse.
	MaxWssUrlLength = 256
	// MaxVersionLength caps the version string.
	MaxVersionLength = 64
)

var (
	_ sdk.Msg = &MsgRegisterRelay{}
	_ sdk.Msg = &MsgUnregisterRelay{}
	_ sdk.Msg = &MsgHeartbeat{}
)

// ─── MsgRegisterRelay ──────────────────────────────────────────────────────

// MsgRegisterRelay registers (or replaces) the calling validator's relay
// endpoint URL. The signer must be the operator-key of a currently bonded
// validator.
type MsgRegisterRelay struct {
	OperatorAddr string `protobuf:"bytes,1,opt,name=operator_addr,json=operatorAddr,proto3" json:"operator_addr"`
	WssUrl       string `protobuf:"bytes,2,opt,name=wss_url,json=wssUrl,proto3" json:"wss_url"`
	Version      string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
}

func NewMsgRegisterRelay(operator, wssUrl, version string) *MsgRegisterRelay {
	return &MsgRegisterRelay{OperatorAddr: operator, WssUrl: wssUrl, Version: version}
}

func (msg *MsgRegisterRelay) Route() string { return RouterKey }
func (msg *MsgRegisterRelay) Type() string  { return TypeMsgRegisterRelay }

func (msg *MsgRegisterRelay) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddr)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(valAddr)}
}

func (msg *MsgRegisterRelay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterRelay) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.OperatorAddr); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if err := validateWssUrl(msg.WssUrl); err != nil {
		return err
	}
	if len(msg.Version) > MaxVersionLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "version too long")
	}
	return nil
}

func (msg *MsgRegisterRelay) ProtoMessage() {}
func (msg *MsgRegisterRelay) Reset()        { *msg = MsgRegisterRelay{} }
func (msg *MsgRegisterRelay) String() string {
	return "MsgRegisterRelay{" + msg.OperatorAddr + " -> " + msg.WssUrl + "}"
}

type MsgRegisterRelayResponse struct{}

func (m *MsgRegisterRelayResponse) Reset()         { *m = MsgRegisterRelayResponse{} }
func (m *MsgRegisterRelayResponse) String() string { return "MsgRegisterRelayResponse" }
func (m *MsgRegisterRelayResponse) ProtoMessage()  {}

// ─── MsgUnregisterRelay ────────────────────────────────────────────────────

// MsgUnregisterRelay removes the calling validator's relay endpoint. Idempotent.
type MsgUnregisterRelay struct {
	OperatorAddr string `protobuf:"bytes,1,opt,name=operator_addr,json=operatorAddr,proto3" json:"operator_addr"`
}

func NewMsgUnregisterRelay(operator string) *MsgUnregisterRelay {
	return &MsgUnregisterRelay{OperatorAddr: operator}
}

func (msg *MsgUnregisterRelay) Route() string { return RouterKey }
func (msg *MsgUnregisterRelay) Type() string  { return TypeMsgUnregisterRelay }

func (msg *MsgUnregisterRelay) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddr)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(valAddr)}
}

func (msg *MsgUnregisterRelay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnregisterRelay) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.OperatorAddr); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	return nil
}

func (msg *MsgUnregisterRelay) ProtoMessage()  {}
func (msg *MsgUnregisterRelay) Reset()         { *msg = MsgUnregisterRelay{} }
func (msg *MsgUnregisterRelay) String() string { return "MsgUnregisterRelay{" + msg.OperatorAddr + "}" }

type MsgUnregisterRelayResponse struct{}

func (m *MsgUnregisterRelayResponse) Reset()         { *m = MsgUnregisterRelayResponse{} }
func (m *MsgUnregisterRelayResponse) String() string { return "MsgUnregisterRelayResponse" }
func (m *MsgUnregisterRelayResponse) ProtoMessage()  {}

// ─── MsgHeartbeat ──────────────────────────────────────────────────────────

// MsgHeartbeat updates the LastHeartbeat timestamp on an existing endpoint.
// Throttled — server rejects if sent too soon after the previous heartbeat
// (see Params.MinHeartbeatInterval — for now hardcoded server-side).
type MsgHeartbeat struct {
	OperatorAddr string `protobuf:"bytes,1,opt,name=operator_addr,json=operatorAddr,proto3" json:"operator_addr"`
}

func NewMsgHeartbeat(operator string) *MsgHeartbeat {
	return &MsgHeartbeat{OperatorAddr: operator}
}

func (msg *MsgHeartbeat) Route() string { return RouterKey }
func (msg *MsgHeartbeat) Type() string  { return TypeMsgHeartbeat }

func (msg *MsgHeartbeat) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddr)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(valAddr)}
}

func (msg *MsgHeartbeat) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgHeartbeat) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.OperatorAddr); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	return nil
}

func (msg *MsgHeartbeat) ProtoMessage()  {}
func (msg *MsgHeartbeat) Reset()         { *msg = MsgHeartbeat{} }
func (msg *MsgHeartbeat) String() string { return "MsgHeartbeat{" + msg.OperatorAddr + "}" }

type MsgHeartbeatResponse struct{}

func (m *MsgHeartbeatResponse) Reset()         { *m = MsgHeartbeatResponse{} }
func (m *MsgHeartbeatResponse) String() string { return "MsgHeartbeatResponse" }
func (m *MsgHeartbeatResponse) ProtoMessage()  {}

// ─── Helpers ───────────────────────────────────────────────────────────────

func validateWssUrl(raw string) error {
	if len(raw) == 0 || len(raw) > MaxWssUrlLength {
		return ErrInvalidWssUrl
	}
	if !strings.HasPrefix(raw, "wss://") {
		return ErrInvalidWssUrl
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return ErrInvalidWssUrl
	}
	return nil
}

// ─── gRPC Service Wiring (hand-written, mirrors x/lockup pattern) ──────────

type MsgServer interface {
	RegisterRelay(context.Context, *MsgRegisterRelay) (*MsgRegisterRelayResponse, error)
	UnregisterRelay(context.Context, *MsgUnregisterRelay) (*MsgUnregisterRelayResponse, error)
	Heartbeat(context.Context, *MsgHeartbeat) (*MsgHeartbeatResponse, error)
}

func RegisterMsgServer(s gogogrpc.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_RegisterRelay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterRelay)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterRelay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mytc.relay.v1.Msg/RegisterRelay"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterRelay(ctx, req.(*MsgRegisterRelay))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UnregisterRelay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUnregisterRelay)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UnregisterRelay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mytc.relay.v1.Msg/UnregisterRelay"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UnregisterRelay(ctx, req.(*MsgUnregisterRelay))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgHeartbeat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mytc.relay.v1.Msg/Heartbeat"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Heartbeat(ctx, req.(*MsgHeartbeat))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mytc.relay.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "RegisterRelay", Handler: _Msg_RegisterRelay_Handler},
		{MethodName: "UnregisterRelay", Handler: _Msg_UnregisterRelay_Handler},
		{MethodName: "Heartbeat", Handler: _Msg_Heartbeat_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mytc/relay/v1/tx.proto",
}

func init() {
	proto.RegisterType((*RelayEndpoint)(nil), "mytc.relay.v1.RelayEndpoint")
	proto.RegisterType((*MsgRegisterRelay)(nil), "mytc.relay.v1.MsgRegisterRelay")
	proto.RegisterType((*MsgRegisterRelayResponse)(nil), "mytc.relay.v1.MsgRegisterRelayResponse")
	proto.RegisterType((*MsgUnregisterRelay)(nil), "mytc.relay.v1.MsgUnregisterRelay")
	proto.RegisterType((*MsgUnregisterRelayResponse)(nil), "mytc.relay.v1.MsgUnregisterRelayResponse")
	proto.RegisterType((*MsgHeartbeat)(nil), "mytc.relay.v1.MsgHeartbeat")
	proto.RegisterType((*MsgHeartbeatResponse)(nil), "mytc.relay.v1.MsgHeartbeatResponse")
}
