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
	Lock(context.Context, *QueryLockRequest) (*QueryLockResponse, error)
	LocksByOwner(context.Context, *QueryLocksByOwnerRequest) (*QueryLocksByOwnerResponse, error)
}

func RegisterQueryServer(s gogogrpc.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Lock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryLockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Lock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mytc.lockup.v1.Query/Lock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Lock(ctx, req.(*QueryLockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_LocksByOwner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryLocksByOwnerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).LocksByOwner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mytc.lockup.v1.Query/LocksByOwner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).LocksByOwner(ctx, req.(*QueryLocksByOwnerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mytc.lockup.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Lock",
			Handler:    _Query_Lock_Handler,
		},
		{
			MethodName: "LocksByOwner",
			Handler:    _Query_LocksByOwner_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mytc/lockup/v1/query.proto",
}

// Query Request/Response types
type QueryLockRequest struct {
	LockId uint64 `protobuf:"varint,1,opt,name=lock_id,json=lockId,proto3" json:"lock_id"`
}

func (m *QueryLockRequest) Reset()         { *m = QueryLockRequest{} }
func (m *QueryLockRequest) String() string { return "QueryLockRequest" }
func (m *QueryLockRequest) ProtoMessage()  {}

type QueryLockResponse struct {
	Lock *Lock `protobuf:"bytes,1,opt,name=lock,proto3" json:"lock,omitempty"`
}

func (m *QueryLockResponse) Reset()         { *m = QueryLockResponse{} }
func (m *QueryLockResponse) String() string { return "QueryLockResponse" }
func (m *QueryLockResponse) ProtoMessage()  {}

type QueryLocksByOwnerRequest struct {
	Owner      string             `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryLocksByOwnerRequest) Reset()         { *m = QueryLocksByOwnerRequest{} }
func (m *QueryLocksByOwnerRequest) String() string { return "QueryLocksByOwnerRequest" }
func (m *QueryLocksByOwnerRequest) ProtoMessage()  {}

type QueryLocksByOwnerResponse struct {
	Locks      []*Lock             `protobuf:"bytes,1,rep,name=locks,proto3" json:"locks"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryLocksByOwnerResponse) Reset()         { *m = QueryLocksByOwnerResponse{} }
func (m *QueryLocksByOwnerResponse) String() string { return "QueryLocksByOwnerResponse" }
func (m *QueryLocksByOwnerResponse) ProtoMessage()  {}

func init() {
	proto.RegisterType((*QueryLockRequest)(nil), "mytc.lockup.v1.QueryLockRequest")
	proto.RegisterType((*QueryLockResponse)(nil), "mytc.lockup.v1.QueryLockResponse")
	proto.RegisterType((*QueryLocksByOwnerRequest)(nil), "mytc.lockup.v1.QueryLocksByOwnerRequest")
	proto.RegisterType((*QueryLocksByOwnerResponse)(nil), "mytc.lockup.v1.QueryLocksByOwnerResponse")
}
