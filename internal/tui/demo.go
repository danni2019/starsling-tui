package tui

import (
	"github.com/danni2019/starsling-tui/internal/router"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// RunDemo renders a read-only synthetic snapshot using the live workspace widgets.
func RunDemo() error {
	return newDemo().app.Run()
}

func newDemo() *UI {
	ui := &UI{app: tview.NewApplication(), useArbMonitor: true, overviewSortBy: defaultOverviewSortBy}
	workspace := ui.buildLiveScreen()
	// The sample chain is fixed, so selection never starts a live refresh.
	ui.liveMarket.SetSelectionChangedFunc(nil)

	state := router.NewState()
	state.UpdateMarket(router.MarketSnapshot{
		Seq: 1,
		Rows: []map[string]any{
			{"ctp_contract": "cu2701", "symbol": "cu", "exchange": "SHFE", "last": 74500, "pre_settlement": 74000, "volume": 2140, "open_interest": 8600},
			{"ctp_contract": "cu2702", "symbol": "cu", "exchange": "SHFE", "last": 74680, "pre_settlement": 74300, "volume": 1250, "open_interest": 5400},
		},
	})
	state.UpdateOptions(router.OptionsSnapshot{
		Seq: 1,
		Rows: []map[string]any{
			{"ctp_contract": "cu2701C74500", "underlying": "cu2701", "symbol": "cu", "option_type": "c", "strike": 74500.0, "iv": 0.22, "delta": 0.51, "volume": 48.0},
			{"ctp_contract": "cu2701P74500", "underlying": "cu2701", "symbol": "cu", "option_type": "p", "strike": 74500.0, "iv": 0.24, "delta": -0.49, "volume": 76.0},
		},
	})
	state.UpdateCurve(router.CurveSnapshot{
		Seq: 1,
		Rows: []map[string]any{
			{"ctp_contract": "cu2701", "forward": 74590.0, "volume": 2140.0, "open_interest": 8600.0, "vix": 0.23, "call_skew": 0.02, "put_skew": 0.03},
		},
	})
	snapshot := state.GetViewSnapshot("cu2701")
	ui.applyOverviewSnapshot(snapshot.Overview)
	ui.marketRows = convertMarketRows(snapshot.Market.Rows)
	fillMarketTable(ui.liveMarket, ui.marketRows)
	ui.liveOpts.SetText(renderOptionsPanel(nil, snapshot.Options.Rows, "cu2701", optionRenderFilter{}))
	ui.liveCurve.SetText(renderCurvePanel(snapshot.Curve.Rows))
	ui.applyUnusualSnapshot(router.UnusualSnapshot{
		Seq: 1, Rows: []map[string]any{
			{"ctp_contract": "cu2701P74500", "underlying": "cu2701", "option_type": "p", "strike": 74500.0, "iv": 0.24, "volume": 76.0, "last": 690.0},
		},
	})
	ui.setArbitrageFormula("cu2702 - cu2701")
	ui.renderArbitrageMonitor()

	banner := tview.NewTextView().SetText("starSling demo | synthetic data | Tab change panel · q quit")
	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(banner, 1, 0, false).
		AddItem(workspace, 0, 1, true)
	ui.app.SetRoot(content, true).SetFocus(ui.liveOverview)
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyRune && event.Rune() == 'q':
			ui.app.Stop()
			return nil
		case event.Key() == tcell.KeyTab:
			ui.focusIndex = (ui.focusIndex + 1) % len(ui.focusables)
			ui.app.SetFocus(ui.focusables[ui.focusIndex])
			return nil
		default:
			return event
		}
	})
	return ui
}
