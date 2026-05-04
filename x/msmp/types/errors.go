package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInsufficientPoints = sdkerrors.Register(ModuleName, 1, "insufficient activity points")
	ErrInvalidFeeSplit    = sdkerrors.Register(ModuleName, 2, "invalid fee split configuration")
)
