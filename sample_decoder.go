package byteSample

import (
	"github.com/bytebuffer_parser/parsers"
	"github.com/golang/snappy"
	"log"
)

func SAMPLE_Decoder(compression bool, byteArr []byte) [][]Person {

	bb := parsers.Buffer{
		FloatIntEncoderVal: 10000.0,
		Endian:             "big",
	}

	if compression {
		decompressedData, err := snappy.Decode(nil, byteArr)

		if err != nil {
			log.Fatal("Failed to decompress data...")
		}

		bb.Wrap(decompressedData)
	} else {
		bb.Wrap(byteArr)
	}

	arrLen := bb.GetShort()
	obj := make([][]Person, arrLen)

	for i0 := 0; i0 < arrLen; i0++ {

		arrLen1 := bb.GetShort()

		obj[i0] = make([]Person, arrLen1)

		for i1 := 0; i1 < arrLen1; i1++ {

			obj[i0][i1].Epoch = bb.GetLongInteger()

			WatchArrLen0 := bb.GetShort()
			obj[i0][i1].Watch = make([][]int, WatchArrLen0)
			for index00 := 0; index00 < WatchArrLen0; index00++ {

				WatchArrLen1 := bb.GetShort()
				obj[i0][i1].Watch[index00] = make([]int, WatchArrLen1)
				for index01 := 0; index01 < WatchArrLen1; index01++ {

					obj[i0][i1].Watch[index00][index01] = bb.GetInteger()

				}

			}

			obj[i0][i1].Xyz = bb.GetInteger()

			obj[i0][i1].Salary = bb.GetFloatUsingIntEncoding()

			EmployeeArrLen0 := bb.GetShort()
			obj[i0][i1].Employee = make([][][][]Employees, EmployeeArrLen0)
			for index00 := 0; index00 < EmployeeArrLen0; index00++ {

				EmployeeArrLen1 := bb.GetShort()
				obj[i0][i1].Employee[index00] = make([][][]Employees, EmployeeArrLen1)
				for index11 := 0; index11 < EmployeeArrLen1; index11++ {

					EmployeeArrLen2 := bb.GetShort()
					obj[i0][i1].Employee[index00][index11] = make([][]Employees, EmployeeArrLen2)
					for index22 := 0; index22 < EmployeeArrLen2; index22++ {

						EmployeeArrLen3 := bb.GetShort()
						obj[i0][i1].Employee[index00][index11][index22] = make([]Employees, EmployeeArrLen3)
						for index33 := 0; index33 < EmployeeArrLen3; index33++ {

							obj[i0][i1].Employee[index00][index11][index22][index33].Name = bb.GetString()

							obj[i0][i1].Employee[index00][index11][index22][index33].Salary = bb.GetFloatUsingIntEncoding()

							StudentArrLen0 := bb.GetShort()
							obj[i0][i1].Employee[index00][index11][index22][index33].Student = make([]StudentClass, StudentArrLen0)
							for index40 := 0; index40 < StudentArrLen0; index40++ {

								obj[i0][i1].Employee[index00][index11][index22][index33].Student[index40].Name = bb.GetString()

							}

						}

					}

				}

			}

		}

	}

	return obj
}
