package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/mytherra/mytc/x/lockup/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) LockTokens(goCtx context.Context, msg *types.MsgLockTokens) (*types.MsgLockTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	var duration time.Duration
	if msg.Duration != nil {
		duration, err = gogotypes.DurationFromProto(msg.Duration)
		if err != nil {
			return nil, err
		}
	}
	lockID, err := k.Keeper.LockTokens(ctx, owner, msg.Amount, duration)
	if err != nil {
		return nil, err
	}

	return &types.MsgLockTokensResponse{
		LockId: lockID,
	}, nil
}

func (k msgServer) Unlock(goCtx context.Context, msg *types.MsgUnlock) (*types.MsgUnlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	err = k.Keeper.UnlockTokens(ctx, owner, msg.LockId)
	if err != nil {
		return nil, err
	}

	return &types.MsgUnlockResponse{}, nil
}
