# x/relay — open work

## Proto generation (BLOCKER for tx submission)

The hand-written types in `types/tx.go` and `types/query.go` follow the same
style as `x/lockup`, but the Cosmos SDK proto codec rejects them at tx
submission time:

```
cannot protobuf JSON encode unsupported type: *types.MsgRegisterRelay
*types.MsgRegisterRelay does not have a Descriptor() method: tx parse error
```

To fix:

1. Write `.proto` files under `proto/mytc/relay/v1/`:
   - `tx.proto` — Msg service + MsgRegisterRelay/Unregister/Heartbeat + responses
   - `query.proto` — Query service + Endpoint/Endpoints
   - `relay.proto` — RelayEndpoint message

2. Install plugins:
   - `protoc-gen-gocosmos` (regen-network) — ran into go.mod replace conflicts
     on the VPS; pin to a tagged release or use `buf` instead.
   - `protoc-gen-gogo`

3. Generate `.pb.go` files:
   ```sh
   protoc --proto_path=proto --gocosmos_out=plugins=interfacetype+grpc,\
     Mgogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto:./ \
     proto/mytc/relay/v1/*.proto
   ```

4. Delete the manual marshaling stanzas in `types/tx.go` and `types/query.go`
   (everything below the "gRPC Service Wiring" comments) — keep only the
   business-logic constructors (`NewMsgRegisterRelay`, `ValidateBasic`, etc.).

5. Same fix should be applied to `x/lockup` and `x/msmp` — they're affected by
   the same issue but it's never been triggered because nobody has tried to
   submit a lockup tx.

## After proto fix

- Genesis migration plan for mainnet (relay module gets added in a
  software-upgrade proposal)
- Multi-node testnet to test inter-validator independence
- Client-side: messenger reads endpoint set from Cosmos REST instead of
  hardcoded WSS URL

## Phase B / C / D (separate work)

- Multi-relay client (subscribe to N validators simultaneously)
- Inter-relay message gossip with dedup
- Fee routing to delivery validator
