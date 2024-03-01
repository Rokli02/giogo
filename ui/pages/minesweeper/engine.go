package minesweeper

import (
	"fmt"
	"image"
	"math/rand"
	"time"

	"gioui.org/io/pointer"
)

type MinesweeperEngine struct {
	Width    uint16
	Height   uint16
	MaxMines uint16
	State    MinesweeperState
	Revealed uint16
	Marked   uint16

	mineChannel       chan *MineElement
	acks              chan uint8
	mines             uint16
	matrix            [][]*MineElement
	animationDuration time.Duration
}

func NewMinesweeperEngine() *MinesweeperEngine {
	me := &MinesweeperEngine{
		State:             UNDEFINED,
		animationDuration: time.Millisecond * 40,
	}

	return me
}

func (me *MinesweeperEngine) Resize(width uint16, height uint16, mines uint16) {
	if width != me.Width || height != me.Height {
		me.Width = width
		me.Height = height

		clear(me.matrix)
		me.matrix = make([][]*MineElement, height)

		for rowIndex := range me.matrix {
			me.matrix[rowIndex] = make([]*MineElement, width)

			for colIndex := range me.matrix[rowIndex] {
				element := &me.matrix[rowIndex][colIndex]
				*element = &MineElement{Value: 0, Props: hiddenBits, Pos: image.Point{colIndex, rowIndex}}
			}
		}
	} else {
		for rowIndex := range me.matrix {
			for colIndex := range me.matrix[rowIndex] {
				me.matrix[rowIndex][colIndex] = &MineElement{Value: 0, Props: hiddenBits, Pos: image.Point{colIndex, rowIndex}}
			}
		}
	}

	me.MaxMines = mines
	me.State = START
	me.Revealed = 0
	me.Marked = 0

	me.mines = mines
}

func (me *MinesweeperEngine) generateMines(clickPos image.Point) [][]*MineElement {
	minePositions := make([]image.Point, 0, me.MaxMines)
	// Calculate them

	for i, calcTries := 0, 0; i < cap(minePositions) || calcTries > 5; calcTries++ {
		validPos := true

		// Random num
		minePos := image.Point{int(rand.Int31n(int32(me.Width))), int(rand.Int31n(int32(me.Height)))}

		// Check if is it stored already or at 'clickPos'
		if minePos == clickPos {
			continue
		}

		for i := range minePositions {
			if minePositions[i] == minePos {
				validPos = false
				break
			}
		}

		if validPos {
			minePositions = append(minePositions, minePos)
			i++
			calcTries = 0
		}
	}

	// Clear minefield
	for rowIndex := range me.matrix {
		for colIndex := range me.matrix[rowIndex] {
			me.matrix[rowIndex][colIndex].Value = 0
		}
	}

	// Plant mines
	for i := range minePositions {
		me.matrix[minePositions[i].Y][minePositions[i].X].Value = -1
	}
	me.mines = uint16(len(minePositions))

	// Find mines in the neighborhood, MAAAN!
	for rowIndex := range me.matrix {
		for colIndex := range me.matrix[rowIndex] {
			element := me.matrix[rowIndex][colIndex]

			if element.Value == -1 {
				continue
			}

			element.Value = me.neighboringMines(rowIndex, colIndex)
		}
	}

	return me.matrix
}

func (me *MinesweeperEngine) neighboringMines(rowIndexParam, colIndexParam int) int8 {
	var sum int8 = 0
	// Row loop
	for i := -1; i <= 1; i++ {
		rowIndex := rowIndexParam + i
		if rowIndex < 0 || rowIndex > int(me.Height-1) {
			continue
		}

		// Column loop
		for j := -1; j <= 1; j++ {
			if j == 0 && i == 0 {
				continue
			}

			colIndex := colIndexParam + j
			if colIndex < 0 || colIndex > int(me.Width-1) {
				continue
			}

			if me.matrix[rowIndex][colIndex].Value == -1 {
				sum++
			}
		}
	}

	return sum
}

func (me *MinesweeperEngine) OnButtonClick(pos image.Point, clickType pointer.Buttons) (MinesweeperState, *MineElement) {
	element := me.matrix[pos.Y][pos.X]
	var returnState MinesweeperState = RUNNING
	var returnElement *MineElement = element

	switch me.State {
	case START:
		fmt.Println("--- Game Start ---")
		me.State = RUNNING

		// Legenerálni a bombákat
		me.generateMines(pos)

		fallthrough
	case RUNNING:
		if element.IsHidden() && clickType == pointer.ButtonSecondary {
			element.ToggleProp(markedBits)

			if element.IsMarked() {
				me.Marked++
			} else {
				me.Marked--
			}

			break
		}

		me.Revealed++

		switch element.Value {
		case -1:
			me.State = LOSE

			returnState = LOSE
		case 0:
			element.PropOff(hiddenBits)

			go me.floodNeutralCells(pos)
			returnState = LOADING
		default:
			element.PropOff(hiddenBits)

			returnState = RUNNING
			returnElement = element
		}

		if me.State == RUNNING && me.Revealed >= me.Width*me.Height-me.mines {
			fmt.Println("--- GG, WIN ---")
			me.State = WIN
			me.Marked = me.mines

			returnState = WIN
			returnElement = element
		}

	}

	return returnState, returnElement
}

func (me *MinesweeperEngine) floodNeutralCells(startingPoint image.Point) {
	me.State = LOADING
	floodedPos := make([]image.Point, 0, 8)
	floodedPos = append(floodedPos, startingPoint)

	me.mineChannel <- me.matrix[startingPoint.Y][startingPoint.X]
	time.Sleep(me.animationDuration)
	<-me.acks

	countOfFloodedCells := uint16(0)

	for iterator := 0; iterator < len(floodedPos); iterator++ {
		// Venni a jelenlegi elem rejtett környezetét és azokat hozzáadni egy listához
		for i := -1; i <= 1; i++ {
			rowIndex := floodedPos[iterator].Y + i

			// Kilóg felül, vagy alul
			if rowIndex < 0 || rowIndex > int(me.Height-1) {
				continue
			}

			// Column loop
			for j := -1; j <= 1; j++ {
				if me.State != LOADING {
					return
				}

				if j == 0 && i == 0 {
					continue
				}

				colIndex := floodedPos[iterator].X + j

				// Kilóg bal, vagy jobb oldalt
				if colIndex < 0 || colIndex > int(me.Width-1) {
					continue
				}

				// Az adott elem értékét megvizsgálni
				element := me.matrix[rowIndex][colIndex]
				if !element.IsHidden() || element.IsMarked() {
					continue
				}

				// felfedni és növelni a felfedettek számát
				element.PropOff(hiddenBits)
				countOfFloodedCells++

				// Ha 0, akkor feldeni, listához hozzáadni
				if element.Value == 0 {
					floodedPos = append(floodedPos, element.Pos)
				}

				// Lassú animáció
				if me.mineChannel == nil {
					return
				}

				me.mineChannel <- element
				time.Sleep(me.animationDuration)
				<-me.acks
			}
		}
	}

	me.State = RUNNING
	me.Revealed += countOfFloodedCells
}
