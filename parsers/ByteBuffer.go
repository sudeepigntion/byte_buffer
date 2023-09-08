package parsers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

type Buffer struct {
	packetBuffer       bytes.Buffer
	byteData           []byte
	Endian             string
	Offset             int
	FloatIntEncoderVal float64
}

func (obj *Buffer) Wrap(data []byte) {
	obj.byteData = data
}

func (obj *Buffer) GetInt() []byte {

	intValue := obj.byteData[obj.Offset : obj.Offset+4]

	obj.Offset += 4

	return intValue

}

func (obj *Buffer) GetShortInt() []byte {

	intValue := obj.byteData[obj.Offset : obj.Offset+2]

	obj.Offset += 2

	return intValue

}

func (obj *Buffer) GetInteger() int {

	intValue := obj.byteData[obj.Offset : obj.Offset+4]

	obj.Offset += 4

	return int(binary.BigEndian.Uint32(intValue))
}

func (obj *Buffer) GetLongInteger() int {

	longValue := obj.byteData[obj.Offset : obj.Offset+8]

	obj.Offset += 8

	return int(binary.BigEndian.Uint64(longValue))
}

func (obj *Buffer) GetShort() int {

	shortValue := obj.byteData[obj.Offset : obj.Offset+2]

	obj.Offset += 2

	return int(binary.BigEndian.Uint16(shortValue))
}

func (obj *Buffer) GetLong() []byte {

	longValue := obj.byteData[obj.Offset : obj.Offset+8]

	obj.Offset += 8

	return longValue

}

func (obj *Buffer) Get(size int) []byte {

	var value = obj.byteData[obj.Offset : obj.Offset+size]

	obj.Offset += size

	return value
}

func (obj *Buffer) GetByte() []byte {

	var value = obj.byteData[obj.Offset : obj.Offset+1]

	obj.Offset += 1

	return value

}

func (obj *Buffer) GetString() string {

	typeString := obj.GetByte()[0]
	stringData := ""
	switch typeString {
	case 1:
		strLen := obj.GetByte()[0]
		stringData = string(obj.Get(int(strLen)))
	case 2:
		strLen := int(binary.BigEndian.Uint16(obj.GetShortInt()))
		stringData = string(obj.Get(strLen))
	case 3:
		strLen := int(binary.BigEndian.Uint32(obj.GetInt()))
		stringData = string(obj.Get(strLen))
	case 4:
		strLen := int(binary.BigEndian.Uint64(obj.GetLong()))
		stringData = string(obj.Get(strLen))
	default:
		fmt.Println("Invalid string length type...")
	}

	if stringData == "X" {
		return ""
	} else {
		return stringData
	}
}

func (obj *Buffer) PutString(value string) {
	strLen := len(value)
	if strLen > 0 {
		if strLen < 128 {
			obj.PutByte(1)
			obj.PutByte(byte(strLen))
			obj.Put([]byte(value))
		} else if strLen < 32768 {
			obj.PutByte(2)
			obj.PutShort(strLen)
			obj.Put([]byte(value))
		} else if strLen < 2147483648 {
			obj.PutByte(3)
			obj.PutInt(strLen)
			obj.Put([]byte(value))
		} else {
			obj.PutByte(4)
			obj.PutLong(strLen)
			obj.Put([]byte(value))
		}
	} else {
		obj.PutByte(1)
		obj.PutByte(1)
		obj.Put([]byte("X"))
	}
}

func (obj *Buffer) GetBoolean() bool {

	return obj.GetByte()[0] == 1
}

func (obj *Buffer) PutShort(value int) {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return

	}

	var buff = make([]byte, 2)

	if obj.Endian == "big" {

		binary.BigEndian.PutUint16(buff, uint16(value))

	} else if obj.Endian == "little" {

		binary.LittleEndian.PutUint16(buff, uint16(value))

	}

	obj.packetBuffer.Write(buff)

}

func (obj *Buffer) PutInt(value int) {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return

	}

	var buff = make([]byte, 4)

	if obj.Endian == "big" {

		binary.BigEndian.PutUint32(buff, uint32(value))

	} else if obj.Endian == "little" {

		binary.LittleEndian.PutUint32(buff, uint32(value))

	}

	obj.packetBuffer.Write(buff)

}

func (obj *Buffer) PutLong(value int) {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return

	}

	var buff = make([]byte, 8)

	if obj.Endian == "big" {

		binary.BigEndian.PutUint64(buff, uint64(value))

	} else if obj.Endian == "little" {

		binary.LittleEndian.PutUint64(buff, uint64(value))

	}

	obj.packetBuffer.Write(buff)

}

func (obj *Buffer) PutFloat(value float32) {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return

	}

	var bits = math.Float32bits(value)

	var buff = make([]byte, 4)

	if obj.Endian == "big" {

		binary.BigEndian.PutUint32(buff, bits)

	} else if obj.Endian == "little" {

		binary.LittleEndian.PutUint32(buff, bits)

	}

	obj.packetBuffer.Write(buff)

}

func (obj *Buffer) PutFloatUsingIntEncoding(value float64) {
	obj.PutLong(int(value * obj.FloatIntEncoderVal))
}

func (obj *Buffer) GetFloatUsingIntEncoding() float64 {
	return float64(binary.BigEndian.Uint64(obj.GetLong())) / obj.FloatIntEncoderVal
}

func (obj *Buffer) PutDouble(value float64) {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return

	}

	var bits = math.Float64bits(value)

	var buff = make([]byte, 8)

	if obj.Endian == "big" {

		binary.BigEndian.PutUint64(buff, bits)

	} else if obj.Endian == "little" {

		binary.LittleEndian.PutUint64(buff, bits)

	}

	obj.packetBuffer.Write(buff)

}

