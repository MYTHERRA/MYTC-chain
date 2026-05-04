# MYTC-chain

Mytherra Communication — Cosmos-SDK based PoS chain for the Mytherra ecosystem.

Powers payments and (in the future) the message-relay federation for **MYT Messenger**.

## Stack

- Cosmos SDK v0.45.16
- Tendermint v0.34.27
- Go 1.19+
- Module path: `github.com/mytherra/mytc`

## Existing modules

- `x/lockup` — token lockup / vesting
- `x/msmp` — Mytherra-specific module

## Building

```sh
go build -o build/mytcd ./cmd/mytcd
```

## Running a node

```sh
./build/mytcd start --home ~/.mytc
```

## Roadmap

- **`x/relay`** — on-chain registry of validator-operated message relays. First step toward making MYTC the transport layer for MYT Messenger (Phase A federation).

## License

TBD
