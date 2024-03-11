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
	previousKeys []Key
	currentKey   Key
}

func NewRouter[Route interface{}, Key int | string](w *app.Window) (r *Router[Route, Key]) {
	r = &Router[Route, Key]{
		w:            w,
		routes:       make(map[Key]Route),
		previousKeys: make([]Key, 0, 1),
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
	previousKey := r.currentKey

	if err := r.Select(key); err != nil {
		e = err

		return
	}

	r.previousKeys = append(r.previousKeys, previousKey)

	r.Rerender()

	return
}

func (r *Router[Route, Key]) GoBack() (e error) {

	if len(r.previousKeys) == 0 {
		return
	}

	lastIndex := len(r.previousKeys) - 1
	lastKey := r.previousKeys[lastIndex]
	r.previousKeys = r.previousKeys[:lastIndex]

	if err := r.Select(lastKey); err != nil {
		e = err

		return
	}
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

func (r *Router[Route, Key]) WipeHistory() {
	r.previousKeys = make([]Key, 0, 1)
}

func (r *Router[Route, Key]) WipeHistoryLastN(n int) {
	r.previousKeys = r.previousKeys[:len(r.previousKeys)-n]
}

func (r *Router[Route, Key]) WipeHistoryUntilKey(key Key) {
	for i := len(r.previousKeys) - 1; i >= 0; i-- {
		if r.previousKeys[i] == key {
			r.previousKeys = r.previousKeys[:i]

			return
		}
	}
}
