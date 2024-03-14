package server

import (
	"context"
	"fmt"
	"giogo/utils"
	"io"
	"net/http"

	"nhooyr.io/websocket"
)

type MinesweeperServerClient struct {
	Host               string
	Port               uint
	OnClosedConnection func()

	ctx        context.Context
	connection *websocket.Conn
}

func NewMinesweeperServerClient(host string, port uint) *MinesweeperServerClient {
	msc := &MinesweeperServerClient{
		Host: host,
		Port: port,
		ctx:  nil,
	}

	return msc
}

func (msc *MinesweeperServerClient) Join() {
	msc.ctx = context.Background()
	connection, resp, err := websocket.Dial(msc.ctx, fmt.Sprintf("ws://%s:%d%s", msc.Host, msc.Port, websocket_action_path), nil)
	if err != nil {
		// resp.body-ban ott van a részletesebb indok
		if resp != nil {
			fmt.Println("Couldn't join to server [", resp.Status, "]:", err)
		} else {
			fmt.Println("Couldn't join to server:", err)
		}

		msc.OnClosedConnection()

		return
	}

	msc.connection = connection
}

func (msc *MinesweeperServerClient) IsJoined() bool {
	return msc.connection != nil && msc.ctx != nil
}

// mainChannel chan []byte
func (msc *MinesweeperServerClient) ReadData() []byte {
	if !msc.IsJoined() {
		fmt.Println("Websocket is not joined to server!")
		msc.OnClosedConnection()

		return nil
	}

	_, data, err := msc.connection.Read(msc.ctx)
	if err != nil {
		fmt.Println("Client | read error from server:", err)
		msc.Disconnect()
		msc.OnClosedConnection()

		return nil
	}

	return data
}

func (msc *MinesweeperServerClient) WriteData(data []byte) {
	if !msc.IsJoined() {
		fmt.Println("Websocket is not joined to server!")
		msc.OnClosedConnection()

		return
	}

	err := msc.connection.Write(msc.ctx, websocket.MessageText, data)
	if err != nil {
		fmt.Println("Client | write error from server:", err)
		msc.Disconnect()
		msc.OnClosedConnection()

		return
	}
}

func (msc *MinesweeperServerClient) GetStatus() (joined, limit uint8, canJoin bool) {
	response, err := http.Get(fmt.Sprintf("http://%s:%d%s", msc.Host, msc.Port, websocket_status_path))
	if err != nil {
		fmt.Println("Client | error getting status from server", err)

		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Client | error:", err)

		return
	}

	joined = body[0]
	limit = body[1]
	canJoin = utils.ByteConverter.BytesToBool(body, 2)

	return
}

func (msc *MinesweeperServerClient) Disconnect() {
	if !msc.IsJoined() {
		fmt.Println("Websocket is not joined to server!")

		return
	}

	// defer msc.connection.CloseNow()
	msc.connection.Close(websocket.StatusNormalClosure, "Disconnected from server")
	msc.connection = nil
}
