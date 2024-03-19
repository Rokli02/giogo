package server

import (
	"context"
	"errors"
	"fmt"
	"giogo/server/models"
	"giogo/utils"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"gioui.org/io/pointer"
	"nhooyr.io/websocket"
)

type MinesweeperServer struct {
	server              http.Server
	host                string
	port                uint
	hostConnection      *websocket.Conn
	clients             []*models.ClientData
	connectionsToRemove []*websocket.Conn
	connectionMutext    sync.Mutex
	clientToServer      chan models.ClientMessage
	healthCheckChan     chan uint8

	engine *MinesweeperServerEngine

	ClientLimit uint8
	CanJoin     bool
}

func NewMinesweeperServer(host string, port uint, clientLimit uint8) *MinesweeperServer {
	ms := &MinesweeperServer{
		host:        host,
		port:        port,
		ClientLimit: clientLimit,
		CanJoin:     true,
	}

	ms.engine = NewMinesweeperServerEngine(ms.broadcastToClient)
	ms.engine.SetAnimationDuration(time.Millisecond * 0)

	return ms
}

func (ms *MinesweeperServer) Open() {
	ms.healthCheckChan = make(chan uint8)
	var usedPort uint = ms.port

	for i := 0; i < count_of_port_reservation_tries; i++ {
		if !isPortAvailable(ms.host, usedPort) {
			usedPort += 1
		}
	}

	ms.port = usedPort

	mux := http.NewServeMux()
	mux.HandleFunc(server_status, ms.statusRoute)
	mux.HandleFunc(websocket_status_path, ms.socketStatusRoute)
	mux.HandleFunc(websocket_action_path, ms.socketsRoute)
	mux.HandleFunc(server_health_check, ms.healthCheckRoute)

	ms.server = http.Server{
		Addr:    fmt.Sprintf("%s:%d", ms.host, ms.port),
		Handler: mux,
	}

	ms.clients = make([]*models.ClientData, 0, ms.ClientLimit)
	ms.connectionsToRemove = make([]*websocket.Conn, 0, ms.ClientLimit)
	ms.clientToServer = make(chan models.ClientMessage)

	go ms.handleClientActions()

	if ms.host == Public_Host {
		if address, err := getLocalIp(); err != nil {
			fmt.Println(err)
		} else {
			ms.host = address
		}
	}

	go func() {
		fmt.Printf("Server is listening on ws://%s\n", ms.server.Addr)

		go func() {
			healthCheckTimes := 0
			factorialRate := 1.3
			factorial := factorialRate

			for ms.healthCheckChan != nil && healthCheckTimes < 8 {
				_, err := http.Get(fmt.Sprintf("http://%s:%d%s", ms.host, ms.port, websocket_status_path))
				if err == nil {
					ms.healthCheckChan <- 1

					return
				}

				fmt.Println("Healthcheck error:", err)

				waitTime := time.Millisecond * time.Duration(250*factorial)
				if healthCheckTimes < 5 {
					factorial *= factorialRate
				}

				time.Sleep(waitTime)
				healthCheckTimes++
			}

			panic("Nem indult el a szerver😭🤢👿👿")
		}()

		err := ms.server.ListenAndServe()

		switch err.Error() {
		case http.ErrServerClosed.Error():
			fmt.Println("Szerver leállt")

			return
		default:
			fmt.Println("Ismeretlen baj van a szerverrel")
		}
	}()

	if ms.healthCheckChan != nil {
		<-ms.healthCheckChan
	}
}

func (ms *MinesweeperServer) Close() {
	for _, client := range ms.clients {
		client.Conn.Close(websocket.StatusNormalClosure, "Server is closing...")
	}

	ms.clients = nil

	if ms.clientToServer != nil {
		close(ms.clientToServer)
	}

	ms.clientToServer = nil

	if ms.healthCheckChan != nil {
		close(ms.healthCheckChan)
	}

	ms.healthCheckChan = nil

	ms.server.Close()
}

func (ms *MinesweeperServer) GetHost() string {
	return ms.host
}

func (ms *MinesweeperServer) GetPort() uint {
	return ms.port
}

func (ms *MinesweeperServer) DisableJoin() {
	ms.CanJoin = false
}

func (ms *MinesweeperServer) statusRoute(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(fmt.Sprintf("%d/%d | Can join (%t)", len(ms.clients), ms.ClientLimit, ms.CanJoin)))
}

func (ms *MinesweeperServer) socketStatusRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte{byte(len(ms.clients)), ms.ClientLimit, utils.ByteConverter.BoolToByte(ms.CanJoin)})
}

func (ms *MinesweeperServer) healthCheckRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("good"))
}

func (ms *MinesweeperServer) socketsRoute(w http.ResponseWriter, r *http.Request) {
	connection := ms.handleJoin(&w, r)
	if connection == nil {
		return
	}

	for {
		_, data, err := connection.Read(r.Context())

		if err != nil {
			fmt.Println("Server | read error from client:", err)

			break
		}

		message := models.ClientMessage{
			Connection: connection,
			SocketData: &models.SocketData{},
		}
		message.SocketData.FromBytes(data)

		if ms.clientToServer != nil {
			ms.clientToServer <- message
		}
	}

	ms.removeConnection(connection)

	fmt.Printf("%s disconnected from the server", r.RemoteAddr)
}

