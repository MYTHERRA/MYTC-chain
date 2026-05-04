# Running a MYTC Validator

This guide walks you through joining the **`mytherra-1`** Mytherra Communication chain as a validator. You'll sync a node, create a validator, and (optionally) register a relay endpoint that powers MYT Messenger.

> ⚠️ **Status:** Early federation. As of writing, two validators run (Mytherra Foundation + foundation-2). We're actively recruiting independent operators. **Reach out before starting** — you'll need MYTC tokens from the foundation to self-bond.

---

## What you'll need

- **Hardware**: Any cheap Linux VPS works. Minimums: 2 vCPU, 4 GB RAM, 50 GB SSD, stable network. €5–10/month is enough.
- **Software**: Ubuntu 22.04 / Debian 12. Go 1.19+ installed.
- **MYTC tokens for self-bond** (currently provided by the foundation; later via dedicated faucet).
- **Open ports**: 26656 (P2P) inbound. Optionally 26657 (RPC) and 1317 (REST) if you want to expose query APIs.

> 💡 **If your VPS is behind a cloud-provider firewall (IONOS, Hetzner Cloud, AWS Security Group, etc.):** the OS-level UFW or iptables rules are *not enough*. You also have to open **TCP 26656 inbound** in the provider's web panel. Changes can take 5–10 minutes to propagate; sometimes a server power-cycle from the panel is needed before they take effect.

---

## 1. Install Go and git

```sh
sudo apt update
sudo apt install -y git build-essential

curl -L https://go.dev/dl/go1.22.2.linux-amd64.tar.gz | sudo tar -C /usr/local -xzf -
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile.d/go.sh
source /etc/profile.d/go.sh
go version  # should print "go version go1.22.2 ..."
```

## 2. Create a dedicated user (recommended)

Don't run the validator as root. Create a system user:

```sh
sudo useradd -r -m -d /home/mytc -s /bin/bash mytc
```

All subsequent commands assume you're working as `mytc`:

```sh
sudo -i -u mytc
```

## 3. Clone and build mytcd

```sh
git clone https://github.com/MYTHERRA/MYTC-chain.git
cd MYTC-chain
go build -o build/mytcd ./cmd/mytcd

# Install system-wide so systemd can find it
sudo cp build/mytcd /usr/local/bin/
mytcd version 2>&1 | head -1
```

## 4. Initialize your node

Pick a moniker — this is the public name people will see for your validator.

```sh
mytcd init "your-moniker-here" --chain-id mytherra-1 --home ~/.mytc
```

> ⚠️ **Backup `~/.mytc/config/priv_validator_key.json` and `~/.mytc/config/node_key.json` immediately.** Encrypt and store offline. If you lose `priv_validator_key.json`, you lose your validator (and any future double-sign with a re-created key gets you slashed).

## 5. Get the live genesis.json

```sh
curl -o ~/.mytc/config/genesis.json \
  https://mytmessenger.mytherrablockchain.org/mytherra-1-genesis.json

# Verify the hash:
sha256sum ~/.mytc/config/genesis.json
# Expected: 6d160fd45af10de839919d4d179633829afa32956b6362d60d8923ab489eafec
```

If the hash doesn't match — **stop**. Don't start with a tampered genesis. Reach out via GitHub issues.

## 6. Configure peers + your external address

Edit `~/.mytc/config/config.toml`:

```toml
# Find this in [p2p] section
persistent_peers = "5c67756024ba0566e5408330e9055b0009bae23a@217.154.114.75:26656"

# IMPORTANT — set this to YOUR public IP so other validators can dial in.
external_address = "YOUR.PUBLIC.IP:26656"
```

(`5c67…@217.154.114.75:26656` is the Mytherra Foundation bootstrap node. The list will grow as more validators join.)

## 7. Sync the chain

You have two options:

### 7a. Block sync from genesis (slow but pure)

Just start the node and let it replay every block:

```sh
mytcd start --home ~/.mytc
```

This works fine for a small chain. **However**: in some IONOS-to-IONOS / cross-cloud-provider setups, the block sync stalls with `invalid peer` errors due to packet fragmentation or MTU mismatches on big block-data messages. If you see `module=blockchain` errors and your height stays at 1, fall back to 7b.

### 7b. Snapshot sync (fast workaround)

Get a copy of the data directory from a healthy operator. Coordinate with the foundation contact — they can `rsync` you a snapshot over SSH:

```sh
# On YOUR validator machine (run before starting mytcd):
mytcd init "your-moniker-here" --chain-id mytherra-1 --home ~/.mytc
# (then preserve YOUR own keys before importing data)
cp ~/.mytc/config/priv_validator_key.json /tmp/my-priv.json
cp ~/.mytc/config/node_key.json /tmp/my-node-key.json
cp ~/.mytc/data/priv_validator_state.json /tmp/my-priv-state.json

# Foundation operator runs (on mainnet):
#   rsync -az /root/.mytc/data/ user@your.vps:/home/mytc/.mytc/data/

# Then on your machine, restore YOUR validator keys (must NOT be the same as another validator's):
cp /tmp/my-priv-state.json ~/.mytc/data/priv_validator_state.json
chown -R mytc:mytc ~/.mytc
```

> 🔒 **Critical**: never reuse another validator's `priv_validator_key.json`. Two nodes signing with the same key cause double-sign slashing. Always preserve and restore your own.

## 8. Run as a systemd service

`/etc/systemd/system/mytc.service`:

