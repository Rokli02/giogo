package engine

import (
	"fmt"
	"image"
	"time"

	"giogo/ui/pages/minesweeper/engine/logic"
	"giogo/ui/pages/minesweeper/model"

	"gioui.org/io/pointer"
)

type MinesweeperLocalEngine struct {
	width    uint16
	height   uint16
	maxMines uint16
	state    model.MinesweeperState
	revealed uint16
	marked   uint16

	mineChannel       chan model.MineElement
	acks              chan uint8
	mines             uint16
	matrix            [][]*model.MineElement
	animationDuration time.Duration
}

// Check if 'MinesweeperLocalEngine' implements every methods of interface 'MinesweeperEngine'
var _ MinesweeperEngine = (*MinesweeperLocalEngine)(nil)

func NewMinesweeperLocalEngine() *MinesweeperLocalEngine {
	me := &MinesweeperLocalEngine{
		state: model.UNDEFINED,
	}

	return me
}

func (me *MinesweeperLocalEngine) Resize(width uint16, height uint16, mines uint16) {
	if width != me.width || height != me.height {
		me.width = width
		me.height = height

		clear(me.matrix)
		me.matrix = make([][]*model.MineElement, height)

		for rowIndex := range me.matrix {
			me.matrix[rowIndex] = make([]*model.MineElement, width)

			for colIndex := range me.matrix[rowIndex] {
				me.matrix[rowIndex][colIndex] = &model.MineElement{Value: 0, Props: model.HiddenBits, Pos: image.Point{colIndex, rowIndex}}
			}
		}

		me.maxMines = mines
		me.state = model.START
		me.revealed = 0
		me.marked = 0

		me.mines = mines
	} else {
		me.Restart()
	}
}

func (me *MinesweeperLocalEngine) Restart() {
	for rowIndex := range me.matrix {
		for colIndex := range me.matrix[rowIndex] {
			me.matrix[rowIndex][colIndex] = &model.MineElement{Value: 0, Props: model.HiddenBits, Pos: image.Point{colIndex, rowIndex}}
		}
	}

	me.state = model.START
	me.revealed = 0
	me.marked = 0
}

func (me *MinesweeperLocalEngine) OnButtonClick(pos image.Point, clickType pointer.Buttons) model.MinesweeperState {
	element := me.matrix[pos.Y][pos.X]
	var returnState model.MinesweeperState = model.RUNNING
	var returnElement *model.MineElement = element

	switch me.state {
	case model.START:
		fmt.Println("--- Game Start ---")
		me.state = model.RUNNING
		returnState = model.RUNNING

		// Legenerálni a bombákat
		me.mines = logic.GenerateMines(pos, me.matrix, me.maxMines)

		fallthrough
	case model.RUNNING:
		if element.IsHidden() && clickType == pointer.ButtonSecondary {
			element.ToggleProp(model.MarkedBits)

			if element.IsMarked() {
				me.marked++
			} else {
				me.marked--
			}

			break
		}

		me.revealed++

		switch element.Value {
		case -1:
			me.state = model.LOSE

			returnState = model.LOSE
		case 0:
			element.PropOff(model.HiddenBits)

			go me.floodNeutralCells(pos)
			returnState = model.LOADING
		default:
			element.PropOff(model.HiddenBits)

			returnState = model.RUNNING
			returnElement = element
		}

		if me.state == model.RUNNING && me.revealed >= me.width*me.height-me.mines {
			fmt.Println("--- GG, WIN ---")
			me.state = model.WIN
			me.marked = me.mines

			me.mineChannel <- *returnElement
			<-me.acks

			return model.WIN
		}
	}

	switch returnState {
	case model.RUNNING:
		me.mineChannel <- *returnElement
		<-me.acks
	}

	return returnState
}

func (me *MinesweeperLocalEngine) Close() {
	me.mineChannel = nil
}

func (me *MinesweeperLocalEngine) SetChannels(mainChannel chan model.MineElement, acks chan uint8) MinesweeperEngine {
	me.mineChannel = mainChannel
	me.acks = acks

	return me
}

func (me *MinesweeperLocalEngine) SetAnimationDuration(animationDuration time.Duration) MinesweeperEngine {
	me.animationDuration = animationDuration

	return me
}

func (me *MinesweeperLocalEngine) GetWidth() int {
	return int(me.width)
}

func (me *MinesweeperLocalEngine) GetHeight() int {
	return int(me.height)
}

func (me *MinesweeperLocalEngine) GetState() model.MinesweeperState {
	return me.state
}

func (me *MinesweeperLocalEngine) GetRevealed() uint16 {
	return me.revealed
}

func (me *MinesweeperLocalEngine) GetMarked() uint16 {
	return me.marked
}

func (me *MinesweeperLocalEngine) GetMines() uint16 {
	return me.mines
}

func (me *MinesweeperLocalEngine) GetRemainingMines() *[]*model.MineElement {
	matrix := make([]*model.MineElement, 0, (me.height*me.width)>>2)

	for rowIndex := range me.matrix {
		for colIndex := range me.matrix[rowIndex] {
			if me.matrix[rowIndex][colIndex].IsHidden() {
				me.matrix[rowIndex][colIndex].PropOff(model.HiddenBits)

				matrix = append(matrix, me.matrix[rowIndex][colIndex])
			}
		}
	}

	return &matrix
}

func (me *MinesweeperLocalEngine) floodNeutralCells(startingPoint image.Point) {
	me.state = model.LOADING
	floodedPos := make([]image.Point, 0, 8)
	floodedPos = append(floodedPos, startingPoint)

	if me.mineChannel == nil {
		return
	}

	me.mineChannel <- *me.matrix[startingPoint.Y][startingPoint.X]
	time.Sleep(me.animationDuration)
	<-me.acks

	countOfFloodedCells := uint16(0)

	for iterator := 0; iterator < len(floodedPos); iterator++ {
		// Venni a jelenlegi elem rejtett környezetét és azokat hozzáadni egy listához
		for i := -1; i <= 1; i++ {
			rowIndex := floodedPos[iterator].Y + i

			// Kilóg felül, vagy alul
			if rowIndex < 0 || rowIndex > int(me.height-1) {
				continue
			}

			// Column loop
			for j := -1; j <= 1; j++ {
				if me.state != model.LOADING {
					return
				}

				if j == 0 && i == 0 {
					continue
				}

				colIndex := floodedPos[iterator].X + j

				// Kilóg bal, vagy jobb oldalt
				if colIndex < 0 || colIndex > int(me.width-1) {
					continue
				}

				// Az adott elem értékét megvizsgálni
				element := me.matrix[rowIndex][colIndex]
				if !element.IsHidden() || element.IsMarked() {
					continue
				}

				// felfedni és növelni a felfedettek számát
				element.PropOff(model.HiddenBits)
				countOfFloodedCells++

				// Ha 0, akkor feldeni, listához hozzáadni
				if element.Value == 0 {
					floodedPos = append(floodedPos, element.Pos)
				}

				// Lassú animáció
				if me.mineChannel == nil {
					return
				}

				me.mineChannel <- *element
				time.Sleep(me.animationDuration)
				<-me.acks
			}
		}
	}

	me.state = model.RUNNING
	me.revealed += countOfFloodedCells
}
