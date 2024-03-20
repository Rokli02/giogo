package models

import (
	"fmt"
	"giogo/utils"
)

type ServerStatus struct {
	Joined      int
	Limit       int
	CanJoin     bool
	PlayerNames []string
}

func (s *ServerStatus) ToBytes() []byte {
	sizeOfStringArray := 4 + len(s.PlayerNames)
	for _, name := range s.PlayerNames {
		sizeOfStringArray += len(name)
	}
	b := make([]byte, 0, 9+sizeOfStringArray)

	b = append(b, utils.ByteConverter.IntToBytes(s.Joined)...)
	b = append(b, utils.ByteConverter.IntToBytes(s.Limit)...)
	b = append(b, utils.ByteConverter.BoolToByte(s.CanJoin))
	fmt.Println("Playernames in bytes:", utils.ByteConverter.StringArrayToBytes(s.PlayerNames))
	b = append(b, utils.ByteConverter.StringArrayToBytes(s.PlayerNames)...)

	return b
}

func (s *ServerStatus) FromBytes(b []byte, fromIndex int) {
	s.Joined = utils.ByteConverter.BytesToInt(b, fromIndex)
	fromIndex += 4
	s.Limit = utils.ByteConverter.BytesToInt(b, fromIndex)
	fromIndex += 4
	s.CanJoin = utils.ByteConverter.BytesToBool(b, fromIndex)
	fromIndex += 1
	s.PlayerNames = utils.ByteConverter.BytesToStringArray(b, fromIndex)
}
