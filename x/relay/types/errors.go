package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNotBondedValidator = sdkerrors.Register(ModuleName, 2, "signer is not a currently bonded validator")
	ErrInvalidWssUrl      = sdkerrors.Register(ModuleName, 3, "invalid wss URL (must start with wss:// and have a host)")
	ErrEndpointNotFound   = sdkerrors.Register(ModuleName, 4, "relay endpoint not found")
	ErrHeartbeatTooSoon   = sdkerrors.Register(ModuleName, 5, "heartbeat sent too soon after the previous one")
)
