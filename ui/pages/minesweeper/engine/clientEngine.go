package engine

import (
	"fmt"
	serverModule "giogo/server"
	"giogo/server/models"
	"giogo/utils"
	"image"
	"time"

	"gioui.org/app"
	"gioui.org/io/pointer"

	"giogo/ui/pages/minesweeper/model"
)

type MinesweeperClientEngine struct {
	w *app.Window

	width        uint16
	height       uint16
	maxMines     uint16
	state        model.MinesweeperState
	revealed     uint16
	marked       uint16
	Client       *serverModule.MinesweeperServerClient
	ServerStatus *models.ServerStatus

	mineChannel    chan model.MineElement
	acks           chan uint8
	serverToClient chan models.SocketData
	engineCommand  chan EngineCommand
	mines          uint16
	elementList    []*model.MineElement
}

// Interface implementation check
var _ MinesweeperEngine = (*MinesweeperClientEngine)(nil)

func NewMinesweeperClientEngine(w *app.Window, host string, port uint) *MinesweeperClientEngine {
	m := &MinesweeperClientEngine{
		w:            w,
		state:        model.WAITING,
		ServerStatus: &models.ServerStatus{},
	}

	m.Client = serverModule.NewMinesweeperServerClient(host, port)
	m.Client.OnClosedConnection = func() {
		m.Close()
		m.engineCommand <- GO_BACK
	}

	return m
}

func (m *MinesweeperClientEngine) Initialize() {
	m.serverToClient = make(chan models.SocketData)
	m.Client.Join()

	go func() {
		for {
			if !m.Client.IsJoined() {
				m.Close()

				return
			}

			readData := m.Client.ReadData()
			if readData == nil {
				m.Client.OnClosedConnection()

				return
			}

			socketData := models.SocketData{}
			socketData.FromBytes(readData)

			if socketData.DataType == models.SERVER_STATUS {
				m.ServerStatus.FromBytes(socketData.Data, 0)

				m.w.Invalidate()

				continue
			}

			m.serverToClient <- socketData
		}
	}()

	go func() {
		fmt.Println("Starting listeing to server responses...")

		for socketData := range m.serverToClient {
			// fmt.Printf("Client received socketData with type (%s), data length in bytes (%d)\n", socketData.DataType.ToString(), len(socketData.Data))
			// fmt.Println("\tData content", socketData.Data)

			switch socketData.DataType {
			case models.STATE:
				state := model.MinesweeperState(socketData.Data[0])

				switch state {
				case model.WAITING:
				case model.START:
					m.width = utils.ByteConverter.BytesToUint16(socketData.Data, 1)
					m.height = utils.ByteConverter.BytesToUint16(socketData.Data, 3)
					m.maxMines = utils.ByteConverter.BytesToUint16(socketData.Data, 5)
					m.mines = utils.ByteConverter.BytesToUint16(socketData.Data, 5)
				case model.LOADING:
					if state != m.state {
						m.state = state
						m.w.Invalidate()
					}

					m.marked = utils.ByteConverter.BytesToUint16(socketData.Data, 1)
					element := utils.ByteConverter.BytesToMineElement(socketData.Data, 3)

					m.mineChannel <- element
					<-m.acks
				case model.RUNNING:
					m.marked = utils.ByteConverter.BytesToUint16(socketData.Data, 1)
				case model.LOSE:
					fallthrough
				case model.WIN:
					m.w.Invalidate()
				}

				m.state = state
			case models.END_OF_GAME:
				m.state = model.MinesweeperState(socketData.Data[0])
				m.marked = utils.ByteConverter.BytesToUint16(socketData.Data, 1)

				clear(m.elementList)
				m.elementList = make([]*model.MineElement, 0)

				for i := 3; i < len(socketData.Data); i += model.SizeOfMineElementInBytes {
					element := utils.ByteConverter.BytesToMineElement(socketData.Data, i)

					m.elementList = append(m.elementList, &element)
				}

				if m.state == model.WIN {
					m.engineCommand <- AFTER_CLICK_WIN

					break
				}

				m.engineCommand <- AFTER_CLICK_LOSE
			case models.RESIZE:
				m.state = model.MinesweeperState(socketData.Data[0])
				m.marked = utils.ByteConverter.BytesToUint16(socketData.Data, 1)
				m.width = utils.ByteConverter.BytesToUint16(socketData.Data, 3)
				m.height = utils.ByteConverter.BytesToUint16(socketData.Data, 5)
				m.mines = utils.ByteConverter.BytesToUint16(socketData.Data, 7)
				m.maxMines = m.mines

				m.engineCommand <- RESIZE
			case models.RESTART:
				m.state = model.MinesweeperState(socketData.Data[0])
				m.marked = utils.ByteConverter.BytesToUint16(socketData.Data, 1)

				m.engineCommand <- RESTART
			case models.POSITION:
				m.state = model.MinesweeperState(socketData.Data[0])
				m.marked = utils.ByteConverter.BytesToUint16(socketData.Data, 1)
				element := utils.ByteConverter.BytesToMineElement(socketData.Data, 3)

				m.mineChannel <- element
				<-m.acks
				m.w.Invalidate()
			case models.UNKNOWN:
				fmt.Println("Ismeretlen adattípus érkezett")
			}
		}

		fmt.Println("Stoped listeing to server...")
	}()
}

