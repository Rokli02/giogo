package utils

import (
	"giogo/ui/pages/minesweeper/model"
	"image"
	"reflect"
	"unsafe"
)

type byteConverter uint8

var ByteConverter byteConverter = 0

func (byteConverter) IntToBytes(num int) []byte {
	var b []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh.Len = 4
	sh.Cap = 4
	sh.Data = uintptr(unsafe.Pointer(&num))

	return b[:]
}

func (byteConverter) BytesToInt(b []byte, fromIndex int) int {
	return int(*(*int32)(unsafe.Pointer(&b[fromIndex])))
}

func (byteConverter) Uint16ToBytes(num uint16) []byte {
	var b []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh.Len = 2
	sh.Cap = 2
	sh.Data = uintptr(unsafe.Pointer(&num))

	return b[:]
}

func (byteConverter) BytesToUint16(b []byte, fromIndex int) uint16 {
	return *(*uint16)(unsafe.Pointer(&b[fromIndex]))
}

func (byteConverter) MineElementToBytes(mineElement model.MineElement) []byte {
	var b []byte = make([]byte, 0, 10)

	b = append(b, ByteConverter.PointToBytes(mineElement.Pos)...)
	b = append(b, byte(mineElement.Props))
	b = append(b, byte(mineElement.Value))

	return b
}

func (byteConverter) BytesToMineElement(b []byte, fromIndex int) model.MineElement {
	return model.MineElement{
		Pos:   ByteConverter.BytesToPoint(b, fromIndex),
		Props: model.MineElementProps(b[fromIndex+8]),
		Value: int8(b[fromIndex+9]),
	}
}

func (byteConverter) PointToBytes(point image.Point) []byte {
	var b []byte = make([]byte, 0, 8)

	b = append(b, ByteConverter.IntToBytes(point.X)...)
	b = append(b, ByteConverter.IntToBytes(point.Y)...)

	return b
}

func (byteConverter) BytesToPoint(b []byte, fromIndex int) image.Point {
	return image.Point{
		X: ByteConverter.BytesToInt(b, fromIndex),
		Y: ByteConverter.BytesToInt(b, fromIndex+4),
	}
}

func (byteConverter) BoolToByte(value bool) byte {
	if value {
		return 1
	}

	return 0
}

func (byteConverter) BytesToBool(b []byte, fromIndex int) bool {
	return b[fromIndex] == 1
}
