package server

import (
	"fmt"
	"giogo/ui/pages/minesweeper/engine/logic"
	"giogo/ui/pages/minesweeper/model"
	"giogo/utils"
	"image"
	"time"

	"gioui.org/io/pointer"
)

type MinesweeperServerEngine struct {
	width    uint16
	height   uint16
	maxMines uint16
	state    model.MinesweeperState
	revealed uint16
	marked   uint16

	mines             uint16
	matrix            [][]*model.MineElement
	animationDuration time.Duration

	broadcastToClient func(data SocketData)
}

type _minesweeperEngine interface {
	Resize(width uint16, height uint16, mines uint16, isHost bool)
	Restart()
	OnPositionAction(pos image.Point, clickType pointer.Buttons)
	GetRemainingMines() []*model.MineElement
	SetAnimationDuration(animationDuration time.Duration) _minesweeperEngine
}

// Interface implementation check
var _ _minesweeperEngine = (*MinesweeperServerEngine)(nil)

func NewMinesweeperServerEngine(broadcastToClient func(data SocketData)) *MinesweeperServerEngine {
	m := &MinesweeperServerEngine{
		state:             model.WAITING,
		broadcastToClient: broadcastToClient,
	}

	return m
}

func (m *MinesweeperServerEngine) Resize(width uint16, height uint16, mines uint16, isHost bool) {
	fmt.Printf("Server resize (host=%t)\n", isHost)

	if isHost {
		m.width = width
		m.height = height

		clear(m.matrix)
		m.matrix = make([][]*model.MineElement, height)

		for rowIndex := range m.matrix {
			m.matrix[rowIndex] = make([]*model.MineElement, width)

			for colIndex := range m.matrix[rowIndex] {
				m.matrix[rowIndex][colIndex] = &model.MineElement{Value: 0, Props: model.HiddenBits, Pos: image.Point{colIndex, rowIndex}}
			}
		}

		m.mines = mines
		m.maxMines = mines
		m.state = model.START
		m.revealed = 0
		m.marked = 0
	}

	socketData := SocketData{DataType: RESIZE}
	socketData.Data = make([]byte, 0, 9)
	socketData.Data = append(socketData.Data, byte(model.START))
	socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(m.marked)...)
	socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(width)...)
	socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(height)...)
	socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(mines)...)

	m.broadcastToClient(socketData)
}

func (m *MinesweeperServerEngine) Restart() {
	for rowIndex := range m.matrix {
		for colIndex := range m.matrix[rowIndex] {
			m.matrix[rowIndex][colIndex] = &model.MineElement{Value: 0, Props: model.HiddenBits, Pos: image.Point{colIndex, rowIndex}}
		}
	}

	m.state = model.START
	m.revealed = 0
	m.marked = 0

	socketData := SocketData{DataType: RESTART}
	socketData.Data = make([]byte, 0, 3)
	socketData.Data = append(socketData.Data, byte(model.START))
	socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(m.marked)...)

	m.broadcastToClient(socketData)
}

