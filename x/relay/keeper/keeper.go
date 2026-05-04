package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/mytherra/mytc/x/relay/types"
)

// MinHeartbeatInterval is the minimum number of seconds between consecutive
// heartbeats from the same validator. Sending earlier is rejected. Hardcoded
// for now; promote to module Params once the module needs governance-tunable
// behaviour.
const MinHeartbeatInterval int64 = 60

type Keeper struct {
	cdc      codec.Codec
	storeKey storetypes.StoreKey

	stakingKeeper types.StakingKeeper
}

func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, sk types.StakingKeeper) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		stakingKeeper: sk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