```ini
[Unit]
Description=Mytherra MYTC Validator
After=network.target

[Service]
Type=simple
User=mytc
Group=mytc
WorkingDirectory=/home/mytc
ExecStart=/usr/local/bin/mytcd start --home /home/mytc/.mytc
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

Wait until `catching_up: false` — then you're synced:

```sh
mytcd status --home ~/.mytc 2>&1 | grep -E 'latest_block_height|catching_up'
```

## 9. Create a wallet key and get MYTC

```sh
mytcd keys add validator --keyring-backend test --home ~/.mytc
```

> ⚠️ **Save the printed mnemonic offline.** That's the only recovery for this account's funds.

Get your address:

```sh
mytcd keys show validator -a --keyring-backend test --home ~/.mytc
```

**Send your address to the foundation contact** — they'll send the MYTC you need to self-bond. Recommended: at least 2,000,000,000 umytc (2,000 MYTC) so you have buffer for fees + future delegations.

Verify funds arrived:

```sh
mytcd query bank balances $(mytcd keys show validator -a --keyring-backend test --home ~/.mytc) \
  --node tcp://127.0.0.1:26657
```

## 10. Become a validator

```sh
mytcd tx staking create-validator \
  --amount=1000000000umytc \
  --pubkey=$(mytcd tendermint show-validator --home ~/.mytc) \
  --moniker="your-moniker-here" \
  --website="https://yourwebsite.example" \
  --details="Short description of who you are" \
  --commission-rate=0.10 \
  --commission-max-rate=0.20 \
  --commission-max-change-rate=0.01 \
  --min-self-delegation=1 \
  --from validator \
  --keyring-backend test \
  --home ~/.mytc \
  --chain-id mytherra-1 \
  --node tcp://127.0.0.1:26657 \
  --fees 1000umytc \
  -y
```

Verify within ~10 seconds:

```sh
mytcd query staking validators --node tcp://127.0.0.1:26657 \
  | grep -A2 your-moniker-here
```

You should see `status: BOND_STATUS_BONDED`. **Congratulations — you're signing blocks.** Confirm with:

```sh
mytcd status --home ~/.mytc 2>&1 | grep VotingPower
# Should print a non-zero number
```

## 11. (Optional but encouraged) Register a relay endpoint

This is what makes MYTC actually *useful* beyond payments — your validator becomes a transport node for MYT Messenger.

You need to run a `myt-relay` server (separate from `mytcd`). Setup is described in the [MYT-Messenger repo](https://github.com/MYTHERRA/MYT-Messenger) (TODO: relay-setup doc). Once your relay is running at e.g. `wss://relay.your-domain.example/ws`, register it on-chain:

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

The MYT Messenger client picks up registered endpoints automatically (within 30 minutes due to the registry cache TTL).

Send a heartbeat occasionally (cron'd hourly is reasonable) so clients know your relay is fresh:

```sh
mytcd tx relay heartbeat \
  --from validator --keyring-backend test --home ~/.mytc \
  --chain-id mytherra-1 --node tcp://127.0.0.1:26657 \
  --fees 100umytc -y
```

> ℹ️ **Phase A status**: the on-chain registry exists and is live, but inter-relay message gossip (so users on different relays can talk to each other) is **not implemented yet** (Phase C). Until that's built, registering an isolated relay endpoint will route some messenger users to a server that can't reach the rest of the federation. For that reason, we currently recommend operating only as a **validator** for the protocol layer; relay-server registration will be encouraged once Phase C is shipped.

## Verify your endpoint

```sh
mytcd query relay endpoint $(mytcd keys show validator --bech val -a --keyring-backend test --home ~/.mytc) \
  --node tcp://127.0.0.1:26657
```

Or via REST (any node):

```sh
curl -s https://explorer.mytherrablockchain.org/mytc-api/mytc/relay/v1/endpoints | jq
```

---

## Operational tips

- **Backups (the only thing that really matters)**:
  - `~/.mytc/config/priv_validator_key.json` — losing this = losing the validator
  - `~/.mytc/config/node_key.json` — losing this = changing your node-id (peers won't recognize you)
  - The mnemonic from step 9 — losing this = losing your funds
  Encrypt and store all three offline.
- **Monitoring**: `journalctl -u mytc.service -f`. Watch for "missed signing" warnings. Persistent missed blocks slowly slash your stake.
- **Don't run the same validator on two machines.** Same `priv_validator_key.json` on two nodes = double-sign slashing = lose your stake.
- **Updates**: when the foundation announces a software upgrade, swap the binary at the announced height. Watch the GitHub releases / governance forum.

## Common pitfalls (lessons learned during the dry-run setup)

1. **Cloud-firewall takes minutes to apply.** UFW rules at the OS level alone are *not enough* if your provider has their own perimeter firewall. Open 26656 in the provider panel; if changes don't take effect within ~5 minutes, power-cycle the server from the panel.
2. **Block sync vs snapshot sync.** Some inter-provider network paths have MTU issues that break Tendermint's block-sync protocol on real-world chains. If you see `module=blockchain` `invalid peer` errors and zero progress, switch to snapshot sync (Section 7b).
3. **External_address is essential.** If you don't set it correctly, other validators can dial *in* to you (because you advertise nothing) and you'll appear unreachable in the peer mesh.
4. **Don't create a new firewall policy when you just want to add a rule.** Most cloud panels have separate "edit existing policy" vs "create new policy" flows; using the wrong one *replaces* the existing rules and locks you out of SSH.
5. **Snapshot-sync hygiene.** When importing another node's data dir, *always* swap back your own `priv_validator_key.json`, `node_key.json`, and reset `priv_validator_state.json` before starting. Two nodes signing with the same validator key get slashed.

## Getting help

- Open an issue at [MYTHERRA/MYTC-chain](https://github.com/MYTHERRA/MYTC-chain/issues)
- Mytherra community Discord (TBD link)
