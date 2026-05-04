package types

import (
	"github.com/gogo/protobuf/proto"
	gogotypes "github.com/gogo/protobuf/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Lock represents a single lockup position
type Lock struct {
	Id      uint64              `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	Owner   string              `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	Amount  sdk.Coin            `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
	EndTime *gogotypes.Timestamp `protobuf:"bytes,4,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
}

func (l *Lock) ProtoMessage() {}
func (l *Lock) Reset() { *l = Lock{} }
func (l *Lock) String() string { return "Lock" }

func init() {
	proto.RegisterType((*Lock)(nil), "mytc.lockup.v1.Lock")
}
