package main

import (
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bytebuffer_parser/parsers"
)

func main() {

	fileName := flag.String("fileName", "sample", "a string")
	packageName := flag.String("package", "main", "a string")
	language := flag.String("language", "golang", "a string")

	flag.Parse()

	// Define the input class definitions
	fileContent, err := os.ReadFile(*fileName + ".bb")
	if err != nil {
		log.Fatal(err)
	}

	// Convert the byte slice to a string
	contentAsString := string(fileContent)

	// Split the input into separate class definitions
	// classDefs := strings.Split(contentAsString, "\n\n")

	// temporary variables to write file names package names etc...
	totalContent := ""
	finalStruct := ""
	modelFileName := ""
	encoderFileName := ""
	decoderFileName := ""

	rootClassName := parsers.RootClass{}
	globalMap := make(map[string][]string)

	// specific to golang
	switch *language {
	case "golang":
		finalStruct = `package ` + *packageName
		modelFileName = *fileName + ".go"
		encoderFileName = *fileName + "_encoder.go"
		decoderFileName = *fileName + "_decoder.go"

		// generate struct out of it
		totalContent, rootClassName = parsers.GenerateGolangStruct(contentAsString, rootClassName, &globalMap)
		break
	default:
		log.Fatal("Invalid language...")
	}

	// writing to final content variable
	finalStruct += "\n" + totalContent

	// Write the content to the file
	switch *language {
	case "golang":
		parsers.WriteStructData(modelFileName, finalStruct)
		break
	default:
		log.Fatal("Invalid language...")
	}

	// Create the root node
	treeNode := &parsers.TreeNode{Value: rootClassName.Name}
	createTreeNode(treeNode, globalMap, rootClassName.Name)

	currentIterate := 0
	stringDataEncoder := ""
	stringDataDecoder := ""

	switch *language {
	case "golang":

		squareBrackets := ""
		for i := 0; i < rootClassName.ArrayCount; i++ {
			squareBrackets += "[]"
		}

		totalParentBraces := ""
		stringDataEncoder = `

		package ` + *packageName + `

		import(
			"github.com/bytebuffer_parser/parsers"
		)

		func ` + strings.ToUpper(*fileName) + `_Encoder(obj ` + squareBrackets + rootClassName.Name + `) []byte{

			bb := parsers.Buffer{
				FloatIntEncoderVal: 10000.0,
				Endian: "big",
			}
	`

		if squareBrackets != "" {
			for i := 0; i < rootClassName.ArrayCount; i++ {
				if i == 0 {
					stringDataEncoder += `
			bb.PutShort(len(obj))
	`
					stringDataEncoder += `
				for i` + strconv.Itoa(i) + `:=0;i` + strconv.Itoa(i) + `<len(obj);i` + strconv.Itoa(i) + `++{
	`
				} else {
					stringDataEncoder += `
			bb.PutShort(len(obj` + totalParentBraces + `))
	`
					stringDataEncoder += `
				for i` + strconv.Itoa(i) + `:=0;i` + strconv.Itoa(i) + `<len(obj` + totalParentBraces + `);i` + strconv.Itoa(i) + `++{
	`
				}

				totalParentBraces += "[i" + strconv.Itoa(i) + "]"
			}

		}

		totalParentBraces = "obj" + totalParentBraces + "."

		parsers.GenerateGolangEncodeCode(&currentIterate, &stringDataEncoder, treeNode, totalParentBraces)

		if squareBrackets != "" {
			for i := 0; i < rootClassName.ArrayCount; i++ {
				stringDataEncoder += `
			}
				`
			}
		}

		stringDataEncoder += `
			return bb.Array()
		}
	`
		break

	default:
		log.Fatal("Invalid language...")
	}

	// writing encoder data
	// Write the content to the file
	switch *language {
	case "golang":
		parsers.WriteEncoderData(encoderFileName, stringDataEncoder)
		break
	default:
		log.Fatal("Invalid language...")
	}

	switch *language {
	case "golang":

		squareBrackets := ""

		for i := 0; i < rootClassName.ArrayCount; i++ {
			squareBrackets += "[]"
		}

		totalParentBraces := ""

		currentIterate = 0
		rootArrayClass := ""

		if squareBrackets != "" {
			rootArrayClass = "obj := make(" + squareBrackets + rootClassName.Name + ", " + strconv.Itoa(rootClassName.ArrayCount) + ")"
		} else {
			rootArrayClass = "obj := " + rootClassName.Name + "{}"
		}

		stringDataDecoder = `

		package ` + *packageName + `

		import(
			"github.com/bytebuffer_parser/parsers"
		)

		func ` + strings.ToUpper(*fileName) + `_Decoder(byteArr []byte) ` + squareBrackets + rootClassName.Name + `{

			bb := parsers.Buffer{
				FloatIntEncoderVal: 10000.0,
				Endian: "big",
			}

			bb.Wrap(byteArr)

			` + rootArrayClass + `
	`
		if squareBrackets != "" {
			innerbracesCount := rootClassName.ArrayCount - 1
			for i := 0; i < rootClassName.ArrayCount; i++ {
				if i == 0 {
					stringDataDecoder += `
			for i` + strconv.Itoa(i) + `:=0;i` + strconv.Itoa(i) + `<len(obj);i` + strconv.Itoa(i) + `++{
`
				} else {

					nestedSquareBrackets := ""

					for j := 0; j < innerbracesCount; j++ {
						nestedSquareBrackets += "[]"
					}

					innerbracesCount -= 1

					stringDataDecoder += `
					obj` + totalParentBraces + ` = make(` + nestedSquareBrackets + rootClassName.Name + `, len(obj` + totalParentBraces + `))
`
					stringDataDecoder += `
			for i` + strconv.Itoa(i) + `:=0;i` + strconv.Itoa(i) + `<len(obj` + totalParentBraces + `);i` + strconv.Itoa(i) + `++{
`
				}

				totalParentBraces += "[i" + strconv.Itoa(i) + "]"
			}

		}

		totalParentBraces = "obj" + totalParentBraces + "."

		parsers.GenerateGolangDecoderCode(&currentIterate, &stringDataDecoder, treeNode, totalParentBraces)

		if squareBrackets != "" {
			for i := 0; i < rootClassName.ArrayCount; i++ {
				stringDataDecoder += `
			}
				`
			}
		}

		stringDataDecoder += `
			return obj
		}
	`
		break
	default:
		log.Fatal("Invalid language...")
	}

	// write decode code here TODO
	switch *language {
	case "golang":
		parsers.WriteDecoderData(decoderFileName, stringDataDecoder)
		break
	default:
		log.Fatal("Invalid language...")
	}

}

