package router

import (
	"errors"
	"fmt"

	"gioui.org/app"
)

type Router[Route interface{}, Key int | string] struct {
	w            *app.Window
	routes       map[Key]Route
	currentRoute *Route
	currentKey   Key
}

func NewRouter[Route interface{}, Key int | string](w *app.Window) (r *Router[Route, Key]) {
	r = &Router[Route, Key]{
		w:            w,
		routes:       make(map[Key]Route),
		currentRoute: nil,
	}

	return
}

func (r *Router[Route, Key]) Reset() {
	clear(r.routes)
}

func (r *Router[Route, Key]) Add(key Key, route Route) {
	_, hasItem := r.routes[key]
	if !hasItem {
		r.routes[key] = route

		return
	}

	fmt.Printf("Route '%v' has already been added!", key)
}

func (r *Router[Route, Key]) Select(key Key) (e error) {
	currentRoute, hasItem := r.routes[key]

	if !hasItem {
		e = errors.New("no route found")

		return
	}

	r.currentRoute = &currentRoute
	r.currentKey = key

	return
}

func (r *Router[Route, Key]) GoTo(key Key) (e error) {
	e = r.Select(key)
	r.Rerender()

	return
}

func (r *Router[Route, Key]) Remove(key Key) {
	delete(r.routes, key)
}

func (r *Router[Route, Key]) CurrentKey() Key {
	return r.currentKey
}

func (r *Router[Route, Key]) CurrentRoute() Route {
	return *r.currentRoute
}

func (r *Router[Route, Key]) Rerender() {
	r.w.Invalidate()
}
