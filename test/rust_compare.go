package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bytebuffer_parser/parsers"
	"github.com/golang/snappy"
)

type Person struct {
	Epoch    [][]int     `json:"Epoch"`
	Watch    [][]int     `json:"Watch"`
	Xyz      int         `json:"Xyz"`
	Salary   float64     `json:"Salary"`
	Employee []Employees `json:"Employee"`
}

type StudentClass struct {
	Name string `json:"Name"`
}

type Employees struct {
	Name    string         `json:"Name"`
	Salary  float64        `json:"Salary"`
	Student []StudentClass `json:"Student"`
}

func SAMPLE_Encoders(compress string, obj [][][]Person) []byte {

	bb := parsers.Buffer{
		FloatIntEncoderVal: 10000.0,
		Endian:             "big",
	}

	bb.PutShort(len(obj))

	for i0 := 0; i0 < len(obj); i0++ {

		bb.PutShort(len(obj[i0]))

		for i1 := 0; i1 < len(obj[i0]); i1++ {

			bb.PutShort(len(obj[i0][i1]))

			for i2 := 0; i2 < len(obj[i0][i1]); i2++ {

				bb.PutShort(len(obj[i0][i1][i2].Epoch))

				for index00 := range obj[i0][i1][i2].Epoch {

					bb.PutShort(len(obj[i0][i1][i2].Epoch[index00]))

					for index01 := range obj[i0][i1][i2].Epoch[index00] {

						bb.PutLong(obj[i0][i1][i2].Epoch[index00][index01])

					}

				}

				bb.PutShort(len(obj[i0][i1][i2].Watch))

				for index00 := range obj[i0][i1][i2].Watch {

					bb.PutShort(len(obj[i0][i1][i2].Watch[index00]))

					for index01 := range obj[i0][i1][i2].Watch[index00] {

						bb.PutInt(obj[i0][i1][i2].Watch[index00][index01])

					}

				}

				bb.PutInt(obj[i0][i1][i2].Xyz)

				bb.PutFloatUsingIntEncoding(obj[i0][i1][i2].Salary)

				bb.PutShort(len(obj[i0][i1][i2].Employee))

				for index00 := range obj[i0][i1][i2].Employee {

					bb.PutString(obj[i0][i1][i2].Employee[index00].Name)

					bb.PutFloatUsingIntEncoding(obj[i0][i1][i2].Employee[index00].Salary)

					bb.PutShort(len(obj[i0][i1][i2].Employee[index00].Student))

					for index10 := range obj[i0][i1][i2].Employee[index00].Student {

						bb.PutString(obj[i0][i1][i2].Employee[index00].Student[index10].Name)

					}

				}

			}

		}

	}

	// Create a buffer to hold the compressed data
	var compressedData bytes.Buffer

	if compress == "gzip" {
		// Create a gzip writer
		gzipWriter := gzip.NewWriter(&compressedData)
		// Close the gzip writer to flush any remaining data
		gzipWriter.Close()

		// Write the data to the gzip writer
		_, err := gzipWriter.Write(bb.Array())
		if err != nil {
			log.Fatal(err)
		}
	} else if compress == "snappy" {
		compressedData.Write(snappy.Encode(nil, bb.Array()))
	} else {
		compressedData.Write(bb.Array())
	}

	return compressedData.Bytes()
}

func SAMPLE_Decoder1(compress string, byteArr []byte) [][][]Person {

	bb := parsers.Buffer{
		FloatIntEncoderVal: 10000.0,
		Endian:             "big",
	}

	var decompressedData []byte
	var err error
	// Create a gzip reader
	if compress == "gzip" {
		compressedBuffer := bytes.NewReader(byteArr)
		gzipReader, err := gzip.NewReader(compressedBuffer)
		if err != nil {
			log.Fatal(err)
		}
		defer gzipReader.Close()

		// Read the decompressed data
		decompressedData, err = io.ReadAll(gzipReader)
		if err != nil {
			log.Fatal(err)
		}
	} else if compress == "snappy" {
		decompressedData, err = snappy.Decode(nil, byteArr)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		decompressedData = byteArr
	}

	bb.Wrap(decompressedData)

	arrLen := bb.GetShort()
	obj := make([][][]Person, arrLen)

	for i0 := 0; i0 < arrLen; i0++ {

		arrLen1 := bb.GetShort()

		obj[i0] = make([][]Person, arrLen1)

		for i1 := 0; i1 < arrLen1; i1++ {

			arrLen2 := bb.GetShort()

			obj[i0][i1] = make([]Person, arrLen2)

			for i2 := 0; i2 < arrLen2; i2++ {

				EpochArrLen0 := bb.GetShort()
				obj[i0][i1][i2].Epoch = make([][]int, EpochArrLen0)
				for index00 := 0; index00 < EpochArrLen0; index00++ {

					EpochArrLen1 := bb.GetShort()
					obj[i0][i1][i2].Epoch[index00] = make([]int, EpochArrLen1)
					for index01 := 0; index01 < EpochArrLen1; index01++ {

						obj[i0][i1][i2].Epoch[index00][index01] = bb.GetLongInteger()

					}

				}

				WatchArrLen0 := bb.GetShort()
				obj[i0][i1][i2].Watch = make([][]int, WatchArrLen0)
				for index00 := 0; index00 < WatchArrLen0; index00++ {

					WatchArrLen1 := bb.GetShort()
					obj[i0][i1][i2].Watch[index00] = make([]int, WatchArrLen1)
					for index01 := 0; index01 < WatchArrLen1; index01++ {

						obj[i0][i1][i2].Watch[index00][index01] = bb.GetInteger()

					}

				}

				obj[i0][i1][i2].Xyz = bb.GetInteger()

				obj[i0][i1][i2].Salary = bb.GetFloatUsingIntEncoding()

				EmployeeArrLen0 := bb.GetShort()
				obj[i0][i1][i2].Employee = make([]Employees, EmployeeArrLen0)
				for index00 := 0; index00 < EmployeeArrLen0; index00++ {

					obj[i0][i1][i2].Employee[index00].Name = bb.GetString()

					obj[i0][i1][i2].Employee[index00].Salary = bb.GetFloatUsingIntEncoding()

					StudentArrLen0 := bb.GetShort()
					obj[i0][i1][i2].Employee[index00].Student = make([]StudentClass, StudentArrLen0)
					for index10 := 0; index10 < StudentArrLen0; index10++ {

						obj[i0][i1][i2].Employee[index00].Student[index10].Name = bb.GetString()

					}

				}

			}

		}

	}

	return obj
}

func main() {

	// Initialize the structs
	student1 := StudentClass{Name: "Alice"}
	student2 := StudentClass{Name: "Bob"}

	employee1 := Employees{
		Name:    "John",
		Salary:  50000.0,
		Student: []StudentClass{student1, student2},
	}

	person1 := Person{
		Epoch:    [][]int{{1, 2, 3}, {4, 5, 6}},
		Watch:    [][]int{{7, 8, 9}, {10, 11, 12}},
		Xyz:      42,
		Salary:   75000.0,
		Employee: []Employees{employee1},
	}

	for i := 0; i < 10000; i++ {
		person1.Employee = append(person1.Employee, employee1)
	}

	person2 := []Person{person1}
	person3 := [][]Person{person2}
	person4 := [][][]Person{person3}

	start := time.Now()
	data := SAMPLE_Encoders("", person4)
	// fmt.Println(data)
	fmt.Println(len(data))

	SAMPLE_Decoder1("", data)
	fmt.Println("bytebuffer encoding-decoding: ", time.Since(start))
	// fmt.Println(per)
}
