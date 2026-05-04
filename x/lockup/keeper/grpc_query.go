package keeper

import (
	"context"
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/mytherra/mytc/x/lockup/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}

func (k queryServer) Lock(goCtx context.Context, req *types.QueryLockRequest) (*types.QueryLockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	lock, found := k.GetLock(ctx, req.LockId)
	if !found {
		return nil, status.Error(codes.NotFound, "lock not found")
	}

	return &types.QueryLockResponse{Lock: &lock}, nil
}

func (k queryServer) LocksByOwner(goCtx context.Context, req *types.QueryLocksByOwnerRequest) (*types.QueryLocksByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid owner address")
	}

	store := ctx.KVStore(k.storeKey)
	lockStore := prefix.NewStore(store, append(types.KeyLockByOwnerPrefix, owner.Bytes()...))

	var locks []*types.Lock
	pageRes, err := query.Paginate(lockStore, req.Pagination, func(key []byte, value []byte) error {
		// Key is just the lock ID (8 bytes)
		if len(key) < 8 {
			return nil
		}
		lockID := binary.BigEndian.Uint64(key)
		lock, found := k.GetLock(ctx, lockID)
		if found {
			lockCopy := lock
			locks = append(locks, &lockCopy)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLocksByOwnerResponse{Locks: locks, Pagination: pageRes}, nil
}