func (ms *MinesweeperServer) handleJoin(w *http.ResponseWriter, r *http.Request) *websocket.Conn {
	ms.connectionMutext.Lock()
	defer ms.connectionMutext.Unlock()

	if !ms.CanJoin {
		(*w).WriteHeader(http.StatusLocked)
		(*w).Write([]byte("Can't join to game (server is locked)\n"))

		return nil
	}

	if len(ms.clients) >= int(ms.ClientLimit) {
		(*w).WriteHeader(http.StatusLocked)
		(*w).Write([]byte("Server is full\n"))

		return nil
	}

	connection, err := websocket.Accept(*w, r, nil)
	if err != nil {
		fmt.Println("Error occured in \"handleAction\":", err)

		return nil
	}

	fmt.Printf("%s joined to the lobby\n", r.RemoteAddr)

	if len(ms.clients) == 0 {
		ms.hostConnection = connection
	}

	var username string = ""
	if usernameHeader, hasHeader := r.Header["User-Name"]; hasHeader {
		username = usernameHeader[0]
	}

	ms.clients = append(ms.clients, &models.ClientData{Username: username, Conn: connection})

	var usernames []string = make([]string, 0, len(ms.clients))
	for _, client := range ms.clients {
		usernames = append(usernames, client.Username)
	}

	socketData := models.SocketData{DataType: models.SERVER_STATUS}
	serverStatus := models.ServerStatus{Joined: len(ms.clients), Limit: int(ms.ClientLimit), CanJoin: ms.CanJoin, PlayerNames: usernames}
	socketData.Data = serverStatus.ToBytes()

	ms.broadcastToClient(socketData)

	return connection
}

func (ms *MinesweeperServer) handleClientActions() {
	for message := range ms.clientToServer {
		data := *message.SocketData
		if len(ms.connectionsToRemove) > 0 {
			for i := 0; i < len(ms.connectionsToRemove); i++ {
				ms.connectionsToRemove[i] = nil
			}

			ms.connectionsToRemove = ms.connectionsToRemove[:0]
		}

		fmt.Printf("Server received socketData with type (%s), data length in bytes (%d)\n", data.DataType.ToString(), len(data.Data))
		fmt.Println("\tData content", data.Data)

		switch data.DataType {
		case models.TEXT:
			fmt.Printf("Message from client [%s]: ", data.DataType.ToString())
			fmt.Println(string(data.Data))

			ms.broadcastToClient(data)
		case models.POSITION:
			clickType := pointer.Buttons(data.Data[0])
			pos := utils.ByteConverter.BytesToPoint(data.Data, 1)

			ms.engine.OnPositionAction(pos, clickType)
		case models.RESIZE:
			width := utils.ByteConverter.BytesToUint16(data.Data, 0)
			height := utils.ByteConverter.BytesToUint16(data.Data, 2)
			mines := utils.ByteConverter.BytesToUint16(data.Data, 4)

			ms.engine.Resize(width, height, mines, message.Connection == ms.hostConnection)
		case models.RESTART:
			// TODO: Feature: mind a 4 felhasználónak rá kell nyomnia a RESTART-ra, hogy ténylegesen újrakezdődjön
			ms.engine.Restart()
		}

		for _, connection := range ms.connectionsToRemove {
			ms.removeConnection(connection)
		}
	}
}

func (ms *MinesweeperServer) removeConnection(connection *websocket.Conn) {
	ms.connectionMutext.Lock()
	defer ms.connectionMutext.Unlock()

	if connection == nil {
		return
	}

	if connection == ms.hostConnection {
		ms.Close()

		return
	}

	connectionIndex := -1

	for i, sliceConnection := range ms.clients {
		if sliceConnection.Conn != connection {
			continue
		}

		connectionIndex = i

		break
	}

	length := len(ms.clients)
	if length == 0 || connectionIndex == -1 {
		fmt.Println("Couldn't remove client from slice")

		return
	}

	if length > 1 {
		ms.clients[connectionIndex] = ms.clients[length-1]
		ms.clients[length-1] = nil
	}

	ms.clients = ms.clients[:length-1]

	var usernames []string = make([]string, 0, len(ms.clients))
	for _, client := range ms.clients {
		usernames = append(usernames, client.Username)
	}

	socketData := models.SocketData{DataType: models.SERVER_STATUS}
	serverStatus := models.ServerStatus{Joined: len(ms.clients), Limit: int(ms.ClientLimit), CanJoin: ms.CanJoin, PlayerNames: usernames}
	socketData.Data = serverStatus.ToBytes()

	ms.broadcastToClient(socketData)
}

func (ms *MinesweeperServer) broadcastToClient(data models.SocketData) {
	fmt.Printf("Broadcast %s type data with (%d byte long) to %d clients\n", data.DataType.ToString(), len(data.Data), len(ms.clients))
	fmt.Println("\tData content:", data.Data)
	for _, client := range ms.clients {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)

		err := client.Conn.Write(ctx, websocket.MessageBinary, data.ToBytes())
		cancel()
		if err != nil {
			ms.connectionsToRemove = append(ms.connectionsToRemove, client.Conn)

			continue
		}
	}
}

func isPortAvailable(host string, port uint) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Printf("Port %d is not available!", port)

		return false
	}

	err = listener.Close()
	return err == nil
	// if err != nil {
	// 	return false
	// }

	// return true
}

func getLocalIp() (string, error) {
	conn, err := net.Dial("udp", "192.168.0.0:80")
	if err != nil {
		return Private_Host, errors.New("server | error during local address check")
	}

	defer conn.Close()

	localAddress := strings.Split(conn.LocalAddr().String(), ":")
	if len(localAddress) == 0 {
		return Private_Host, errors.New("server | couldn't split up local address")
	}

	return localAddress[0], nil
}
