package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/mytherra/mytc/x/lockup/types"
)

func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the lockup module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cmdQueryLock(),
		cmdQueryLocksByOwner(),
	)

	return cmd
}

func cmdQueryLock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock [lock-id]",
		Short: "Query a lock by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			qc := types.NewQueryClient(clientCtx)
			res, err := qc.Lock(context.Background(), &types.QueryLockRequest{LockId: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryLocksByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locks-by-owner [owner]",
		Short: "Query locks by owner address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			qc := types.NewQueryClient(clientCtx)
			res, err := qc.LocksByOwner(context.Background(), &types.QueryLocksByOwnerRequest{Owner: args[0]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
