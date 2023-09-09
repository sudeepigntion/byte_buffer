package byteSample

import (
	"github.com/bytebuffer_parser/parsers"
	"github.com/golang/snappy"
)

func SAMPLE_Encoder(compression bool, obj [][]Person) []byte {

	bb := parsers.Buffer{
		FloatIntEncoderVal: 10000.0,
		Endian:             "big",
	}

	bb.PutShort(len(obj))

	for i0 := 0; i0 < len(obj); i0++ {

		bb.PutShort(len(obj[i0]))

		for i1 := 0; i1 < len(obj[i0]); i1++ {

			bb.PutLong(obj[i0][i1].Epoch)

			bb.PutShort(len(obj[i0][i1].Watch))

			for index00 := range obj[i0][i1].Watch {

				bb.PutShort(len(obj[i0][i1].Watch[index00]))

				for index01 := range obj[i0][i1].Watch[index00] {

					bb.PutInt(obj[i0][i1].Watch[index00][index01])

				}

			}

			bb.PutInt(obj[i0][i1].Xyz)

			bb.PutFloatUsingIntEncoding(obj[i0][i1].Salary)

			bb.PutShort(len(obj[i0][i1].Employee))

			for index00 := range obj[i0][i1].Employee {

				bb.PutShort(len(obj[i0][i1].Employee[index00]))

				for index11 := range obj[i0][i1].Employee[index00] {

					bb.PutShort(len(obj[i0][i1].Employee[index00][index11]))

					for index22 := range obj[i0][i1].Employee[index00][index11] {

						bb.PutShort(len(obj[i0][i1].Employee[index00][index11][index22]))

						for index33 := range obj[i0][i1].Employee[index00][index11][index22] {

							bb.PutString(obj[i0][i1].Employee[index00][index11][index22][index33].Name)

							bb.PutFloatUsingIntEncoding(obj[i0][i1].Employee[index00][index11][index22][index33].Salary)

							bb.PutShort(len(obj[i0][i1].Employee[index00][index11][index22][index33].Student))

							for index40 := range obj[i0][i1].Employee[index00][index11][index22][index33].Student {

								bb.PutString(obj[i0][i1].Employee[index00][index11][index22][index33].Student[index40].Name)

							}

						}

					}

				}

			}

		}

	}

	var compressedData []byte

	if compression {
		compressedData = snappy.Encode(nil, bb.Array())
	} else {
		compressedData = bb.Array()
	}

	return compressedData
}
