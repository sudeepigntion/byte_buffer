package parsers

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func generateCSharpClass(className string, properties []string) {
	fmt.Printf("public class %s {\n", className)
	for _, prop := range properties {
		fmt.Printf("    %s\n", prop)
	}
	fmt.Println("}")
}

func GenerateCSharpClass(classDefinitions string, rootClassName RootClass, globalMap *map[string][]string) (string, RootClass) {

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

				fieldType = strings.ReplaceAll(fieldType, "long", "long")
				fieldType = strings.ReplaceAll(fieldType, "short", "short")
				fieldType = strings.ReplaceAll(fieldType, "float", "double")

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
		goCode += fmt.Sprintf("\t public class %s {\n", structName)
		for _, field := range fields {
			count := strings.Count(field.Type, "[]")
			typeStr := ""
			for i := 0; i < count; i++ {
				typeStr += "[]"
			}
			field.Type = strings.ReplaceAll(field.Type, "[]", "")
			field.Type = field.Type + typeStr
			goCode += fmt.Sprintf("\t\t public %s %s { get; set; } \n", field.Type, field.Name)
		}
		goCode += "\t}\n\n"
	}

	return goCode, rootClassName
}

func WriteCSharpClassData(modelFileName string, finalStruct string) {

	// creating model file for example it will contain struct or class file
	file, err := os.Create(modelFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Close the file when we're done

	_, err = file.WriteString(finalStruct)
	if err != nil {
		log.Fatal(err)
	}
}

func EncoderCSharpCodeGeneration(rootClassName RootClass, stringDataEncoder *string, packageName *string, fileName *string, treeNode *TreeNode) {

	currentIterate := 0

	squareBrackets := ""
	for i := 0; i < rootClassName.ArrayCount; i++ {
		squareBrackets += "[]"
	}

	totalParentBraces := ""
	*stringDataEncoder = `
using System;
using System.Collections.Generic;
using System.IO;
using System.Text;

namespace ` + *packageName + ` { 

	class Encoder{

		public byte[] ` + *fileName + `_Encoder(` + squareBrackets + rootClassName.Name + ` obj){

			ByteBuffer bb = new ByteBuffer();
`

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			if i == 0 {
				*stringDataEncoder += `
			bb.PutShort(Convert.ToInt16(obj.Length));
	`
				*stringDataEncoder += `
			for (int i` + strconv.Itoa(i) + `=0;i` + strconv.Itoa(i) + `<obj.Length;i` + strconv.Itoa(i) + `++){
	`
			} else {
				*stringDataEncoder += "bb.PutShort(Convert.ToInt16(obj" + totalParentBraces + ".Length));"

				*stringDataEncoder += `
				for (int i` + strconv.Itoa(i) + `=0;i` + strconv.Itoa(i) + `<obj` + totalParentBraces + `.Length;i` + strconv.Itoa(i) + `++){
	`
			}

			totalParentBraces += "[i" + strconv.Itoa(i) + "]"
		}
	}

	totalParentBraces = "obj" + totalParentBraces + "."

	GenerateCSharpEncodeCode(&currentIterate, stringDataEncoder, treeNode, totalParentBraces)

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			*stringDataEncoder += `
			}
				`
		}
	}

	*stringDataEncoder += `

			byte[] response = bb.ToArray();

			bb.Dispose();
			
			return response;
		}
	}
}
	`
}

