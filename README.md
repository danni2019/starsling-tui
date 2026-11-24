# starSling TUI

starSling is a Go terminal workspace for futures and options. It displays market snapshots, option chains and Greeks, an implied-volatility and forward panel, unusual option volume, and user-defined spread formulas.

## Try the demo

Install Go 1.25 or newer and run the following from a checkout.

```sh
go run ./cmd/starsling demo
```

The demo renders synthetic prices and option quotes using the live workspace panels. Use Tab to move between panels and q to exit. It runs offline with an in-memory snapshot.

## Live market data

Run `go run ./cmd/starsling` to launch the full workspace. The live path uses Python workers and a user-provided OpenCTP-compatible market-data feed. The application guides you through runtime setup, metadata refresh and feed configuration. Run `go run ./cmd/starsling doctor` for a setup report.

```sh
go test ./...
```

This release derives from the February–April 2026 local starSling baseline. It preserves source commit dates and author names. Repository filtering removed local runtime data and images and rewrote commit IDs. The original BSD 3-Clause license is in [LICENSE](LICENSE).
