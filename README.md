<img width="600" height="600" alt="MYTC" src="https://github.com/user-attachments/assets/a6312b64-8441-4ae2-b753-a158f0db4fe2" />
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

MIT License

Copyright (c) 2026 Mytherra

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

