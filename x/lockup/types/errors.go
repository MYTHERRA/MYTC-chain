package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/lockup module sentinel errors
var (
	ErrLockNotFound   = sdkerrors.Register(ModuleName, 1100, "lock not found")
	ErrNotLockOwner   = sdkerrors.Register(ModuleName, 1101, "not lock owner")
	ErrLockNotMatured = sdkerrors.Register(ModuleName, 1102, "lock not matured")
)
