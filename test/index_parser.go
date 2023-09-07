package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/bytebuffer_parser/MySchema"
	"github.com/bytebuffer_parser/parsers"
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

	for index00 := 0; index00 < EmployeeArrLen0; index00++ {
		obj.Employee = append(obj.Employee, Employees{})
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
	for i := 0; i < 100000; i++ {
		person.Employee = append(person.Employee, Employees{Name: "Jane Smith", Salary: 60000.0})
	}

	fmt.Println("Encoding........................")

	start := time.Now()
	data := SAMPLE_Encoder(person)
	fmt.Println("bytebuffer encoding: ", time.Since(start))
	fmt.Println("bytebuffer length: ", len(data))

	start = time.Now()
	flatbufferData := createFlatBuffer(person)
	fmt.Println("flatbuffer encoding: ", time.Since(start))
	fmt.Println("flatbuffer length: ", len(flatbufferData))

	start = time.Now()
	jsonE, _ := json.Marshal(person)
	fmt.Println("json encoding: ", time.Since(start))
	fmt.Println("json length: ", len(jsonE))

	start = time.Now()
	xmlStr, _ := xml.Marshal(person)
	fmt.Println("xml encoding: ", time.Since(start))
	fmt.Println("xml length: ", len(xmlStr))

	fmt.Println("Decoding........................")

	start = time.Now()
	SAMPLE_Decoder(data)
	fmt.Println("bytebuffer decoding: ", time.Since(start))

	start = time.Now()
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
	fmt.Println("flatbuffer decoding: ", time.Since(start))

	start = time.Now()
	p := Person{}
	json.Unmarshal(jsonE, &p)
	fmt.Println("json decoding: ", time.Since(start))

	start = time.Now()
	s := Person{}
	xml.Unmarshal(xmlStr, &s)
	fmt.Println("xml decoding: ", time.Since(start))
}
