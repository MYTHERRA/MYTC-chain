package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values.
const (
	DefaultFeeSplitStaker       uint32 = 80
	DefaultFeeSplitFoundation   uint32 = 20
	DefaultFoundationCapPercent uint32 = 5
)

// Parameter store keys.
var (
	KeyFeeSplitStaker       = []byte("FeeSplitStaker")
	KeyFeeSplitFoundation   = []byte("FeeSplitFoundation")
	KeyFoundationCapPercent = []byte("FoundationCapPercent")
)

// Params holds the on-chain parameters for the MSMP module.
type Params struct {
	FeeSplitStaker       uint32 `json:"fee_split_staker" yaml:"fee_split_staker"`
	FeeSplitFoundation   uint32 `json:"fee_split_foundation" yaml:"fee_split_foundation"`
	FoundationCapPercent uint32 `json:"foundation_cap_percent" yaml:"foundation_cap_percent"`
}

// DefaultParams returns default msmp parameters.
func DefaultParams() Params {
	return Params{
		FeeSplitStaker:       DefaultFeeSplitStaker,
		FeeSplitFoundation:   DefaultFeeSplitFoundation,
		FoundationCapPercent: DefaultFoundationCapPercent,
	}
}

// ParamKeyTable for msmp module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements paramtypes.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFeeSplitStaker, &p.FeeSplitStaker, validateFeeSplit),
		paramtypes.NewParamSetPair(KeyFeeSplitFoundation, &p.FeeSplitFoundation, validateFeeSplit),
		paramtypes.NewParamSetPair(KeyFoundationCapPercent, &p.FoundationCapPercent, validateFeeSplit),
	}
}

func validateFeeSplit(i interface{}) error {
	val, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if val > 100 {
		return fmt.Errorf("fee split must be <= 100")
	}
	return nil
}
