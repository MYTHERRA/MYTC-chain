# Running a MYTC Validator

This guide walks you through joining the **`mytherra-1`** Mytherra Communication chain as a validator. You'll sync a node, create a validator, and (optionally) register a relay endpoint that powers MYT Messenger.

> ⚠️ **Status:** Early federation. Currently only the Mytherra Foundation validator runs. We're actively recruiting 5–10 additional validators. If you're reading this and want to join, **reach out before starting** — you'll need MYTC tokens from us to self-bond.

---

## What you'll need

- **Hardware**: Any cheap Linux VPS works. Minimums: 2 vCPU, 4 GB RAM, 50 GB SSD, stable network. €5–10/month is enough.
- **Software**: Ubuntu 22.04 / Debian 12 / similar. Go 1.19+ installed.
- **MYTC tokens for self-bond** (currently provided by the foundation; later via dedicated faucet).
- **Open ports**: 26656 (P2P) inbound, optionally 26657 (RPC) and 1317 (REST) if you want to expose query APIs.

---

## 1. Install Go and git

```sh
sudo apt update
sudo apt install -y git build-essential

# Go 1.22 (works fine for the v0.45 SDK)
curl -L https://go.dev/dl/go1.22.2.linux-amd64.tar.gz | sudo tar -C /usr/local -xzf -
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile.d/go.sh
source /etc/profile.d/go.sh
go version  # should print "go version go1.22.2 ..."
```

## 2. Clone and build mytcd

```sh
git clone https://github.com/MYTHERRA/MYTC-chain.git
cd MYTC-chain
go build -o build/mytcd ./cmd/mytcd
sudo cp build/mytcd /usr/local/bin/
mytcd version  # smoke test
```

## 3. Initialize your node

Pick a moniker — this is the public name people will see for your validator.

```sh
mytcd init "your-moniker-here" --chain-id mytherra-1 --home ~/.mytc
```

This creates `~/.mytc/config/` with default `genesis.json`, `config.toml`, `app.toml`, and a generated `priv_validator_key.json` and `node_key.json`. **Back up `priv_validator_key.json` and `node_key.json`** — losing them means losing your validator identity.

## 4. Replace genesis.json with the live one

The init step creates a stub genesis. You need the real `mytherra-1` genesis:

```sh
# From the running mainnet node — ask the foundation operator for a copy,
# or download from a public mirror (TBD URL).
curl -o ~/.mytc/config/genesis.json https://example.org/mytherra-1-genesis.json

# Verify the hash matches what other validators have:
sha256sum ~/.mytc/config/genesis.json
```

## 5. Configure persistent peers

Edit `~/.mytc/config/config.toml`. Find `persistent_peers` under `[p2p]`, set it to:

```
persistent_peers = "5c67756024ba0566e5408330e9055b0009bae23a@217.154.114.75:26656"
```

