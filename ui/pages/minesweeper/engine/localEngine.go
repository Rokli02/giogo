package engine

import (
	"fmt"
	"image"
	"math/rand"
	"time"

	"gioui.org/io/pointer"
)

// TODO: Létrehozni egy olyan motort, ami képes a multiplayer kezelésésre.
// TODO: Létrehozni egy szerver komponenst, ami pár gomb megnyomására létrehoz, és bezár egyet

type MinesweeperLocalEngine struct {
	width    uint16
	height   uint16
	maxMines uint16
	state    MinesweeperState
	revealed uint16
	marked   uint16

	mineChannel       chan MineElement
	acks              chan uint8
	mines             uint16
	matrix            [][]*MineElement
	animationDuration time.Duration
}

// Check if 'MinesweeperLocalEngine' implements every methods of interface 'MinesweeperEngine'
var _ MinesweeperEngine = (*MinesweeperLocalEngine)(nil)

func NewMinesweeperLocalEngine() *MinesweeperLocalEngine {
	me := &MinesweeperLocalEngine{
		state: UNDEFINED,
	}

	return me
}

func (me *MinesweeperLocalEngine) Resize(width uint16, height uint16, mines uint16) {
	if width != me.width || height != me.height {
		me.width = width
		me.height = height

		clear(me.matrix)
		me.matrix = make([][]*MineElement, height)

		for rowIndex := range me.matrix {
			me.matrix[rowIndex] = make([]*MineElement, width)

			for colIndex := range me.matrix[rowIndex] {
				element := &me.matrix[rowIndex][colIndex]
				*element = &MineElement{Value: 0, Props: HiddenBits, Pos: image.Point{colIndex, rowIndex}}
			}
		}

		me.maxMines = mines
		me.state = START
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
			me.matrix[rowIndex][colIndex] = &MineElement{Value: 0, Props: HiddenBits, Pos: image.Point{colIndex, rowIndex}}
		}
	}

	me.state = START
	me.revealed = 0
	me.marked = 0
}

func (me *MinesweeperLocalEngine) OnButtonClick(pos image.Point, clickType pointer.Buttons) MinesweeperState {
	element := me.matrix[pos.Y][pos.X]
	var returnState MinesweeperState = RUNNING
	var returnElement *MineElement = element

	switch me.state {
	case START:
		fmt.Println("--- Game Start ---")
		me.state = RUNNING
		returnState = RUNNING

		// Legenerálni a bombákat
		me.generateMines(pos)

		fallthrough
	case RUNNING:
		if element.IsHidden() && clickType == pointer.ButtonSecondary {
			element.ToggleProp(MarkedBits)

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
			me.state = LOSE

			returnState = LOSE
		case 0:
			element.PropOff(HiddenBits)

			go me.floodNeutralCells(pos)
			returnState = LOADING
		default:
			element.PropOff(HiddenBits)

			returnState = RUNNING
			returnElement = element
		}

		if me.state == RUNNING && me.revealed >= me.width*me.height-me.mines {
			fmt.Println("--- GG, WIN ---")
			me.state = WIN
			me.marked = me.mines

			me.mineChannel <- *returnElement
			<-me.acks

			return WIN
		}
	}

	switch returnState {
	case RUNNING:
		me.mineChannel <- *returnElement
		<-me.acks
	}

	return returnState
}

func (me *MinesweeperLocalEngine) Close() {
	me.mineChannel = nil
}

func (me *MinesweeperLocalEngine) SetChannels(mainChannel chan MineElement, acks chan uint8) MinesweeperEngine {
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

func (me *MinesweeperLocalEngine) GetState() MinesweeperState {
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

func (me *MinesweeperLocalEngine) GetRemainingMines() *[]*MineElement {
	matrix := make([]*MineElement, 0, (me.height*me.width)>>2)

	for rowIndex := range me.matrix {
		for colIndex := range me.matrix[rowIndex] {
			if me.matrix[rowIndex][colIndex].IsHidden() {
				me.matrix[rowIndex][colIndex].PropOff(HiddenBits)

				matrix = append(matrix, me.matrix[rowIndex][colIndex])
			}
		}
	}

	return &matrix
}

func (me *MinesweeperLocalEngine) generateMines(clickPos image.Point) [][]*MineElement {
	minePositions := make([]image.Point, 0, me.maxMines)
	// Calculate them

	for i, calcTries := 0, 0; i < cap(minePositions) || calcTries > 5; calcTries++ {
		validPos := true

		// Random num
		minePos := image.Point{int(rand.Int31n(int32(me.width))), int(rand.Int31n(int32(me.height)))}

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

func (me *MinesweeperLocalEngine) neighboringMines(rowIndexParam, colIndexParam int) int8 {
	var sum int8 = 0
	// Row loop
	for i := -1; i <= 1; i++ {
		rowIndex := rowIndexParam + i
		if rowIndex < 0 || rowIndex > int(me.height-1) {
			continue
		}

		// Column loop
		for j := -1; j <= 1; j++ {
			if j == 0 && i == 0 {
				continue
			}

			colIndex := colIndexParam + j
			if colIndex < 0 || colIndex > int(me.width-1) {
				continue
			}

			if me.matrix[rowIndex][colIndex].Value == -1 {
				sum++
			}
		}
	}

	return sum
}

func (me *MinesweeperLocalEngine) floodNeutralCells(startingPoint image.Point) {
	me.state = LOADING
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
				if me.state != LOADING {
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
				element.PropOff(HiddenBits)
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

	me.state = RUNNING
	me.revealed += countOfFloodedCells
}