func (m *MinesweeperClientEngine) Resize(width uint16, height uint16, mines uint16) {
	socketData := models.SocketData{DataType: models.RESIZE}

	socketData.Data = make([]byte, 0, 6)
	socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(width)...)
	socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(height)...)
	socketData.Data = append(socketData.Data, utils.ByteConverter.Uint16ToBytes(mines)...)

	m.Client.WriteData(socketData.ToBytes())
}

func (m *MinesweeperClientEngine) Restart() {
	socketData := models.SocketData{DataType: models.RESTART}

	m.Client.WriteData(socketData.ToBytes())
}

func (m *MinesweeperClientEngine) Close() {
	if m.serverToClient != nil {
		close(m.serverToClient)
	}

	m.serverToClient = nil

	m.mineChannel = nil
	m.Client.Disconnect()
}

func (m *MinesweeperClientEngine) OnButtonClick(pos image.Point, clickType pointer.Buttons) {
	socketData := models.SocketData{DataType: models.POSITION}

	socketData.Data = make([]byte, 0, 9)
	socketData.Data = append(socketData.Data, byte(clickType))
	socketData.Data = append(socketData.Data, utils.ByteConverter.PointToBytes(pos)...)

	fmt.Println("Client Button click pos:", pos, " | data:", socketData.Data)

	m.Client.WriteData(socketData.ToBytes())
}

func (m *MinesweeperClientEngine) GetRemainingMines() []*model.MineElement {
	elementListCopy := make([]*model.MineElement, 0, len(m.elementList))

	for _, element := range m.elementList {
		elementCopy := *element

		elementListCopy = append(elementListCopy, &elementCopy)
	}

	m.elementList = nil
	m.elementList = make([]*model.MineElement, 0)

	return elementListCopy
}

func (m *MinesweeperClientEngine) SetAnimationDuration(_ time.Duration) MinesweeperEngine {
	return m
}

func (m *MinesweeperClientEngine) SetChannels(mainChannel chan model.MineElement, acks chan uint8, engineCommand chan EngineCommand) MinesweeperEngine {
	m.mineChannel = mainChannel
	m.acks = acks
	m.engineCommand = engineCommand

	return m
}

func (m *MinesweeperClientEngine) GetWidth() int {
	return int(m.width)
}

func (m *MinesweeperClientEngine) GetHeight() int {
	return int(m.height)
}

func (m *MinesweeperClientEngine) GetMarked() uint16 {
	return m.marked
}

func (m *MinesweeperClientEngine) GetMines() uint16 {
	return m.mines
}

func (m *MinesweeperClientEngine) GetRevealed() uint16 {
	return m.revealed
}

func (m *MinesweeperClientEngine) GetState() model.MinesweeperState {
	return m.state
}