func (obj *Buffer) Put(value []byte) {

	obj.packetBuffer.Write(value)

}

func (obj *Buffer) PutByte(value byte) {

	var tempByte []byte

	tempByte = append(tempByte, value)

	obj.packetBuffer.Write(tempByte)

}

func (obj *Buffer) PutBoolean(value bool) {
	if value {
		obj.PutByte(1)
	} else {
		obj.PutByte(0)
	}
}

func (obj *Buffer) Array() []byte {

	return obj.packetBuffer.Bytes()

}

func (obj *Buffer) Size() int {

	return len(obj.packetBuffer.Bytes())

}

func (obj *Buffer) Flip() {

	var bytesArr = obj.packetBuffer.Bytes()

	for i, j := 0, len(bytesArr)-1; i < j; i, j = i+1, j-1 {

		bytesArr[i], bytesArr[j] = bytesArr[j], bytesArr[i]

	}

	var byteBuffer bytes.Buffer

	byteBuffer.Write(bytesArr)

	obj.packetBuffer = byteBuffer
}

func (obj *Buffer) Clear() {

	obj.packetBuffer.Reset()

}

func (obj *Buffer) Slice(start int, end int) error {

	var bytesArr = obj.packetBuffer.Bytes()

	if len(bytesArr) < (start + end) {
		return errors.New("Buffer does not contain that much of limit")
	}

	bytesArr = bytesArr[start:end]

	var byteBuffer bytes.Buffer

	byteBuffer.Write(bytesArr)

	obj.packetBuffer = byteBuffer

	return nil

}

func (obj *Buffer) Bytes2Str(data []byte) string {

	return string(data)

}

func (obj *Buffer) Str2Bytes(data string) []byte {

	return []byte(data)

}

func (obj *Buffer) Bytes2Short(data []byte) uint16 {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return 0

	}

	if obj.Endian == "big" {

		return binary.BigEndian.Uint16(data)

	} else if obj.Endian == "little" {

		return binary.LittleEndian.Uint16(data)

	}

	return 0
}

func (obj *Buffer) Bytes2Int(data []byte) uint32 {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return 0

	}

	if obj.Endian == "big" {

		return binary.BigEndian.Uint32(data)

	} else if obj.Endian == "little" {

		return binary.LittleEndian.Uint32(data)

	}

	return 0
}

func (obj *Buffer) Bytes2Long(data []byte) uint64 {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return 0

	}

	if obj.Endian == "big" {

		return binary.BigEndian.Uint64(data)

	} else if obj.Endian == "little" {

		return binary.LittleEndian.Uint64(data)

	}

	return 0
}

func (obj *Buffer) Short2Bytes(data uint16) []byte {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return nil

	}

	bs := make([]byte, 2)

	if obj.Endian == "big" {

		binary.BigEndian.PutUint16(bs, data)

	} else if obj.Endian == "little" {

		binary.LittleEndian.PutUint16(bs, data)

	}

	return bs
}

func (obj *Buffer) Int2Bytes(data uint32) []byte {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return nil

	}

	bs := make([]byte, 4)

	if obj.Endian == "big" {

		binary.BigEndian.PutUint32(bs, data)

	} else if obj.Endian == "little" {

		binary.LittleEndian.PutUint32(bs, data)

	}

	return bs

}

func (obj *Buffer) Long2Bytes(data uint64) []byte {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return nil

	}

	bs := make([]byte, 8)

	if obj.Endian == "big" {

		binary.BigEndian.PutUint64(bs, data)

	} else if obj.Endian == "little" {

		binary.LittleEndian.PutUint64(bs, data)

	}

	return bs

}

func (obj *Buffer) Bytes2Float(bytes []byte) float32 {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return 0.0

	}

	if obj.Endian == "big" {

		bits := binary.BigEndian.Uint32(bytes)

		float := math.Float32frombits(bits)

		return float

	} else if obj.Endian == "little" {

		bits := binary.LittleEndian.Uint32(bytes)

		float := math.Float32frombits(bits)

		return float

	}

	return 0.0
}

func (obj *Buffer) Float2Bytes(float float32) []byte {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return nil

	}

	if obj.Endian == "big" {

		bits := math.Float32bits(float)

		bytes := make([]byte, 4)

		binary.BigEndian.PutUint32(bytes, bits)

		return bytes

	} else if obj.Endian == "little" {

		bits := math.Float32bits(float)

		bytes := make([]byte, 4)

		binary.LittleEndian.PutUint32(bytes, bits)

		return bytes

	}

	return nil
}

func (obj *Buffer) Bytes2Double(bytes []byte) float64 {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return 0

	}

	if obj.Endian == "big" {

		bits := binary.BigEndian.Uint64(bytes)

		float := math.Float64frombits(bits)

		return float

	} else if obj.Endian == "little" {

		bits := binary.LittleEndian.Uint64(bytes)

		float := math.Float64frombits(bits)

		return float

	}

	return 0
}

func (obj *Buffer) Double2Bytes(float float64) []byte {

	if obj.Endian != "big" && obj.Endian != "little" {

		fmt.Println(ColorRed, "Invalid endianness, must be big or little")
		return nil

	}

	if obj.Endian == "big" {

		bits := math.Float64bits(float)

		bytes := make([]byte, 8)

		binary.BigEndian.PutUint64(bytes, bits)

		return bytes

	} else if obj.Endian == "little" {

		bits := math.Float64bits(float)

		bytes := make([]byte, 8)

		binary.LittleEndian.PutUint64(bytes, bits)

		return bytes

	}

	return nil

}
