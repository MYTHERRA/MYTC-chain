package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "relay"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for relay
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_relay"
)

// Store key prefixes
var (
	// KeyEndpointPrefix stores RelayEndpoint records keyed by valoper address.
	// Layout: 0x01 || valoper_bytes -> RelayEndpoint
	KeyEndpointPrefix = []byte{0x01}
)

// GetEndpointKey returns the store key for a single relay endpoint by validator
// operator address.
func GetEndpointKey(valoper sdk.ValAddress) []byte {
	return append(KeyEndpointPrefix, valoper.Bytes()...)
}
