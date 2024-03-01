package minesweeper

import "image"

type MineElementProps uint8

const (
	hiddenBits MineElementProps = 0b1
	markedBits MineElementProps = 0b10
)

type MineElement struct {
	Pos   image.Point
	Props MineElementProps
	Value int8
}

func (m *MineElement) isActive(prop MineElementProps) bool {
	return m.Props&prop != 0
}

func (m *MineElement) IsMarked() bool {
	return m.isActive(markedBits)
}

func (m *MineElement) IsHidden() bool {
	return m.isActive(hiddenBits)
}

func (m *MineElement) ToggleProp(prop MineElementProps) {
	m.Props ^= prop
}

func (m *MineElement) PropOn(prop MineElementProps) {
	if !m.isActive(prop) {
		m.ToggleProp(prop)
	}
}

func (m *MineElement) PropOff(prop MineElementProps) {
	if m.isActive(prop) {
		m.ToggleProp(prop)
	}
}
