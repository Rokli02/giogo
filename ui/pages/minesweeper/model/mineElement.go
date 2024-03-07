package model

import "image"

type MineElementProps uint8

const (
	HiddenBits MineElementProps = 0b1
	MarkedBits MineElementProps = 0b10
)

type MineElement struct {
	Pos   image.Point
	Props MineElementProps
	Value int8
}

func (m *MineElement) isPropActive(prop MineElementProps) bool {
	return m.Props&prop != 0
}

func (m *MineElement) IsMarked() bool {
	return m.isPropActive(MarkedBits)
}

func (m *MineElement) IsHidden() bool {
	return m.isPropActive(HiddenBits)
}

func (m *MineElement) ToggleProp(prop MineElementProps) {
	m.Props ^= prop
}

func (m *MineElement) PropOn(prop MineElementProps) {
	if !m.isPropActive(prop) {
		m.ToggleProp(prop)
	}
}

func (m *MineElement) PropOff(prop MineElementProps) {
	if m.isPropActive(prop) {
		m.ToggleProp(prop)
	}
}
