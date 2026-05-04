package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mytherra/mytc/x/msmp/types"
)

const (
	umytcDenom = "umytc" // micro MYTC denomination
)

// CollectFee moves messaging fees from a user into the MSMP module account.
func (k Keeper) CollectFee(ctx sdk.Context, sender sdk.AccAddress, amount sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount)
}

// DistributeRewards performs the epoch-based fee split:
//   80% to stakers, 20% to foundation (capped at 5% of total supply).
func (k Keeper) DistributeRewards(ctx sdk.Context) error {
	params := k.GetParams(ctx)

	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := k.bankKeeper.GetBalance(ctx, moduleAddr, umytcDenom)
	if moduleBalance.IsZero() {
		return nil // nothing to distribute
	}

	totalFees := moduleBalance.Amount
	stakerShare := totalFees.Mul(sdk.NewInt(int64(params.FeeSplitStaker))).Quo(sdk.NewInt(100))
	foundationShare := totalFees.Sub(stakerShare)

	// --- Foundation safety valve ---
	foundationAddr := k.accountKeeper.GetModuleAddress("foundation")
	if foundationAddr != nil && foundationShare.IsPositive() {
		foundationBalance := k.bankKeeper.GetBalance(ctx, foundationAddr, umytcDenom)
		totalSupply := k.bankKeeper.GetSupply(ctx, umytcDenom)

		capLimit := totalSupply.Amount.Mul(sdk.NewInt(int64(params.FoundationCapPercent))).Quo(sdk.NewInt(100))
		projected := foundationBalance.Amount.Add(foundationShare)

		if projected.GT(capLimit) {
			// Redirect excess to stakers instead of breaching the cap
			excess := projected.Sub(capLimit)
			stakerShare = stakerShare.Add(excess)
			foundationShare = capLimit.Sub(foundationBalance.Amount)
			if foundationShare.IsNegative() {
				foundationShare = sdk.ZeroInt()
			}
		}
	}

	// --- Staker distribution (MVP placeholder) ---
	// TODO: Query x/lockup for active lockups and distribute proportionally.
	// For now the staker share remains in the module account as a reserve.
	if stakerShare.IsPositive() {
		// In production: iterate lockups, calculate weights, send coins.
		_ = stakerShare
	}

	// --- Foundation transfer ---
	if foundationShare.IsPositive() && foundationAddr != nil {
		coins := sdk.NewCoins(sdk.NewCoin(umytcDenom, foundationShare))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, foundationAddr, coins); err != nil {
			return err
		}
	}

	return nil
}

// AddActivityPoints credits trust points to a user for sending messages.
func (k Keeper) AddActivityPoints(ctx sdk.Context, sender sdk.AccAddress, points uint64) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.ActivityPointsKey, sender...)

	var current uint64
	bz := store.Get(key)
	if bz != nil {
		current = sdk.BigEndianToUint64(bz)
	}
	current += points
	store.Set(key, sdk.Uint64ToBigEndian(current))
}

// GetActivityPoints returns the accumulated trust points for a user.
func (k Keeper) GetActivityPoints(ctx sdk.Context, sender sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := append(types.ActivityPointsKey, sender...)
	bz := store.Get(key)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// ClaimActivityPoints mints a small bonus reward for redeemed trust points.
// TODO: define conversion rate in parameters.
func (k Keeper) ClaimActivityPoints(ctx sdk.Context, sender sdk.AccAddress, points uint64) error {
	store := ctx.KVStore(k.storeKey)
	key := append(types.ActivityPointsKey, sender...)

	current := k.GetActivityPoints(ctx, sender)
	if current < points {
		return types.ErrInsufficientPoints
	}

	current -= points
	store.Set(key, sdk.Uint64ToBigEndian(current))

	// MVP: 1 point = 0.0001 MYTC (hardcoded for now)
	reward := sdk.NewDec(int64(points)).Mul(sdk.NewDecWithPrec(1, 4)) // 0.0001
	rewardInt := reward.TruncateInt()
	if rewardInt.IsPositive() {
		coins := sdk.NewCoins(sdk.NewCoin(umytcDenom, rewardInt))
		if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
			return err
		}
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coins); err != nil {
			return err
		}
	}
	return nil
}
