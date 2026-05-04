package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "lockup"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_lockup"
)

var (
	KeyLastLockID        = []byte{0x01}
	KeyLockPrefix        = []byte{0x02}
	KeyLockByOwnerPrefix = []byte{0x03}
)

func GetLockKey(id uint64) []byte {
	return append(KeyLockPrefix, sdk.Uint64ToBigEndian(id)...)
}

func GetLockByOwnerKey(owner sdk.AccAddress, id uint64) []byte {
	return append(append(KeyLockByOwnerPrefix, owner.Bytes()...), sdk.Uint64ToBigEndian(id)...)
}
