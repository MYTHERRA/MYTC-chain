package keeper

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/mytherra/mytc/x/lockup/types"
)

// GetLastLockID gets the last lock ID from the store
func (k Keeper) GetLastLockID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastLockID)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// SetLastLockID sets the last lock ID in the store
func (k Keeper) SetLastLockID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	store.Set(types.KeyLastLockID, bz)
}

// GetNextLockID returns the next lock ID and increments the counter
func (k Keeper) GetNextLockID(ctx sdk.Context) uint64 {
	id := k.GetLastLockID(ctx) + 1
	k.SetLastLockID(ctx, id)
	return id
}

// SetLock sets a lock in the store
func (k Keeper) SetLock(ctx sdk.Context, lock types.Lock) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalJSON(&lock)
	store.Set(types.GetLockKey(lock.Id), bz)

	// Set index
	owner, err := sdk.AccAddressFromBech32(lock.Owner)
	if err == nil {
		store.Set(types.GetLockByOwnerKey(owner, lock.Id), []byte{})
	}
}

// GetLock returns a lock from the store
func (k Keeper) GetLock(ctx sdk.Context, id uint64) (types.Lock, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLockKey(id))
	if bz == nil {
		return types.Lock{}, false
	}
	var lock types.Lock
	k.cdc.MustUnmarshalJSON(bz, &lock)
	return lock, true
}

// RemoveLock removes a lock from the store
func (k Keeper) RemoveLock(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	
	// Get lock to find owner for index deletion
	lock, found := k.GetLock(ctx, id)
	if found {
		owner, err := sdk.AccAddressFromBech32(lock.Owner)
		if err == nil {
			store.Delete(types.GetLockByOwnerKey(owner, id))
		}
	}

	store.Delete(types.GetLockKey(id))
}

// LockTokens locks tokens for a specific duration
func (k Keeper) LockTokens(ctx sdk.Context, owner sdk.AccAddress, amount sdk.Coin, duration time.Duration) (uint64, error) {
	// Send coins to module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(amount)); err != nil {
		return 0, err
	}

	lockID := k.GetNextLockID(ctx)
	endTime := ctx.BlockTime().Add(duration)
	endTS, err := gogotypes.TimestampProto(endTime)
	if err != nil {
		return 0, err
	}

	lock := types.Lock{
		Id:      lockID,
		Owner:   owner.String(),
		Amount:  amount,
		EndTime: endTS,
	}

	k.SetLock(ctx, lock)

	return lockID, nil
}

// UnlockTokens unlocks tokens if the lock has matured
func (k Keeper) UnlockTokens(ctx sdk.Context, owner sdk.AccAddress, lockID uint64) error {
	lock, found := k.GetLock(ctx, lockID)
	if !found {
		return types.ErrLockNotFound
	}

	if lock.Owner != owner.String() {
		return types.ErrNotLockOwner
	}

	if lock.EndTime == nil {
		return types.ErrLockNotMatured
	}
	endTime, err := gogotypes.TimestampFromProto(lock.EndTime)
	if err != nil {
		return err
	}
	if ctx.BlockTime().Before(endTime) {
		return types.ErrLockNotMatured
	}

	// Send coins back to owner
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(lock.Amount)); err != nil {
		return err
	}

	k.RemoveLock(ctx, lockID)

	return nil
}
