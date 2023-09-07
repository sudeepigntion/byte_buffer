package main

import (
	"github.com/bytebuffer_parser/parsers"
)

func SAMPLE_Encoder(obj Person) []byte {

	bb := parsers.Buffer{
		FloatIntEncoderVal: 10000.0,
		Endian:             "big",
	}

	bb.PutLong(obj.Epoch)

	bb.PutString(obj.Watch)

	bb.PutShort(obj.Xyz)

	bb.PutFloatUsingIntEncoding(obj.Salary)

	bb.PutShort(len(obj.Employee))

	for index00 := range obj.Employee {

		bb.PutString(obj.Employee[index00].Name)

		bb.PutFloatUsingIntEncoding(obj.Employee[index00].Salary)

	}

	return bb.Array()
}
