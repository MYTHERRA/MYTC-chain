package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCollectFee{}, "msmp/CollectFee", nil)
	cdc.RegisterConcrete(&MsgDistributeRewards{}, "msmp/DistributeRewards", nil)
	cdc.RegisterConcrete(&MsgClaimActivityPoints{}, "msmp/ClaimActivityPoints", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCollectFee{},
		&MsgDistributeRewards{},
		&MsgClaimActivityPoints{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