func matchDataType(globalMap map[string][]string, fieldType string) bool {

	// Define a regular expression pattern to match special characters and spaces
	pattern := `[^\w]+`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Replace all matches with an empty string
	result := re.ReplaceAllString(fieldType, "")

	_, ok := globalMap[result]

	if !ok && result != "int" && result != "long" && result != "short" && result != "float" && result != "string" && result != "bool" {
		return false
	}

	return true
}

func createTreeNode(node *parsers.TreeNode, globalMap map[string][]string, rootClassName string) {

	val, ok := globalMap[rootClassName]

	if ok {
		for _, field := range val {

			dataType := strings.Split(field, ":")

			fieldName, fieldType := dataType[0], dataType[1]
			array := 0

			status := matchDataType(globalMap, fieldType)

			if !status {
				log.Fatal("Invalid datatype ", fieldType)
			}

			stat, bracketCount := parsers.CountSquareBrackets(fieldType)
			if !stat {
				log.Fatal("Invalid [] in the bytebuffer file...")
			}
			array = bracketCount

			fieldType = strings.ReplaceAll(fieldType, "[]", "")

			childNode := &parsers.TreeNode{Name: fieldName, Value: fieldType, ArrayCount: array}
			node.Children = append(node.Children, childNode)
			createTreeNode(childNode, globalMap, fieldType)
		}
	}
}
