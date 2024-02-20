package main

import (
	"fmt"
	"giogo/assets"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"image"
	"log"
	"os"

	"gioui.org/app"
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
	router := routerModule.NewRouter(w)

	currentPage := router.CurrentPage()
	currentPage.Initialize()

	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx = layout.NewContext(&ops, e)
			gtx.Constraints.Min = image.Point{}

			pageChanged := router.HandleRouteEvents(&gtx)

			if pageChanged {
				currentPage.Close()
				currentPage = router.CurrentPage()
				currentPage.Initialize()
			}

			currentPage.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
