package engine

import (
	"image"
	"time"

	"gioui.org/io/pointer"
)

type MinesweeperEngine interface {
	Resize(width uint16, height uint16, mines uint16)
	Restart()
	Close()
	OnButtonClick(pos image.Point, clickType pointer.Buttons) MinesweeperState
	GetRemainingMines() *[]*MineElement
	SetChannels(mainChannel chan MineElement, acks chan uint8) MinesweeperEngine
	SetAnimationDuration(animationDuration time.Duration) MinesweeperEngine
	GetWidth() int
	GetHeight() int
	GetState() MinesweeperState
	GetRevealed() uint16
	GetMarked() uint16
	GetMines() uint16
}

type MinesweeperState uint8

const (
	UNDEFINED MinesweeperState = iota
	START
	RUNNING
	LOSE
	WIN
	END
	LOADING
)

func (s MinesweeperState) ToString() string {
	switch s {
	case UNDEFINED:
		return "UNDEFINED"
	case START:
		return "START"
	case RUNNING:
		return "RUNNING"
	case LOSE:
		return "LOSE"
	case WIN:
		return "WIN"
	case END:
		return "END"
	case LOADING:
		return "LOADING"
	}

	return "UNKNOWN"
}
