package parsers

import (
	"fmt"
	"go/format"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GenerateGolangStruct(classDefinitions string, rootClassName RootClass, globalMap *map[string][]string) (string, RootClass) {

	tempMap := *globalMap

	structs := make(map[string][]StructField)
	lines := strings.Split(classDefinitions, "\n")

	currentStructName := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "class") {
			// Start of a new struct definition
			currentStructName = strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(line, "class")), "{")
		} else if strings.HasPrefix(line, "}") {
			// End of the current struct definition
			currentStructName = ""
		} else if currentStructName != "" && strings.Contains(line, " ") {
			// This line defines a field
			parts := strings.Fields(line)

			if len(parts) == 2 {
				fieldName := parts[1]
				fieldType := parts[0]
				currentStructName = strings.TrimSpace(currentStructName)

				pattern := `^[^a-zA-Z]`
				re := regexp.MustCompile(pattern)

				if re.MatchString(currentStructName) {
					log.Fatal("Invalid class name ", currentStructName)
				}

				if re.MatchString(fieldName) {
					log.Fatal("Invalid field name ", fieldName)
				}

				_, ok := tempMap[currentStructName]

				if !ok {
					tempMap[currentStructName] = make([]string, 0)
				}

				tempMap[currentStructName] = append(tempMap[currentStructName], fieldName+":"+fieldType)

				fieldType = strings.ReplaceAll(fieldType, "long", "int")
				fieldType = strings.ReplaceAll(fieldType, "short", "int")
				fieldType = strings.ReplaceAll(fieldType, "float", "float64")

				structs[currentStructName] = append(structs[currentStructName], StructField{Name: fieldName, Type: fieldType})
			}
		} else {

			line = strings.ReplaceAll(line, " ", "")

			lines = strings.Split(line, "=")

			if len(lines) != 2 {
				continue
			}

			stat, arrayCount := CountSquareBrackets(lines[1])

			if !stat {
				log.Fatal("Invalid array [] in root export class ", lines[1])
			}
			rootClassName.ArrayCount = arrayCount

			rootClassName.Name = strings.ReplaceAll(lines[1], "[]", "")
		}
	}

	goCode := ""
	for structName, fields := range structs {
		goCode += fmt.Sprintf("type %s struct {\n", structName)
		for _, field := range fields {
			goCode += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", field.Name, field.Type, field.Name)
		}
		goCode += "}\n\n"
	}

	return goCode, rootClassName
}

