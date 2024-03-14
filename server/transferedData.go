package server

import (
	"giogo/utils"

	"nhooyr.io/websocket"
)

type DataType byte

const (
	// Ismeretlen adattípus
	UNKNOWN DataType = iota
	// Kattintáskor fordul elő ez az adattípus | client -> server -> client
	POSITION
	// Amikor egy állapot változás történik | server -> client
	STATE
	// Tetszőleges üzenet küldés a szerver irányába, ami azt broadcastolja a többi felhasználó felé
	// | client -> server -> client
	TEXT
	// Amikor a host újraméretezi a képernyőt | client -> server -> client
	RESIZE
	// Játék induláskor, vagy gomb lenyomásra amikor a felhasználók újra akarják kezdeni a játékot | client -> server -> client
	RESTART
	// A kliens ezt az üzenetet küldi a szervernek,
	// az pedig meghívja a 'GetRemainingMines()' metódust és
	// visszaküldi milyen eredménnyel ért véget a meccs valamint
	// a még fel nem fedett cellákat | client -> server -> client
	END_OF_GAME
	SERVER_STATUS
)

// TODO: 'END_OF_GAME' kivétele és logika áthelyezés a 'STATE'-be, lényegében felesleges

func (d DataType) ToString() string {
	switch d {
	case POSITION:
		return "Position"
	case STATE:
		return "State"
	case TEXT:
		return "Text"
	case RESIZE:
		return "Resize"
	case RESTART:
		return "Restart"
	case END_OF_GAME:
		return "End of Game"
	case SERVER_STATUS:
		return "Server Status"
	case UNKNOWN:
		fallthrough
	default:
		return "Unknown"
	}
}

type SocketData struct {
	DataType DataType
	Data     []byte
}

func (d *SocketData) ToBytes() []byte {
	bytes := make([]byte, 0, len(d.Data)+1)

	bytes = append(bytes, byte(d.DataType))
	bytes = append(bytes, d.Data...)

	return bytes
}

func (d *SocketData) FromBytes(bytes []byte) {
	if len(bytes) == 0 {
		d.DataType = UNKNOWN
	}

	d.DataType = DataType(bytes[0])
	d.Data = bytes[1:]
}

type ClientMessage struct {
	connection *websocket.Conn
	socketData *SocketData
}

type ServerStatus struct {
	Joined  int
	Limit   int
	CanJoin bool
}

func (s *ServerStatus) ToBytes() []byte {
	b := make([]byte, 0, 9)

	b = append(b, utils.ByteConverter.IntToBytes(s.Joined)...)
	b = append(b, utils.ByteConverter.IntToBytes(s.Limit)...)
	b = append(b, utils.ByteConverter.BoolToByte(s.CanJoin))

	return b
}

func (s *ServerStatus) FromBytes(b []byte, fromIndex int) {
	s.Joined = utils.ByteConverter.BytesToInt(b, fromIndex)
	s.Limit = utils.ByteConverter.BytesToInt(b, fromIndex+4)
	s.CanJoin = utils.ByteConverter.BytesToBool(b, fromIndex+8)
}