(That's the Mytherra Foundation node — current bootstrap peer. Will expand as more validators join.)

## 6. Start the node and wait for sync

```sh
mytcd start --home ~/.mytc
```

You'll see "Replay: Done" and then it starts catching up blocks. Initial sync is fast (chain is small). Watch for the height in the logs to reach the live tip — query mainnet to compare:

```sh
curl -s https://explorer.mytherrablockchain.org/mytc-api/blocks/latest | grep '"height"'
```

When your local node's height matches, sync is complete. Stop with Ctrl+C and set up systemd (next step) for production.

## 7. Run as a systemd service

`/etc/systemd/system/mytc.service`:

```ini
[Unit]
Description=Mytherra MYTC Validator
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/mytcd start --home /root/.mytc
Restart=on-failure
RestartSec=20
LimitNOFILE=65535
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

```sh
sudo systemctl daemon-reload
sudo systemctl enable --now mytc.service
sudo journalctl -u mytc.service -f  # tail logs
```

## 8. Create a wallet key and get MYTC

```sh
# Create your operator key (the one that controls validator commission, etc.)
mytcd keys add validator --keyring-backend test --home ~/.mytc

# This prints a mnemonic — WRITE IT DOWN, store it offline. It's the only
# recovery path for your validator's funds.

# See your address:
mytcd keys show validator -a --keyring-backend test --home ~/.mytc
```

**Send your address to the foundation contact** — they'll send you the MYTC tokens you need to self-bond. (Minimum self-bond is currently 1,000,000 umytc / 1 MYTC, but in practice you'll want more for credibility.)

Once you've received them, verify:

```sh
mytcd query bank balances $(mytcd keys show validator -a --keyring-backend test --home ~/.mytc) \
  --node tcp://127.0.0.1:26657
```

## 9. Become a validator

```sh
mytcd tx staking create-validator \
  --amount=1000000000umytc \
  --pubkey=$(mytcd tendermint show-validator --home ~/.mytc) \
  --moniker="your-moniker-here" \
  --website="https://yourwebsite.example" \
  --details="A short description of who you are" \
  --commission-rate=0.10 \
  --commission-max-rate=0.20 \
  --commission-max-change-rate=0.01 \
  --min-self-delegation=1 \
  --from validator \
  --keyring-backend test \
  --home ~/.mytc \
  --chain-id mytherra-1 \
  --node tcp://127.0.0.1:26657 \
  --fees 200umytc \
  -y
```

Verify within ~10 seconds:

```sh
mytcd query staking validators --node tcp://127.0.0.1:26657 | grep -A1 your-moniker-here
```

You should see your moniker in the output with `status: BOND_STATUS_BONDED`. **Congratulations, you're now signing blocks.**

## 10. (Optional but encouraged) Register a relay endpoint

This is the part that makes MYTC actually *useful* beyond payments — your validator becomes a transport node for MYT Messenger.

You'll need to run a `myt-relay` server (separate from `mytcd`). Setup is in [github.com/MYTHERRA/MYT-Messenger/docs/relay-setup.md](https://github.com/MYTHERRA/MYT-Messenger) (TODO). Once your relay is running at e.g. `wss://relay.your-domain.example/ws`, register it on-chain:

```sh
mytcd tx relay register wss://relay.your-domain.example/ws v1.0.0 \
  --from validator \
  --keyring-backend test \
  --home ~/.mytc \
  --chain-id mytherra-1 \
  --node tcp://127.0.0.1:26657 \
  --fees 200umytc \
  -y
```

The MYT Messenger client picks this up automatically on its next reconnect cycle (within 30 minutes due to the registry cache TTL).

To keep your endpoint visible in the active set, send a heartbeat occasionally (e.g. via cron, hourly):

```sh
mytcd tx relay heartbeat \
  --from validator --keyring-backend test --home ~/.mytc \
  --chain-id mytherra-1 --node tcp://127.0.0.1:26657 \
  --fees 100umytc -y
```

To remove it (e.g. when decommissioning):

```sh
mytcd tx relay unregister --from validator --keyring-backend test --home ~/.mytc \
  --chain-id mytherra-1 --node tcp://127.0.0.1:26657 --fees 100umytc -y
```

## Verify everything

Single endpoint:

```sh
mytcd query relay endpoint $(mytcd keys show validator --bech val -a --keyring-backend test --home ~/.mytc) \
  --node tcp://127.0.0.1:26657
```

All endpoints in the federation:

```sh
mytcd query relay endpoints --node tcp://127.0.0.1:26657
```

Or via REST (any API node):

```sh
curl -s https://explorer.mytherrablockchain.org/mytc-api/mytc/relay/v1/endpoints | jq
```

---

## Operational tips

- **Backups**: `~/.mytc/config/priv_validator_key.json` and `~/.mytc/config/node_key.json`. Lose these = lose validator. Encrypt at rest.
- **Monitoring**: tail `journalctl -u mytc.service`. Watch for "missed signing" warnings — that means your node was offline during a block, which slowly slashes your stake if it persists.
- **Updates**: when the chain announces a software upgrade, you'll need to swap the `mytcd` binary at the upgrade height. Watch the chain governance forum (TBD) for upgrade proposals.

## Getting help

- Open an issue at [MYTHERRA/MYTC-chain](https://github.com/MYTHERRA/MYTC-chain/issues)
- Mytherra community Discord (TBD link)
