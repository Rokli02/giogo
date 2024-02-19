package router

import (
	"errors"
	"fmt"
	"giogo/ui"
	"giogo/ui/pages/minesweeper"
	"giogo/ui/pages/scroller"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
)

type Router struct {
	currentRoute Route
	currentPage  ui.ApplicationCycles
	routes       []ui.ApplicationCycles
}

func NewRouter(w *app.Window) (r *Router) {
	r = &Router{
		currentRoute: Minesweeper,
		routes:       make([]ui.ApplicationCycles, 0, 2),
	}

	r.appendRoutes(w)

	r.currentPage = r.routes[r.currentRoute]

	return r
}

type Route int

const (
	Minesweeper Route = iota
	Scroller
	CountOfRoutes
)

func (r *Router) appendRoutes(w *app.Window) {
	r.routes = append(r.routes, minesweeper.NewMinefield(w, 15, 20, 48))
	r.routes = append(r.routes, scroller.NewScrollPage(w, 175))
}

func (r *Router) SelectPage(route Route) error {
	if int(route) >= len(r.routes) || route < 0 {
		return errors.New("no route found")
	}

	page := r.routes[route]
	r.currentRoute = route
	r.currentPage = page

	return nil
}

func (r *Router) CurrentPage() ui.ApplicationCycles {
	return r.currentPage
}

func (r *Router) GetPage(route Route) (ui.ApplicationCycles, error) {
	if int(route) >= len(r.routes) || route < 0 {
		return nil, errors.New("no route found")
	}

	return r.routes[route], nil
}

func (r *Router) CurrentRoute() Route {
	return r.currentRoute
}

func (r *Router) HandleRouteEvents(gtx *layout.Context) bool {
	// currentRoute, nextRoute := r.currentRoute, r.currentRoute
	isChanged := false

	for _, evt := range gtx.Events(0) {
		switch event := evt.(type) {
		case key.Event:
			previousRoute := r.currentRoute
			// currentRoute = r.currentRoute

			switch event.Name {
			case "0":
				r.SelectPage(Minesweeper)
			case "1":
				r.SelectPage(Scroller)
			default:
				fmt.Println("Invalid Page selection!", event.Modifiers, event.Name)
			}

			if previousRoute != r.currentRoute {
				isChanged = true
			}
			// nextRoute = r.currentRoute
		}
	}

	key.InputOp{
		Keys: key.Set("0|1"),
		Tag:  0,
	}.Add(gtx.Ops)

	return isChanged
	// return currentRoute != nextRoute
}
