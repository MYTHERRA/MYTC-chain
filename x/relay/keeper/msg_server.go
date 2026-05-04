package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mytherra/mytc/x/relay/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

var _ types.MsgServer = msgServer{}

// RegisterRelay registers (or replaces) the calling validator's relay endpoint.
func (k msgServer) RegisterRelay(goCtx context.Context, msg *types.MsgRegisterRelay) (*types.MsgRegisterRelayResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddr)
	if err != nil {
		return nil, err
	}

	if !k.IsBondedValidator(ctx, valAddr) {
		return nil, types.ErrNotBondedValidator
	}

	now := ctx.BlockTime().Unix()

	// Preserve the original RegisteredAt if this is a re-register.
	registeredAt := now
	if existing, found := k.GetEndpoint(ctx, valAddr); found {
		registeredAt = existing.RegisteredAt
	}

	ep := types.RelayEndpoint{
		OperatorAddr:  msg.OperatorAddr,
		WssUrl:        msg.WssUrl,
		Version:       msg.Version,
		LastHeartbeat: now,
		RegisteredAt:  registeredAt,
	}
	k.SetEndpoint(ctx, ep)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRegisterRelay,
		sdk.NewAttribute(types.AttributeKeyOperator, msg.OperatorAddr),
		sdk.NewAttribute(types.AttributeKeyWssUrl, msg.WssUrl),
		sdk.NewAttribute(types.AttributeKeyVersion, msg.Version),
	))

	return &types.MsgRegisterRelayResponse{}, nil
}

// UnregisterRelay removes the calling validator's endpoint. Idempotent.
func (k msgServer) UnregisterRelay(goCtx context.Context, msg *types.MsgUnregisterRelay) (*types.MsgUnregisterRelayResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddr)
	if err != nil {
		return nil, err
	}

	k.DeleteEndpoint(ctx, valAddr)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUnregisterRelay,
		sdk.NewAttribute(types.AttributeKeyOperator, msg.OperatorAddr),
	))

	return &types.MsgUnregisterRelayResponse{}, nil
}

// Heartbeat updates the LastHeartbeat timestamp on the existing endpoint.
// Throttled by MinHeartbeatInterval to prevent spam.
func (k msgServer) Heartbeat(goCtx context.Context, msg *types.MsgHeartbeat) (*types.MsgHeartbeatResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddr)
	if err != nil {
		return nil, err
	}

	if !k.IsBondedValidator(ctx, valAddr) {
		return nil, types.ErrNotBondedValidator
	}

	ep, found := k.GetEndpoint(ctx, valAddr)
	if !found {
		return nil, types.ErrEndpointNotFound
	}

	now := ctx.BlockTime().Unix()
	if now-ep.LastHeartbeat < MinHeartbeatInterval {
		return nil, types.ErrHeartbeatTooSoon
	}

	ep.LastHeartbeat = now
	k.SetEndpoint(ctx, ep)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeHeartbeat,
		sdk.NewAttribute(types.AttributeKeyOperator, msg.OperatorAddr),
	))

	return &types.MsgHeartbeatResponse{}, nil
}
