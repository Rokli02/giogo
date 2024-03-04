package main

import (
	"fmt"
	"giogo/assets"
	"giogo/ui"
	"giogo/ui/pages/menu"
	"giogo/ui/pages/minesweeper"
	"giogo/ui/pages/minesweeper/engine"
	"giogo/ui/pages/scroller"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"image"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		fmt.Println("App Starting")

		w := app.NewWindow(func(_ unit.Metric, c *app.Config) {
			c.Title = "Freestyle"
		})

		if err := run(w); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	styles.MaterialTheme = material.NewTheme()
	var ops op.Ops
	var gtx layout.Context

	assets.InitializeAssets()
	router := routerModule.NewRouter[ui.ApplicationCycles, string](w)
	addRoutes(router, w)
	router.Select(routerModule.MenuPage)

	currentPage := router.CurrentRoute()
	currentPage.Initialize()

	previousRouteKey := router.CurrentKey()

	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx = layout.NewContext(&ops, e)
			gtx.Constraints.Min = image.Point{}

			handleRouting(&gtx, router)

			if previousRouteKey != router.CurrentKey() {
				currentPage.Close()
				currentPage = router.CurrentRoute()
				currentPage.Initialize()
			}

			previousRouteKey = router.CurrentKey()

			currentPage.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

func addRoutes(router *routerModule.Router[ui.ApplicationCycles, string], w *app.Window) {
	singlePlayerMinesweeperEngine := engine.NewMinesweeperLocalEngine()
	singlePlayerMinesweeperEngine.Resize(8, 12, 12)

	router.Add(routerModule.MenuPage, menu.NewMenu(w, router))
	router.Add(routerModule.MinesweeperMenuPage, menu.NewMinesweeperMenu(w, router))
	router.Add(routerModule.MinesweeperMultiplayerMenuPage, menu.NewMinesweeperMultiplayerMenu(w, router))
	router.Add(routerModule.MinesweeperPage, minesweeper.NewMinefield(w, router, time.Millisecond*40).SetEngine(singlePlayerMinesweeperEngine))
	router.Add(routerModule.ScrollerPage, scroller.NewScrollPage(w, 181))
}

func handleRouting(gtx *layout.Context, router *routerModule.Router[ui.ApplicationCycles, string]) {
	for _, evt := range gtx.Events(0) {
		switch event := evt.(type) {
		case key.Event:
			router.Select(event.Name)
		}
	}

	key.InputOp{
		Keys: key.Set("0"),
		Tag:  0,
	}.Add(gtx.Ops)
}
