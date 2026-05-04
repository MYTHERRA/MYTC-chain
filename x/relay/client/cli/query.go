package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/mytherra/mytc/x/relay/types"
)

func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the relay registry",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(cmdQueryEndpoint(), cmdQueryEndpoints())
	return cmd
}

func cmdQueryEndpoint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoint [valoper]",
		Short: "Query a single relay endpoint by validator operator address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := newQueryClient(clientCtx)
			res, err := qc.Endpoint(context.Background(), &types.QueryEndpointRequest{OperatorAddr: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func cmdQueryEndpoints() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoints [max-stale-seconds]",
		Short: "List all relay endpoints (optionally filter by max staleness in seconds)",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			var maxStale int64
			if len(args) >= 1 {
				v, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return err
				}
				maxStale = v
			}
			qc := newQueryClient(clientCtx)
			res, err := qc.Endpoints(context.Background(), &types.QueryEndpointsRequest{MaxStaleSeconds: maxStale})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// newQueryClient returns a QueryClient that talks to the running node via the
// gRPC interface that's exposed by client.Context.
func newQueryClient(ctx client.Context) types.QueryClient {
	return types.NewQueryClient(ctx)
}
