package main

import (
	"github.com/bytebuffer_parser/parsers"
)

func SAMPLE_Decoder(byteArr []byte) Person {

	obj := Person{}

	bb := parsers.Buffer{
		FloatIntEncoderVal: 10000.0,
		Endian:             "big",
	}

	bb.Wrap(byteArr)

	obj.Epoch = bb.GetLongInteger()

	obj.Watch = bb.GetString()

	obj.Xyz = bb.GetShort()

	obj.Salary = bb.GetFloatUsingIntEncoding()

	EmployeeArrLen0 := bb.GetShort()
	obj.Employee = make([]Employees, EmployeeArrLen0)
	for index00 := 0; index00 < EmployeeArrLen0; index00++ {

		obj.Employee[index00].Name = bb.GetString()

		obj.Employee[index00].Salary = bb.GetFloatUsingIntEncoding()

	}

	return obj
}
