package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/bytebuffer_parser/MySchema"
	"github.com/bytebuffer_parser/parsers"
	"github.com/golang/snappy"
	flatbuffers "github.com/google/flatbuffers/go"
)

type Employees struct {
	XMLName xml.Name `xml:"Employee"`
	Name    string   `json:"Name" xml:"Name"`
	Salary  float64  `json:"Salary" xml:"Salary"`
}

type Person struct {
	XMLName  xml.Name    `xml:"Person"`
	Epoch    int         `json:"Epoch" xml:"Epoch"`
	Watch    string      `json:"Watch" xml:"Watch"`
	Xyz      int         `json:"Xyz" xml:"Xyz"`
	Salary   float64     `json:"Salary" xml:"Salary"`
	Employee []Employees `json:"Employee" xml:"Employee"`
}

func SAMPLE_Encoder(compress string, obj Person) []byte {

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

	// Create a buffer to hold the compressed data
	var compressedData bytes.Buffer

	if compress == "gzip" {
		// Create a gzip writer
		gzipWriter := gzip.NewWriter(&compressedData)

		// Write the data to the gzip writer
		_, err := gzipWriter.Write(bb.Array())
		if err != nil {
			fmt.Println("Error writing data to gzip writer:", err)
			return nil
		}

		// Close the gzip writer to flush any remaining data
		gzipWriter.Close()
	} else if compress == "snappy" {
		compressedData.Write(snappy.Encode(nil, bb.Array()))
	} else {
		compressedData.Write(bb.Array())
	}

	return compressedData.Bytes()
}

func SAMPLE_Decoder(compress string, byteArr []byte) Person {

	obj := Person{}
	var decompressedData []byte
	// Create a gzip reader
	if compress == "gzip" {
		compressedBuffer := bytes.NewReader(byteArr)
		gzipReader, err := gzip.NewReader(compressedBuffer)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return obj
		}
		defer gzipReader.Close()

		// Read the decompressed data
		decompressedData, err = io.ReadAll(gzipReader)
		if err != nil {
			fmt.Println("Error reading decompressed data:", err)
			return obj
		}
	} else if compress == "snappy" {
		decompressedData, _ = snappy.Decode(nil, byteArr)
	} else {
		decompressedData = byteArr
	}

	bb := parsers.Buffer{
		FloatIntEncoderVal: 10000.0,
		Endian:             "big",
	}

	bb.Wrap(decompressedData)

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

func createFlatBuffer(person Person) []byte {

	builder := flatbuffers.NewBuilder(0)

	// Create an Employees object
	var employeesArr []flatbuffers.UOffsetT

	for _, val := range person.Employee {
		employeesName := builder.CreateString(val.Name)
		employeesSalary := val.Salary
		MySchema.EmployeesStart(builder)
		MySchema.EmployeesAddName(builder, employeesName)
		MySchema.EmployeesAddSalary(builder, float32(employeesSalary))
		employees := MySchema.EmployeesEnd(builder)
		employeesArr = append(employeesArr, employees)

	}

	// Create a vector of Employees
	employeesVector := builder.CreateVectorOfTables(employeesArr)

	// Create a Person object
	epoch := int64(person.Epoch)
	watch := builder.CreateString(person.Watch)
	xyz := int16(person.Xyz)
	salary := person.Salary
	MySchema.PersonStart(builder)
	MySchema.PersonAddEpoch(builder, epoch)
	MySchema.PersonAddWatch(builder, watch)
	MySchema.PersonAddXyz(builder, xyz)
	MySchema.PersonAddSalary(builder, float32(salary))
	MySchema.PersonAddEmployee(builder, employeesVector)
	personFlat := MySchema.PersonEnd(builder)

	builder.Finish(personFlat)

	// Serialize the FlatBuffer to a byte slice
	data := builder.FinishedBytes()

	return data
}

func main() {

	person := Person{
		Epoch:  1630866000,
		Watch:  "Smartwatch",
		Xyz:    42,
		Salary: 75000.0,
		Employee: []Employees{
			{Name: "John Doe", Salary: 50000.0},
			{Name: "Jane Smith", Salary: 60000.0},
		},
	}

	person.Employee = []Employees{}
	for i := 0; i < 10000000; i++ {
		person.Employee = append(person.Employee, Employees{Name: "Jane Smith", Salary: 60000.0})
	}

	iteration := 1

	fmt.Println("Encoding-decoding........................")

	start := time.Now()
	for i := 0; i < iteration; i++ {
		data := SAMPLE_Encoder("gzip", person)

		if i == 0 {
			fmt.Println("bytebuffer length: ", len(data))
		}
		SAMPLE_Decoder("gzip", data)
	}
	fmt.Println("bytebuffer encoding-decoding: ", time.Since(start))

	start = time.Now()
	for i := 0; i < iteration; i++ {
		flatbufferData := createFlatBuffer(person)
		if i == 0 {
			fmt.Println("flatbuffer length: ", len(flatbufferData))
		}
		personRead := MySchema.GetRootAsPerson(flatbufferData, 0)
		personRead.Epoch()
		personRead.Watch()
		personRead.Xyz()
		personRead.Salary()
		for i := 0; i < personRead.EmployeeLength(); i++ {
			employeesRead := new(MySchema.Employees)
			if personRead.Employee(employeesRead, i) {
			}
		}
	}
	fmt.Println("flatbuffer encoding-decoding: ", time.Since(start))

	start = time.Now()
	for i := 0; i < iteration; i++ {
		jsonE, _ := json.Marshal(person)
		// jsonE = snappy.Encode(nil, jsonE)
		if i == 0 {
			fmt.Println("json length: ", len(jsonE))
		}
		p := Person{}
		json.Unmarshal(jsonE, &p)
	}
	fmt.Println("json encoding-decoding: ", time.Since(start))

	start = time.Now()
	for i := 0; i < iteration; i++ {
		xmlStr, _ := xml.Marshal(person)
		// xmlStr = snappy.Encode(nil, xmlStr)
		if i == 0 {
			fmt.Println("xml length: ", len(xmlStr))
		}
		s := Person{}
		xml.Unmarshal(xmlStr, &s)
	}
	fmt.Println("xml encoding-decoding: ", time.Since(start))

}