func GenerateGolangEncodeCode(currentIterate *int, stringData *string, node *TreeNode, parentName string) {

	for _, child := range node.Children {

		path := parentName + child.Name

		switch child.Value {
		case "int":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += `
				    bb.PutShort(len(` + path + `))
				`
						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + `{
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(len(` + path + loopSquares + `))
						`
						*stringData += `
								for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + loopSquares + `{
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutInt(` + path + squares + `)
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutInt(` + path + `)
			`
			}
		case "long":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += `
				    bb.PutShort(len(` + path + `))
				`
						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + `{
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(len(` + path + loopSquares + `))
						`
						*stringData += `
								for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + loopSquares + `{
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutLong(` + path + squares + `)
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutLong(` + path + `)
			`
			}
		case "short":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += `
				    bb.PutShort(len(` + path + `))
				`
						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + `{
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(len(` + path + loopSquares + `))
						`
						*stringData += `
								for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + loopSquares + `{
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutShort(` + path + squares + `)
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutShort(` + path + `)
			`
			}
		case "string":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += `
				    bb.PutShort(len(` + path + `))
				`
						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + `{
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(len(` + path + loopSquares + `))
						`
						*stringData += `
								for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + loopSquares + `{
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutString(` + path + squares + `)
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutString(` + path + `)
			`
			}
		case "float":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += `
				    bb.PutShort(len(` + path + `))
				`
						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + `{
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(len(` + path + loopSquares + `))
						`
						*stringData += `
								for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + loopSquares + `{
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutFloatUsingIntEncoding(` + path + squares + `)
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutFloatUsingIntEncoding(` + path + `)
			`
			}
		case "bool":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += `
				    bb.PutShort(len(` + path + `))
				`
						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + `{
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(len(` + path + loopSquares + `))
						`
						*stringData += `
								for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + loopSquares + `{
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutBoolean(` + path + squares + `)
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
			bb.PutBoolean(` + path + `)
		`
			}
		default:
			if child.ArrayCount > 0 {
				squares := ""
				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				    bb.PutShort(len(` + path + `))
				`
					*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` := range ` + path + `{
						`
					squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

					path += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					*currentIterate += 1
				}
			}

			GenerateGolangEncodeCode(currentIterate, stringData, child, path+".")
		}
	}

	if node.ArrayCount > 0 {
		for i := 0; i < node.ArrayCount; i++ {
			*stringData += `
		}
			`
		}

	}
}

func GenerateGolangDecoderCode(currentIterate *int, stringData *string, node *TreeNode, parentName string) {

	for _, child := range node.Children {

		path := parentName + child.Name

		switch child.Value {
		case "int":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
				    ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
				`
						*stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
						` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
						`

						*stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetInteger()
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.GetInteger()
			`
			}
		case "long":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
				    ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
				`
						*stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
						` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
						`

						*stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetLongInteger()
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.GetLongInteger()
			`
			}
		case "short":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
				    ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
				`
						*stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
						` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
						`

						*stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetShort()
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.GetShort()
			`
			}
		case "string":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
				    ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
				`
						*stringData += path + squares + `= make(` + arrayBraces + `string, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
						` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
						`

						*stringData += path + squares + `= make(` + arrayBraces + `string, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetString()
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.GetString()
			`
			}
		case "float":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
				    ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
				`
						*stringData += path + squares + `= make(` + arrayBraces + `float64, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
						` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
						`

						*stringData += path + squares + `= make(` + arrayBraces + `float64, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetFloatUsingIntEncoding()
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
					` + path + ` = bb.GetFloatUsingIntEncoding()
				`
			}
		case "bool":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
				    ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
				`
						*stringData += path + squares + `= make(` + arrayBraces + `bool, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"
						}

						*stringData += `
						` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
						`

						*stringData += path + squares + `= make(` + arrayBraces + `bool, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetBoolean()
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
						` + path + ` = bb.GetBoolean()
					`
			}
		default:

			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {

					arrayBraces := ""

					for j := 0; j < arrayCount; j++ {
						arrayBraces += "[]"
					}

					*stringData += `
					` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := bb.GetShort()
				`

					*stringData += path + `= make(` + arrayBraces + child.Value + `, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

					*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `:=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++{
						`

					squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

					path += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					*currentIterate += 1
					arrayCount -= 1
				}
			}

			GenerateGolangDecoderCode(currentIterate, stringData, child, path+".")
		}
	}

	if node.ArrayCount > 0 {
		for i := 0; i < node.ArrayCount; i++ {
			*stringData += `
		}
			`
		}

	}
}

func WriteEncoderData(encoderFileName string, stringDataEncoder string) {

	code, err := format.Source([]byte(stringDataEncoder))
	if err != nil {
		log.Fatal(err)
	}

	// creating encoder file
	// Create the file (or truncate it if it already exists)
	file, err := os.Create(encoderFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Close the file when we're done

	_, err = file.WriteString(string(code))
	if err != nil {
		log.Fatal(err)
	}
}

func WriteDecoderData(decoderFileName string, stringDataDecoder string) {

	code, err := format.Source([]byte(stringDataDecoder))
	if err != nil {
		log.Fatal(err)
	}

	// creating encoder file
	// Create the file (or truncate it if it already exists)
	file, err := os.Create(decoderFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Close the file when we're done

	_, err = file.WriteString(string(code))
	if err != nil {
		log.Fatal(err)
	}
}

func WriteStructData(modelFileName string, finalStruct string) {

	code, err := format.Source([]byte(finalStruct))
	if err != nil {
		log.Fatal(err)
	}

	// creating model file for example it will contain struct or class file
	file, err := os.Create(modelFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Close the file when we're done

	_, err = file.WriteString(string(code))
	if err != nil {
		log.Fatal(err)
	}
}

func EncoderCodeGeneration(rootClassName RootClass, stringDataEncoder *string, packageName *string, fileName *string, treeNode *TreeNode) {

	currentIterate := 0

	squareBrackets := ""
	for i := 0; i < rootClassName.ArrayCount; i++ {
		squareBrackets += "[]"
	}

	totalParentBraces := ""
	*stringDataEncoder = `

		package ` + *packageName + `

		import(
			"github.com/golang/snappy"
			"github.com/bytebuffer_parser/parsers"
		)

		func ` + strings.ToUpper(*fileName) + `_Encoder(compression bool, obj ` + squareBrackets + rootClassName.Name + `) []byte{

			bb := parsers.Buffer{
				FloatIntEncoderVal: 10000.0,
				Endian: "big",
			}
	`

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			if i == 0 {
				*stringDataEncoder += `
			bb.PutShort(len(obj))
	`
				*stringDataEncoder += `
				for i` + strconv.Itoa(i) + `:=0;i` + strconv.Itoa(i) + `<len(obj);i` + strconv.Itoa(i) + `++{
	`
			} else {
				*stringDataEncoder += `
			bb.PutShort(len(obj` + totalParentBraces + `))
	`
				*stringDataEncoder += `
				for i` + strconv.Itoa(i) + `:=0;i` + strconv.Itoa(i) + `<len(obj` + totalParentBraces + `);i` + strconv.Itoa(i) + `++{
	`
			}

			totalParentBraces += "[i" + strconv.Itoa(i) + "]"
		}

	}

	totalParentBraces = "obj" + totalParentBraces + "."

	GenerateGolangEncodeCode(&currentIterate, stringDataEncoder, treeNode, totalParentBraces)

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			*stringDataEncoder += `
			}
				`
		}
	}

	*stringDataEncoder += `

			var compressedData []byte

			if compression{
				compressedData =  snappy.Encode(nil, bb.Array())
			}else{
				compressedData = bb.Array()
			}

			return compressedData
		}
	`
}

func DecoderCodeGeneration(rootClassName RootClass, stringDataDecoder *string, packageName *string, fileName *string, treeNode *TreeNode) {
	squareBrackets := ""

	for i := 0; i < rootClassName.ArrayCount; i++ {
		squareBrackets += "[]"
	}

	totalParentBraces := ""

	currentIterate := 0
	rootArrayClass := ""

	if squareBrackets != "" {
		rootArrayClass = `
			arrLen := bb.GetShort()
		`
		rootArrayClass += "obj := make(" + squareBrackets + rootClassName.Name + ", arrLen)"
	} else {
		rootArrayClass = "obj := " + rootClassName.Name + "{}"
	}

	*stringDataDecoder = `

	package ` + *packageName + `

	import(
		"log"
		"github.com/golang/snappy"
		"github.com/bytebuffer_parser/parsers"
	)

	func ` + strings.ToUpper(*fileName) + `_Decoder(compression bool, byteArr []byte) ` + squareBrackets + rootClassName.Name + `{

		bb := parsers.Buffer{
			FloatIntEncoderVal: 10000.0,
			Endian: "big",
		}

		if compression{
			decompressedData, err := snappy.Decode(nil, byteArr)

			if err != nil {
				log.Fatal("Failed to decompress data...")
			}

			bb.Wrap(decompressedData)
		}else{
			bb.Wrap(byteArr)
		}

		` + rootArrayClass + `
`
	if squareBrackets != "" {
		innerbracesCount := rootClassName.ArrayCount - 1
		for i := 0; i < rootClassName.ArrayCount; i++ {
			if i == 0 {
				*stringDataDecoder += `
		for i` + strconv.Itoa(i) + `:=0;i` + strconv.Itoa(i) + `<arrLen;i` + strconv.Itoa(i) + `++{
`
			} else {

				nestedSquareBrackets := ""

				for j := 0; j < innerbracesCount; j++ {
					nestedSquareBrackets += "[]"
				}

				innerbracesCount -= 1

				*stringDataDecoder += `

				arrLen` + strconv.Itoa(i) + ` := bb.GetShort()

				obj` + totalParentBraces + ` = make(` + nestedSquareBrackets + rootClassName.Name + `, arrLen` + strconv.Itoa(i) + `)
`
				*stringDataDecoder += `
		for i` + strconv.Itoa(i) + `:=0;i` + strconv.Itoa(i) + `<arrLen` + strconv.Itoa(i) + `;i` + strconv.Itoa(i) + `++{
`
			}

			totalParentBraces += "[i" + strconv.Itoa(i) + "]"
		}

	}

	totalParentBraces = "obj" + totalParentBraces + "."

	GenerateGolangDecoderCode(&currentIterate, stringDataDecoder, treeNode, totalParentBraces)

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			*stringDataDecoder += `
		}
			`
		}
	}

	*stringDataDecoder += `
		return obj
	}
`
}
