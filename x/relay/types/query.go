package types

import (
	"context"

	"github.com/gogo/protobuf/proto"
	gogogrpc "github.com/gogo/protobuf/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc"
)

// QueryServer is the server API for Query service.
type QueryServer interface {
	Endpoint(context.Context, *QueryEndpointRequest) (*QueryEndpointResponse, error)
	Endpoints(context.Context, *QueryEndpointsRequest) (*QueryEndpointsResponse, error)
}

func RegisterQueryServer(s gogogrpc.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Endpoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryEndpointRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Endpoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mytc.relay.v1.Query/Endpoint"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Endpoint(ctx, req.(*QueryEndpointRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Endpoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryEndpointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Endpoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/mytc.relay.v1.Query/Endpoints"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Endpoints(ctx, req.(*QueryEndpointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mytc.relay.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Endpoint", Handler: _Query_Endpoint_Handler},
		{MethodName: "Endpoints", Handler: _Query_Endpoints_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mytc/relay/v1/query.proto",
}

// ─── Request / Response types ──────────────────────────────────────────────

type QueryEndpointRequest struct {
	OperatorAddr string `protobuf:"bytes,1,opt,name=operator_addr,json=operatorAddr,proto3" json:"operator_addr"`
}

func (m *QueryEndpointRequest) Reset()         { *m = QueryEndpointRequest{} }
func (m *QueryEndpointRequest) String() string { return "QueryEndpointRequest" }
func (m *QueryEndpointRequest) ProtoMessage()  {}

type QueryEndpointResponse struct {
	Endpoint *RelayEndpoint `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
}

func (m *QueryEndpointResponse) Reset()         { *m = QueryEndpointResponse{} }
func (m *QueryEndpointResponse) String() string { return "QueryEndpointResponse" }
func (m *QueryEndpointResponse) ProtoMessage()  {}

type QueryEndpointsRequest struct {
	// MaxStaleSeconds: if > 0, only return endpoints whose LastHeartbeat is
	// within this window. 0 means no filter.
	MaxStaleSeconds int64              `protobuf:"varint,1,opt,name=max_stale_seconds,json=maxStaleSeconds,proto3" json:"max_stale_seconds"`
	Pagination      *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryEndpointsRequest) Reset()         { *m = QueryEndpointsRequest{} }
func (m *QueryEndpointsRequest) String() string { return "QueryEndpointsRequest" }
func (m *QueryEndpointsRequest) ProtoMessage()  {}

type QueryEndpointsResponse struct {
	Endpoints  []*RelayEndpoint    `protobuf:"bytes,1,rep,name=endpoints,proto3" json:"endpoints"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryEndpointsResponse) Reset()         { *m = QueryEndpointsResponse{} }
func (m *QueryEndpointsResponse) String() string { return "QueryEndpointsResponse" }
func (m *QueryEndpointsResponse) ProtoMessage()  {}

func init() {
	proto.RegisterType((*QueryEndpointRequest)(nil), "mytc.relay.v1.QueryEndpointRequest")
	proto.RegisterType((*QueryEndpointResponse)(nil), "mytc.relay.v1.QueryEndpointResponse")
	proto.RegisterType((*QueryEndpointsRequest)(nil), "mytc.relay.v1.QueryEndpointsRequest")
	proto.RegisterType((*QueryEndpointsResponse)(nil), "mytc.relay.v1.QueryEndpointsResponse")
}
