package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/mytherra/mytc/x/relay/types"
)

func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Relay-registry transactions",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(cmdRegisterRelay(), cmdUnregisterRelay(), cmdHeartbeat())
	return cmd
}

func cmdRegisterRelay() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [wss-url] [version]",
		Short: "Register your validator's relay endpoint URL on-chain",
		Long: `Register (or replace) your validator's WebSocket relay URL on-chain.
The signing account must be the operator-key of a currently bonded validator.

Example:
  mytcd tx relay register wss://relay.example.org:443 myt-relay-v1.0.0 \
    --from validator --chain-id mytc-1
`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			version := ""
			if len(args) >= 2 {
				version = args[1]
			}

			// Operator address derived from the from-account: signing account is
			// AccAddress(valoper). We rebuild the bech32 valoper from it.
			valAddr := sdk.ValAddress(clientCtx.FromAddress).String()

			msg := types.NewMsgRegisterRelay(valAddr, args[0], version)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdUnregisterRelay() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unregister",
		Short: "Remove your validator's relay endpoint from the on-chain registry",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			valAddr := sdk.ValAddress(clientCtx.FromAddress).String()
			msg := types.NewMsgUnregisterRelay(valAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdHeartbeat() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "heartbeat",
		Short: "Send a heartbeat for your validator's relay endpoint",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			valAddr := sdk.ValAddress(clientCtx.FromAddress).String()
			msg := types.NewMsgHeartbeat(valAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
