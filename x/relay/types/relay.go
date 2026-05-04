package types

// RelayEndpoint is a validator's published relay-server WSS URL.
// One per bonded validator. Lookups by validator operator address.
type RelayEndpoint struct {
	OperatorAddr  string `protobuf:"bytes,1,opt,name=operator_addr,json=operatorAddr,proto3" json:"operator_addr"`
	WssUrl        string `protobuf:"bytes,2,opt,name=wss_url,json=wssUrl,proto3" json:"wss_url"`
	Version       string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	LastHeartbeat int64  `protobuf:"varint,4,opt,name=last_heartbeat,json=lastHeartbeat,proto3" json:"last_heartbeat"`
	RegisteredAt  int64  `protobuf:"varint,5,opt,name=registered_at,json=registeredAt,proto3" json:"registered_at"`
}

func (m *RelayEndpoint) Reset()         { *m = RelayEndpoint{} }
func (m *RelayEndpoint) String() string { return m.OperatorAddr + " -> " + m.WssUrl }
func (m *RelayEndpoint) ProtoMessage()  {}
