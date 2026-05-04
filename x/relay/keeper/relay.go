package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/mytherra/mytc/x/relay/types"
)

// SetEndpoint stores or replaces an endpoint record.
func (k Keeper) SetEndpoint(ctx sdk.Context, ep types.RelayEndpoint) {
	store := ctx.KVStore(k.storeKey)
	valAddr, err := sdk.ValAddressFromBech32(ep.OperatorAddr)
	if err != nil {
		// SetEndpoint should only be called with already-validated input.
		panic(err)
	}
	bz := k.cdc.MustMarshal(&ep)
	store.Set(types.GetEndpointKey(valAddr), bz)
}

// GetEndpoint returns the stored endpoint or (zero, false) when missing.
func (k Keeper) GetEndpoint(ctx sdk.Context, valoper sdk.ValAddress) (types.RelayEndpoint, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetEndpointKey(valoper))
	if bz == nil {
		return types.RelayEndpoint{}, false
	}
	var ep types.RelayEndpoint
	k.cdc.MustUnmarshal(bz, &ep)
	return ep, true
}

// DeleteEndpoint removes an endpoint. No-op if it doesn't exist.
func (k Keeper) DeleteEndpoint(ctx sdk.Context, valoper sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetEndpointKey(valoper))
}

// IterateEndpoints walks every stored endpoint. Stops if cb returns true.
func (k Keeper) IterateEndpoints(ctx sdk.Context, cb func(types.RelayEndpoint) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyEndpointPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ep types.RelayEndpoint
		k.cdc.MustUnmarshal(iter.Value(), &ep)
		if cb(ep) {
			return
		}
	}
}

// IsBondedValidator returns true if valoper currently corresponds to a
// validator in the bonded set.
func (k Keeper) IsBondedValidator(ctx sdk.Context, valoper sdk.ValAddress) bool {
	v := k.stakingKeeper.Validator(ctx, valoper)
	if v == nil {
		return false
	}
	return v.GetStatus() == stakingtypes.Bonded
}