func (m *MinesweeperServerEngine) OnPositionAction(pos image.Point, clickType pointer.Buttons) {
	fmt.Println("\n - Pos:", pos)
	element := m.matrix[pos.Y][pos.X]

	switch m.state {
	case model.START:
		m.state = model.RUNNING

		m.mines = logic.GenerateMines(pos, m.matrix, m.maxMines)

		socketData := SocketData{DataType: STATE}
		socketData.Data = make([]byte, 0, 7)
		socketData.Data = append(socketData.Data, byte(model.START))
		socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(m.width)...)
		socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(m.height)...)
		socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(m.maxMines)...)

		m.broadcastToClient(socketData)

		fallthrough
	case model.RUNNING:
		if element.IsHidden() && clickType == pointer.ButtonSecondary {
			element.ToggleProp(model.MarkedBits)

			if element.IsMarked() {
				m.marked++
			} else {
				m.marked--
			}

			socketData := SocketData{DataType: POSITION}
			socketData.Data = make([]byte, 0, 13)
			socketData.Data = append(socketData.Data, byte(model.RUNNING))
			socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(m.marked)...)
			socketData.Data = append(socketData.Data, utils.ByteConverter.MineElementToBytes(*element)...)

			m.broadcastToClient(socketData)

			break
		}

		m.revealed++

		switch element.Value {
		case -1:
			m.state = model.LOSE

			socketData := SocketData{DataType: END_OF_GAME}
			remainingMines := m.GetRemainingMines()

			socketData.Data = make([]byte, 0, 3+len(remainingMines)*model.SizeOfMineElementInBytes)
			socketData.Data = append(socketData.Data, byte(model.LOSE))
			socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(m.marked)...)

			for _, mine := range remainingMines {
				socketData.Data = append(socketData.Data, utils.ByteConverter.MineElementToBytes(*mine)...)
			}

			m.broadcastToClient(socketData)

			return
		case 0:
			m.state = model.LOADING

			revealedCells, countOfFloodedCells := logic.RevealedCells(pos, m.matrix)
			m.revealed += countOfFloodedCells

			go func() {
				for _, revealedCell := range revealedCells {
					element := m.matrix[revealedCell.Y][revealedCell.X]

					socketData := SocketData{DataType: STATE}
					socketData.Data = make([]byte, 0, 13)
					socketData.Data = append(socketData.Data, byte(model.LOADING))
					socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(m.marked)...)
					socketData.Data = append(socketData.Data, utils.ByteConverter.MineElementToBytes(*element)...)

					m.broadcastToClient(socketData)
					time.Sleep(m.animationDuration)
				}

				m.state = model.RUNNING
				socketData := SocketData{DataType: STATE}
				socketData.Data = make([]byte, 0, 3) // 13
				socketData.Data = append(socketData.Data, byte(model.RUNNING))
				socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(m.marked)...)

				m.broadcastToClient(socketData)
			}()
		default:
			element.PropOff(model.HiddenBits)
			m.state = model.RUNNING

			socketData := SocketData{DataType: POSITION}
			socketData.Data = make([]byte, 0, 13)
			socketData.Data = append(socketData.Data, byte(m.state))
			socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(m.marked)...)
			socketData.Data = append(socketData.Data, utils.ByteConverter.MineElementToBytes(*element)...)

			m.broadcastToClient(socketData)
		}

		fmt.Println("After click states:", m.state.ToString(), "revealed:", m.revealed, "width|height|mines:", m.width, m.height, m.mines)
		if (m.state == model.RUNNING || m.state == model.LOADING) && m.revealed >= m.width*m.height-m.mines {
			fmt.Println("--- GG, WIN ---")
			m.state = model.WIN
			m.marked = m.mines

			data := SocketData{DataType: END_OF_GAME}
			remainingMines := m.GetRemainingMines()

			data.Data = make([]byte, 0, 3+len(remainingMines)*model.SizeOfMineElementInBytes)
			data.Data = append(data.Data, byte(model.WIN))
			data.Data = append(data.Data, utils.ByteConverter.Uint16ToBytes(m.marked)...)

			for _, mine := range remainingMines {
				data.Data = append(data.Data, utils.ByteConverter.MineElementToBytes(*mine)...)
			}

			m.broadcastToClient(data)

			return
		}
	}
}

func (m *MinesweeperServerEngine) GetRemainingMines() []*model.MineElement {
	matrix := make([]*model.MineElement, 0, (m.height*m.width)>>2)

	for rowIndex := range m.matrix {
		for colIndex := range m.matrix[rowIndex] {
			if m.matrix[rowIndex][colIndex].IsHidden() {
				m.matrix[rowIndex][colIndex].PropOff(model.HiddenBits)

				matrix = append(matrix, m.matrix[rowIndex][colIndex])
			}
		}
	}

	return matrix
}

func (m *MinesweeperServerEngine) SetAnimationDuration(animationDuration time.Duration) _minesweeperEngine {
	m.animationDuration = animationDuration

	return m
}

// func (m *MinesweeperServerEngine) GetWidth() int {
// 	return int(m.width)
// }

// func (m *MinesweeperServerEngine) GetHeight() int {
// 	return int(m.height)
// }

// func (m *MinesweeperServerEngine) GetState() model.MinesweeperState {
// 	return m.state
// }

// func (m *MinesweeperServerEngine) GetRevealed() uint16 {
// 	return m.revealed
// }

// func (m *MinesweeperServerEngine) GetMarked() uint16 {
// 	return m.marked
// }

// func (m *MinesweeperServerEngine) GetMines() uint16 {
// 	return m.mines
// }
