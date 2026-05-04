package types

import (
	"net/url"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// sdk.Msg type assertions — concrete types live in tx.pb.go.
var (
	_ sdk.Msg = &MsgRegisterRelay{}
	_ sdk.Msg = &MsgUnregisterRelay{}
	_ sdk.Msg = &MsgHeartbeat{}
)

const (
	TypeMsgRegisterRelay   = "register_relay"
	TypeMsgUnregisterRelay = "unregister_relay"
	TypeMsgHeartbeat       = "heartbeat"

	MaxWssUrlLength  = 256
	MaxVersionLength = 64
)

// ─── Constructors ──────────────────────────────────────────────────────────

func NewMsgRegisterRelay(operator, wssUrl, version string) *MsgRegisterRelay {
	return &MsgRegisterRelay{OperatorAddr: operator, WssUrl: wssUrl, Version: version}
}

func NewMsgUnregisterRelay(operator string) *MsgUnregisterRelay {
	return &MsgUnregisterRelay{OperatorAddr: operator}
}

func NewMsgHeartbeat(operator string) *MsgHeartbeat {
	return &MsgHeartbeat{OperatorAddr: operator}
}

// ─── MsgRegisterRelay sdk.Msg methods ──────────────────────────────────────

func (msg *MsgRegisterRelay) Route() string { return RouterKey }
func (msg *MsgRegisterRelay) Type() string  { return TypeMsgRegisterRelay }

func (msg *MsgRegisterRelay) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddr)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(valAddr)}
}

func (msg *MsgRegisterRelay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterRelay) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.OperatorAddr); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if err := validateWssUrl(msg.WssUrl); err != nil {
		return err
	}
	if len(msg.Version) > MaxVersionLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "version too long")
	}
	return nil
}

// ─── MsgUnregisterRelay sdk.Msg methods ────────────────────────────────────

func (msg *MsgUnregisterRelay) Route() string { return RouterKey }
func (msg *MsgUnregisterRelay) Type() string  { return TypeMsgUnregisterRelay }

func (msg *MsgUnregisterRelay) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddr)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(valAddr)}
}

func (msg *MsgUnregisterRelay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnregisterRelay) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.OperatorAddr); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	return nil
}

// ─── MsgHeartbeat sdk.Msg methods ──────────────────────────────────────────

func (msg *MsgHeartbeat) Route() string { return RouterKey }
func (msg *MsgHeartbeat) Type() string  { return TypeMsgHeartbeat }

func (msg *MsgHeartbeat) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.ValAddressFromBech32(msg.OperatorAddr)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(valAddr)}
}

func (msg *MsgHeartbeat) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgHeartbeat) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.OperatorAddr); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	return nil
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func validateWssUrl(raw string) error {
	if len(raw) == 0 || len(raw) > MaxWssUrlLength {
		return ErrInvalidWssUrl
	}
	if !strings.HasPrefix(raw, "wss://") {
		return ErrInvalidWssUrl
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return ErrInvalidWssUrl
	}
	return nil
}
