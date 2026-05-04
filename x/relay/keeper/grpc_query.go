package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mytherra/mytc/x/relay/types"
)

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{Keeper: k}
}

var _ types.QueryServer = queryServer{}

// Endpoint returns a single endpoint by operator address.
func (q queryServer) Endpoint(goCtx context.Context, req *types.QueryEndpointRequest) (*types.QueryEndpointResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(req.OperatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	ep, found := q.GetEndpoint(ctx, valAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "endpoint not found")
	}
	return &types.QueryEndpointResponse{Endpoint: &ep}, nil
}

// Endpoints returns every registered endpoint, optionally filtered by
// staleness. Pagination is intentionally simple — we iterate the full set
// (validator counts in the dozens at most for the foreseeable future) and
// slice in memory rather than wiring full Cosmos pagination.
func (q queryServer) Endpoints(goCtx context.Context, req *types.QueryEndpointsRequest) (*types.QueryEndpointsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	now := ctx.BlockTime().Unix()
	maxStale := req.MaxStaleSeconds

	out := make([]*types.RelayEndpoint, 0)
	q.IterateEndpoints(ctx, func(ep types.RelayEndpoint) bool {
		if maxStale > 0 && (now-ep.LastHeartbeat) > maxStale {
			return false
		}
		// Copy to avoid aliasing the loop var.
		copy := ep
		out = append(out, &copy)
		return false
	})

	return &types.QueryEndpointsResponse{Endpoints: out}, nil
}
