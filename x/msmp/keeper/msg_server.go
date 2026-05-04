package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mytherra/mytc/x/msmp/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CollectFee(goCtx context.Context, msg *types.MsgCollectFee) (*types.MsgCollectFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.CollectFee(ctx, sender, msg.Amount); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCollectFee,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	return &types.MsgCollectFeeResponse{}, nil
}

func (k msgServer) DistributeRewards(goCtx context.Context, msg *types.MsgDistributeRewards) (*types.MsgDistributeRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.DistributeRewards(ctx); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDistributeRewards,
			sdk.NewAttribute(types.AttributeKeyDistributor, msg.Distributor),
		),
	)

	return &types.MsgDistributeRewardsResponse{}, nil
}

func (k msgServer) ClaimActivityPoints(goCtx context.Context, msg *types.MsgClaimActivityPoints) (*types.MsgClaimActivityPointsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.ClaimActivityPoints(ctx, sender, msg.Points); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimActivityPoints,
			sdk.NewAttribute(types.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyPoints, sdk.NewUint(msg.Points).String()),
		),
	)

	return &types.MsgClaimActivityPointsResponse{}, nil
}
