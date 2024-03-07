package engine

import (
	"image"
	"time"

	"giogo/ui/pages/minesweeper/model"

	"gioui.org/io/pointer"
)

type MinesweeperEngine interface {
	Resize(width uint16, height uint16, mines uint16)
	Restart()
	Close()
	OnButtonClick(pos image.Point, clickType pointer.Buttons) model.MinesweeperState
	GetRemainingMines() *[]*model.MineElement
	SetChannels(mainChannel chan model.MineElement, acks chan uint8) MinesweeperEngine
	SetAnimationDuration(animationDuration time.Duration) MinesweeperEngine
	GetWidth() int
	GetHeight() int
	GetState() model.MinesweeperState
	GetRevealed() uint16
	GetMarked() uint16
	GetMines() uint16
}
