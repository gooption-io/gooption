init tendermint

```bash
tendermint init --home ~/.gooptiond
```

init gooptiond

```bash
gooptiond --home ~/.gooptiond init
```

start tendermint

```bash
tendermint node --home ~/.gooptiond  --p2p.skip_upnp
```

start gooptiond

```bash
gooptiond --home ~/.gooptiond start --bind tcp://127.0.0.1:26658
```