func GenerateCSharpEncodeCode(currentIterate *int, stringData *string, node *TreeNode, parentName string) {

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
				    bb.PutShort(Convert.ToInt16(` + path + `.Length));
				`
						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + `.Length;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(Convert.ToInt16(` + path + loopSquares + `.Length));
						`
						*stringData += `
								for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + loopSquares + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutInt(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutInt(` + path + `);
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
				    bb.PutShort(Convert.ToInt16(` + path + `.Length));
				`
						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + `.Length;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(Convert.ToInt16(` + path + loopSquares + `.Length));
						`
						*stringData += `
								for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + loopSquares + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutLong(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutLong(` + path + `);
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
				    bb.PutShort(Convert.ToInt16(` + path + `.Length));
				`
						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + `.Length;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(Convert.ToInt16(` + path + loopSquares + `.Length));
						`
						*stringData += `
								for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + loopSquares + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutShort(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutShort(` + path + `);
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
				    bb.PutShort(Convert.ToInt16(` + path + `.Length));
				`
						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + `.Length;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(Convert.ToInt16(` + path + loopSquares + `.Length));
						`
						*stringData += `
								for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + loopSquares + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutString(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutString(` + path + `);
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
				    bb.PutShort(Convert.ToInt16(` + path + `.Length));
				`
						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + `.Length;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(Convert.ToInt16(` + path + loopSquares + `.Length));
						`
						*stringData += `
								for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + loopSquares + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutDouble(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutDouble(` + path + `);
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
				    bb.PutShort(Convert.ToInt16(` + path + `.Length));
				`
						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + `.Length;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.PutShort(Convert.ToInt16(` + path + loopSquares + `.Length));
						`
						*stringData += `
								for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + loopSquares + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
							`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.PutBoolean(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.PutBoolean(` + path + `);
			`
			}
		default:
			if child.ArrayCount > 0 {
				squares := ""
				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				    bb.PutShort(Convert.ToInt16(` + path + `.Length));
				`
					*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` =0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + path + `.Length;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`
					squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

					path += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					*currentIterate += 1
				}
			}

			GenerateCSharpEncodeCode(currentIterate, stringData, child, path+".")
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

func DecoderCSharpCodeGeneration(rootClassName RootClass, stringDataDecoder *string, packageName *string, fileName *string, treeNode *TreeNode) {
	squareBrackets := ""
	firstSquareBrackets := ""

	for i := 0; i < rootClassName.ArrayCount; i++ {
		squareBrackets += "[]"
		if i == 0 {
			firstSquareBrackets += "[arrLen]"
		} else {
			firstSquareBrackets += "[]"
		}
	}

	totalParentBraces := ""

	currentIterate := 0
	rootArrayClass := ""

	if squareBrackets != "" {
		rootArrayClass = `
			int arrLen = (int)bb.GetShort();
		`
		rootArrayClass += rootClassName.Name + squareBrackets + " obj = new " + rootClassName.Name + firstSquareBrackets + ";"
	} else {
		rootArrayClass = rootClassName.Name + " obj = new " + rootClassName.Name + "();"
	}

	*stringDataDecoder = `
using System;
using System.Collections.Generic;
using System.IO;
using System.Text;

namespace ` + *packageName + ` { 

	class Decoder{

		public ` + squareBrackets + rootClassName.Name + ` ` + *fileName + `_Decoder(byte[] byteArr){

			ByteBuffer bb = new ByteBuffer();

			bb.Wrap(byteArr);

		` + rootArrayClass + `
`
	if squareBrackets != "" {
		innerbracesCount := rootClassName.ArrayCount - 1
		for i := 0; i < rootClassName.ArrayCount; i++ {
			if i == 0 {
				*stringDataDecoder += `
		for (int i` + strconv.Itoa(i) + `=0;i` + strconv.Itoa(i) + `<arrLen;i` + strconv.Itoa(i) + `++){
`
			} else {

				nestedSquareBrackets := ""

				for j := 0; j < innerbracesCount; j++ {
					if j == 0 {
						nestedSquareBrackets += "[arrLen" + strconv.Itoa(i) + "]"
					} else {
						nestedSquareBrackets += "[]"
					}
				}

				innerbracesCount -= 1

				*stringDataDecoder += `
					int arrLen` + strconv.Itoa(i) + ` = (int)bb.GetShort();
				`
				*stringDataDecoder += ` obj` + totalParentBraces + ` = new ` + rootClassName.Name + nestedSquareBrackets + `;`
				*stringDataDecoder += `
		for (int i` + strconv.Itoa(i) + `=0;i` + strconv.Itoa(i) + `<arrLen` + strconv.Itoa(i) + `;i` + strconv.Itoa(i) + `++){
`
			}

			totalParentBraces += "[i" + strconv.Itoa(i) + "]"
		}

	}

	totalParentBraces = "obj" + totalParentBraces + "."

	GenerateCSharpDecoderCode(&currentIterate, stringDataDecoder, treeNode, totalParentBraces)

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			*stringDataDecoder += `
		}
			`
		}
	}

	*stringDataDecoder += `

		bb.Dispose();

		return obj;
		}
	}
}
`
}

func GenerateCSharpDecoderCode(currentIterate *int, stringData *string, node *TreeNode, parentName string) {

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
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = (int)bb.GetShort();
				`
						*stringData += child.Value + arrayBraces + " " + path + squares + `= new int` + firstArrBraces + `;`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `
						int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := (int)bb.GetShort();
						`

						*stringData += child.Value + arrayBraces + " " + path + squares + `= new int` + firstArrBraces + `;`
						// *stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetInt();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.GetInt();
			`
			}
		case "long":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = (int)bb.GetShort();
				`
						*stringData += child.Value + arrayBraces + " " + path + squares + `= new long` + firstArrBraces + `;`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `
						int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := (int)bb.GetShort();
						`

						*stringData += child.Value + arrayBraces + " " + path + squares + `= new long` + firstArrBraces + `;`
						// *stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetLong();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.GetLong();
			`
			}
		case "short":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = (int)bb.GetShort();
				`
						*stringData += path + squares + `= new short` + firstArrBraces + `;`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `
						int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := (int)bb.GetShort();
						`

						*stringData += path + squares + `= new short` + firstArrBraces + `;`
						// *stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetShort();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.GetShort();
			`
			}
		case "string":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = (int)bb.GetShort();
				`
						*stringData += path + squares + `= new string` + firstArrBraces + `;`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `
						int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := (int)bb.GetShort();
						`

						*stringData += path + squares + `= new string` + firstArrBraces + `;`
						// *stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetString();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.GetString();
			`
			}
		case "float":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = (int)bb.GetShort();
				`
						// child.Value + arrayBraces + " " +
						*stringData += path + squares + `= new double` + firstArrBraces + `;`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `
						int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := (int)bb.GetShort();
						`

						*stringData += path + squares + `= new double` + firstArrBraces + `;`
						// *stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetDouble();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.GetDouble();
			`
			}
		case "bool":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = (int)bb.GetShort();
				`
						*stringData += path + squares + `= new bool` + firstArrBraces + `;`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""
						firstArrBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "[]"

							if j == 0 {
								firstArrBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
							} else {
								firstArrBraces += "[]"
							}
						}

						*stringData += `
						int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` := (int)bb.GetShort();
						`

						*stringData += path + squares + `= new bool` + firstArrBraces + `;`
						// *stringData += path + squares + `= make(` + arrayBraces + `int, ` + child.Name + `ArrLen` + strconv.Itoa(i) + `)`

						*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.GetBool();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.GetBool();
			`
			}
		default:

			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {

					arrayBraces := ""

					for j := 0; j < arrayCount; j++ {
						if j == 0 {
							arrayBraces += "[" + child.Name + "ArrLen" + strconv.Itoa(i) + "]"
						} else {
							arrayBraces += "[]"
						}
					}

					*stringData += `
					int ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = (int)bb.GetShort();
				`

					*stringData += path + `= new ` + child.Value + arrayBraces + `;`

					*stringData += `
							for (int index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `=0;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `<` + child.Name + `ArrLen` + strconv.Itoa(i) + `;index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `++){
						`

					squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

					path += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					*currentIterate += 1
					arrayCount -= 1
				}
			}

			GenerateCSharpDecoderCode(currentIterate, stringData, child, path+".")
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

func WriteCsharpDecoderData(decoderFileName string, stringDataDecoder string) {

	// creating encoder file
	// Create the file (or truncate it if it already exists)
	file, err := os.Create(decoderFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Close the file when we're done

	_, err = file.WriteString(stringDataDecoder)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteCsharpEncoderData(encoderFileName string, stringDataEncoder string) {
	// creating encoder file
	// Create the file (or truncate it if it already exists)
	file, err := os.Create(encoderFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Close the file when we're done

	_, err = file.WriteString(stringDataEncoder)
	if err != nil {
		log.Fatal(err)
	}
}
