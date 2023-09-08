package main

import (
	"flag"
	"log"
	"os"
	"regexp"
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

	// temporary variables to write file names package names etc...
	totalContent := ""
	finalStruct := ""
	modelFileName := ""
	encoderFileName := ""
	decoderFileName := ""

	rootClassName := parsers.RootClass{}
	globalMap := make(map[string][]string)

	switch *language {
	case "golang":

		stringDataEncoder := ""
		stringDataDecoder := ""

		finalStruct = `package ` + *packageName
		modelFileName = *fileName + ".go"
		encoderFileName = *fileName + "_encoder.go"
		decoderFileName = *fileName + "_decoder.go"

		// generate struct out of it
		totalContent, rootClassName = parsers.GenerateGolangStruct(contentAsString, rootClassName, &globalMap)

		// writing to final content variable
		finalStruct += "\n" + totalContent

		// Create the root node
		treeNode := &parsers.TreeNode{Value: rootClassName.Name}
		createTreeNode(treeNode, globalMap, rootClassName.Name)

		parsers.WriteStructData(modelFileName, finalStruct)

		parsers.EncoderCodeGeneration(rootClassName, &stringDataEncoder, packageName, fileName, treeNode)
		parsers.WriteEncoderData(encoderFileName, stringDataEncoder)

		parsers.DecoderCodeGeneration(rootClassName, &stringDataDecoder, packageName, fileName, treeNode)
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
