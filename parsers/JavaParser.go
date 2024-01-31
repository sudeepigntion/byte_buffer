package parsers

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type JavaParser struct{}

func (p *JavaParser) GenerateJavaClass(packageName *string, classDefinitions string, rootClassName RootClass, globalMap *map[string][]string) (string, RootClass) {
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
				fieldType = strings.ReplaceAll(fieldType, "float", "float")
				fieldType = strings.ReplaceAll(fieldType, "string", "String")
				fieldType = strings.ReplaceAll(fieldType, "bool", "boolean")

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
		goCode += fmt.Sprintf("--split--package %s;\n\npublic class %s {\n", *packageName, structName)
		getterAndSetter := ""

		for _, field := range fields {
			count := strings.Count(field.Type, "[]")
			typeStr := ""

			for i := 0; i < count; i++ {
				typeStr += "[]"
			}

			field.Type = strings.ReplaceAll(field.Type, "[]", "")
			field.Type = field.Type + typeStr
			fieldName := strings.ToLower(field.Name)
			goCode += fmt.Sprintf("\t\tprivate %s %s;\n", field.Type, fieldName)

			// creating setter
			getterAndSetter += fmt.Sprintf("\t\tpublic void set%s(%s %s)\n\t\t{\n\t\t\tthis.%s = %s;\n\t\t}\n\n", field.Name, field.Type, fieldName, fieldName, fieldName)

			// creating getter
			getterAndSetter += fmt.Sprintf("\t\tpublic %s get%s()\n\t\t{\n\t\t\treturn %s;\n\t\t}\n\n", field.Type, field.Name, fieldName)
		}

		goCode += "\n" + getterAndSetter
		goCode = strings.TrimRight(goCode, "\t")
		goCode += "\t}\n\n"
	}
	return goCode, rootClassName
}

func (p *JavaParser) EncoderJavaCodeGeneration(rootClassName RootClass, stringDataEncoder *string, packageName *string, fileName *string, treeNode *TreeNode) {
	currentIterate := 0

	squareBrackets := ""
	for i := 0; i < rootClassName.ArrayCount; i++ {
		squareBrackets += "[]"
	}

	totalParentBraces := ""

	*stringDataEncoder = fmt.Sprintf("package %s;\n\n", *packageName)
	// import statement
	*stringDataEncoder += fmt.Sprintf("import java.io.IOException;\n\n")

	// class declaration
	*stringDataEncoder += fmt.Sprintf("public class Encoder {\n")

	// functions
	*stringDataEncoder += fmt.Sprintf("\tpublic byte[] %s_Encoder(%s obj) throws IOException {\n", strings.ToUpper(*fileName), rootClassName.Name+squareBrackets)
	*stringDataEncoder += fmt.Sprintf("\t\tJavaBuffer bb = new JavaBuffer();\n\n")

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			if i == 0 {
				*stringDataEncoder += fmt.Sprintf("\t\tbb.putShort((short)obj.length);\n\n")
				*stringDataEncoder += fmt.Sprintf("\t\tfor(int i%d = 0; i%d < obj.length; i%d++) {\n", i, i, i)
			} else {
				*stringDataEncoder += fmt.Sprintf("%sbb.putShort((short)obj%s.length);\n\n", strings.Repeat("\t", i+2), totalParentBraces)
				*stringDataEncoder += fmt.Sprintf("%sfor (int i%d = 0; i%d < obj%s.length; i%d++) {\n", strings.Repeat("\t", i+2), i, i, totalParentBraces, i)
			}

			totalParentBraces += "[i" + strconv.Itoa(i) + "]"
		}

	}

	totalParentBraces = "obj" + totalParentBraces + "."

	GenerateJavaEncodeCode(&currentIterate, stringDataEncoder, treeNode, totalParentBraces)

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			*stringDataEncoder += `
}
`
		}
	}

	*stringDataEncoder += `

byte[] response = bb.toArray();
bb.close();
return response;
}
}
`
}

func GenerateJavaEncodeCode(currentIterate *int, stringData *string, node *TreeNode, parentName string) {

	for _, child := range node.Children {

		path := parentName + "get" + child.Name + "()"

		switch child.Value {
		case "int":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s.length);\n\n", strings.Repeat("\t", 2+i), path)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s.length; index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares

					} else {
						dec += 1
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s%s.length);\n\n", strings.Repeat("\t", 2+i), path, loopSquares)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s%s.length; index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, loopSquares, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					}
				}

				*stringData += fmt.Sprintf("%sbb.putInt(%s%s);\n", strings.Repeat("\t", child.ArrayCount+2), path, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("%s}\n", strings.Repeat("\t", child.ArrayCount-i+2))
				}

			} else {
				*stringData += fmt.Sprintf("%sbb.putInt(%s);\n", strings.Repeat("\t", 4), path)
			}
		case "long":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s.length);\n\n", strings.Repeat("\t", 2+i), path)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s.length; index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					} else {
						dec += 1
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s%s.length);\n\n", strings.Repeat("\t", 2+i), path, loopSquares)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s%s.length; index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, loopSquares, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					}
				}
				*stringData += fmt.Sprintf("%sbb.putLong(%s%s);\n", strings.Repeat("\t", child.ArrayCount+2), path, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("%s}\n", strings.Repeat("\t", child.ArrayCount-i))
				}
			} else {
				*stringData += fmt.Sprintf("%sbb.putLong(%s);\n", strings.Repeat("\t", 2), path)
			}
		case "short":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s.length);\n\n", strings.Repeat("\t", 2+i), path)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s.length; index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					} else {
						dec += 1
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s%s.length);\n\n", strings.Repeat("\t", 2+i), path, loopSquares)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s%s.length; index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, loopSquares, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					}
				}

				*stringData += fmt.Sprintf("%sbb.putShort(%s%s);\n", strings.Repeat("\t", child.ArrayCount+2), path, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("%s}\n", strings.Repeat("\t", child.ArrayCount-i))
				}
			} else {
				*stringData += fmt.Sprintf("%sbb.putShort(%s);\n", strings.Repeat("\t", 2), path)
			}
		case "string":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s.length);\n\n", strings.Repeat("\t", 2+i), path)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s.length; index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					} else {
						dec += 1
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s%s.length);\n\n", strings.Repeat("\t", 2+i), path, loopSquares)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s%s.length; index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, loopSquares, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					}
				}

				*stringData += fmt.Sprintf("%sbb.putString(%s%s);\n", strings.Repeat("\t", child.ArrayCount+2), path, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("%s}\n", strings.Repeat("\t", child.ArrayCount-i))
				}
			} else {
				*stringData += fmt.Sprintf("%sbb.putString(%s);\n", strings.Repeat("\t", 4), path)
			}
		case "float":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s.length);\n\n", strings.Repeat("\t", 2+i), path)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s.length; index%d%d++) {\n", strings.Repeat("\t", 4+i), *currentIterate, i, *currentIterate, i, path, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					} else {
						dec += 1
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s%s.length);\n\n", strings.Repeat("\t", 2+i), path, loopSquares)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s%s.length; index%d%d++) {\n", strings.Repeat("\t", 4+i), *currentIterate, i, *currentIterate, i, path, loopSquares, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					}
				}

				*stringData += fmt.Sprintf("%sbb.putFloatUsingIntEncoding(%s%s);\n", strings.Repeat("\t", child.ArrayCount+2), path, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("%s}\n", strings.Repeat("\t", child.ArrayCount-i))
				}
			} else {
				*stringData += fmt.Sprintf("%sbb.putFloatUsingIntEncoding(%s);\n", strings.Repeat("\t", 2), path)
			}
		case "bool":
			if child.ArrayCount > 0 {
				squares := ""
				loopSquares := ""
				dec := 0
				for i := 0; i < child.ArrayCount; i++ {
					if i == 0 {
						*stringData += fmt.Sprintf("%sbb.putShort((short)%s.length);\n\n", strings.Repeat("\t", 2+i), path)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s.length; index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					} else {
						dec += 1

						*stringData += fmt.Sprintf("%sbb.putShort((short)%s%s.length);\n\n", strings.Repeat("\t", 2+i), path, loopSquares)

						*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s%s.length; index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, loopSquares, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

						loopSquares = squares
					}
				}

				*stringData += fmt.Sprintf("%sbb.putBoolean(%s%s);\n", strings.Repeat("\t", child.ArrayCount+2), path, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("%s}\n", strings.Repeat("\t", child.ArrayCount-i))
				}
			} else {
				*stringData += fmt.Sprintf("%sbb.putBoolean(%s);\n", strings.Repeat("\t", 2), path)
			}
		default:
			if child.ArrayCount > 0 {
				squares := ""
				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("%sbb.putShort((short)%s.length);\n", strings.Repeat("\t", 2+i), path)

					*stringData += fmt.Sprintf("%sfor (int index%d%d = 0; index%d%d < %s.length;index%d%d++) {\n", strings.Repeat("\t", 2+i), *currentIterate, i, *currentIterate, i, path, *currentIterate, i)

					squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

					path += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

					*currentIterate += 1
				}
			}

			GenerateJavaEncodeCode(currentIterate, stringData, child, path+".")
		}

	}

	if node.ArrayCount > 0 {
		for i := 0; i < node.ArrayCount; i++ {
			*stringData += fmt.Sprintf("\n}")
		}
	}
}

func (p *JavaParser) WriteJavaModalClass(packageName string, modelClasses string) {
	tempModelClasses := strings.Split(modelClasses, "--split--")

	for _, modelClass := range tempModelClasses {
		modelClass = strings.TrimSpace(modelClass)

		if modelClass == "" {
			continue
		}

		pattern := `class (\w+)`

		// Compile the regular expression pattern
		regex := regexp.MustCompile(pattern)

		// Find the class name using the regular expression
		matches := regex.FindStringSubmatch(modelClass)

		className := ""
		// Check if a match was found
		if len(matches) > 1 {
			className = matches[1]
		}

		// creating model file for example it will contain struct or class file
		dirPath := "./src/" + packageName

		_, err := os.Stat(dirPath)

		if os.IsNotExist(err) {
			// Directory does not exist, so create it
			err := os.MkdirAll(dirPath, 0755) // 0755 is a common directory permission
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
		}

		filePath := dirPath + "/" + className + ".java"
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close() // Close the file when we're done

		_, err = file.WriteString(modelClass)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func (p *JavaParser) WriteJavaEncoderData(dirPath string, encoderFileName string, stringDataEncoder string) {

	// creating encoder file
	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {
		// Directory does not exist, so create it
		err := os.MkdirAll(dirPath, 0755) // 0755 is a common directory permission
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}
	// Create the file (or truncate it if it already exists)
	file, err := os.Create(dirPath + "/" + encoderFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Close the file when we're done

	_, err = file.WriteString(stringDataEncoder)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *JavaParser) DecoderJavaCodeGeneration(rootClassName RootClass, stringDataDecoder *string, packageName *string, fileName *string, treeNode *TreeNode) {
	squareBrackets := ""
	firstSquareBrackets := ""
	parentArrObj := ""

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
			int arrLen = bb.getShort();
		`
		rootArrayClass += rootClassName.Name + squareBrackets + " obj = new " + rootClassName.Name + firstSquareBrackets + ";"
		parentArrObj = fmt.Sprintf(" = new %s();", rootClassName.Name)
	} else {
		rootArrayClass = rootClassName.Name + " obj = new " + rootClassName.Name + "();"
	}

	// writing code
	*stringDataDecoder = fmt.Sprintf("package %s;\n", *packageName)
	*stringDataDecoder += fmt.Sprintf("import java.io.IOException;\n")
	*stringDataDecoder += fmt.Sprintf("\npublic class Decoder {\n")

	*stringDataDecoder += fmt.Sprintf("\tpublic %s%s %s_Decoder(byte[] byteArr) throws IOException {\n", rootClassName.Name, squareBrackets, *fileName)
	*stringDataDecoder += fmt.Sprintf("\t\tJavaBuffer bb = new JavaBuffer();\n")
	*stringDataDecoder += fmt.Sprintf("\t\tbb.wrap(byteArr);\n")
	*stringDataDecoder += fmt.Sprintf("\t\t%s\n\n", rootArrayClass)
	if squareBrackets != "" {
		innerBracesCount := rootClassName.ArrayCount - 1

		for i := 0; i < rootClassName.ArrayCount; i++ {
			if i == 0 {
				*stringDataDecoder += fmt.Sprintf("\t\tfor(int i%d = 0; i%d < arrLen; i%d++) {\n", i, i, i)
			} else {
				nestedSquareBrackets := ""

				for j := 0; j < innerBracesCount; j++ {
					if j == 0 {
						nestedSquareBrackets += "[arrLen" + strconv.Itoa(i) + "]"
					} else {
						nestedSquareBrackets += "[]"
					}
				}

				innerBracesCount -= 1

				*stringDataDecoder += fmt.Sprintf("\t\tint arrLen%d = (int)bb.getShort();\n", i)
				*stringDataDecoder += fmt.Sprintf("\t\tobj%s = new %s%s;", totalParentBraces, rootClassName.Name, nestedSquareBrackets)
				*stringDataDecoder += fmt.Sprintf("\n\t\tfor(int i%d = 0; i%d < arrLen%d; i%d++) {", i, i, i, i)
			}

			totalParentBraces += "[i" + strconv.Itoa(i) + "]"
		}
	}

	totalParentBraces = "obj" + totalParentBraces + "."

	GenerateJavaDecoderCode(&currentIterate, stringDataDecoder, treeNode, totalParentBraces, parentArrObj)

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			*stringDataDecoder += fmt.Sprintf("\n}\n")
		}
	}
	*stringDataDecoder += fmt.Sprintf("\n\t\tbb.close();\n")
	*stringDataDecoder += fmt.Sprintf("\n\t\treturn obj;\n")
	// closing method bracket
	*stringDataDecoder += fmt.Sprintf("\t}\n")
	// closing class bracket
	*stringDataDecoder += fmt.Sprintf("}\n")
}

func GenerateJavaDecoderCode(currentIterate *int, stringData *string, node *TreeNode, parentName string, parentArrObj string) {
	if parentArrObj != "" {
		parentArrObj = fmt.Sprintf("\n\t\t%s%s", strings.TrimRight(parentName, "."), parentArrObj)
	}

	arrClassObjName := ""
	setObjToClass := ""

	for _, child := range node.Children {

		switch child.Value {
		case "int":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				tempArrObjName := ""
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

						*stringData += fmt.Sprintf("\n\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)
						tempArrObjName = fmt.Sprintf("%s%d", strings.ToLower(child.Name), i)
						*stringData += fmt.Sprintf("\t\t%s%s %s%s = new int%s;\n", child.Value, arrayBraces, tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

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

						*stringData += fmt.Sprintf("\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)

						*stringData += fmt.Sprintf("\t\t%s%s = new int %s;\n", tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += fmt.Sprintf("\n\t\t%s%s = bb.getInt();\n", tempArrObjName, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("\t}\n")
				}

				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(%s);\n", parentName, child.Name, tempArrObjName)
			} else {
				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(bb.getInt());\n", parentName, child.Name)
			}
		case "long":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				tempArrObjName := ""
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

						*stringData += fmt.Sprintf("\n\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)
						tempArrObjName = fmt.Sprintf("%s%d", strings.ToLower(child.Name), i)
						*stringData += fmt.Sprintf("\t\t%s%s %s%s = new long%s;\n", child.Value, arrayBraces, tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

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

						*stringData += fmt.Sprintf("\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)

						*stringData += fmt.Sprintf("\t\t%s%s = new long%s;\n", tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += fmt.Sprintf("\n\t\t%s%s = bb.getLong();\n", tempArrObjName, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("\t}\n")
				}

				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(%s);\n", parentName, child.Name, tempArrObjName)
			} else {
				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(bb.getLong());\n", parentName, child.Name)
			}
		case "short":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				tempArrObjName := ""
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

						*stringData += fmt.Sprintf("\n\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)
						tempArrObjName = fmt.Sprintf("%s%d", strings.ToLower(child.Name), i)
						*stringData += fmt.Sprintf("\t\t%s%s %s%s = new short%s;\n", child.Value, arrayBraces, tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

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

						*stringData += fmt.Sprintf("\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)

						*stringData += fmt.Sprintf("\t\t%s%s = new short%s;\n", tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += fmt.Sprintf("\n\t\t%s%s = bb.getShort();\n", tempArrObjName, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("\t}\n")
				}

				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(%s);\n", parentName, child.Name, tempArrObjName)
			} else {
				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(bb.getShort());\n", parentName, child.Name)
			}
		case "string":
			child.Value = "String"
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				tempArrObjName := ""
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

						*stringData += fmt.Sprintf("\n\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)
						tempArrObjName = fmt.Sprintf("%s%d", strings.ToLower(child.Name), i)
						*stringData += fmt.Sprintf("\t\t%s%s %s%s = new String%s;\n", child.Value, arrayBraces, tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

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

						*stringData += fmt.Sprintf("\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)

						*stringData += fmt.Sprintf("\t\t%s%s = new String%s;\n", tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += fmt.Sprintf("\n\t\t%s%s = bb.getString();\n", tempArrObjName, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("\t}\n")
				}

				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(%s);\n", parentName, child.Name, tempArrObjName)
			} else {
				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(bb.getString());\n", parentName, child.Name)
			}
		case "float":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				tempArrObjName := ""
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

						*stringData += fmt.Sprintf("\n\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)
						tempArrObjName = fmt.Sprintf("%s%d", strings.ToLower(child.Name), i)
						*stringData += fmt.Sprintf("\t\t%s%s %s%s = new float%s;\n", child.Value, arrayBraces, tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

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

						*stringData += fmt.Sprintf("\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)

						*stringData += fmt.Sprintf("\t\t%s%s = new float%s;\n", tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += fmt.Sprintf("\n\t\t%s%s = bb.getFloatUsingIntEncoding();\n", tempArrObjName, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("\t}\n")
				}

				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(%s);\n", parentName, child.Name, tempArrObjName)
			} else {
				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(bb.getFloatUsingIntEncoding());\n", parentName, child.Name)
			}
		case "bool":
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				tempArrObjName := ""
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

						*stringData += fmt.Sprintf("\n\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)
						tempArrObjName = fmt.Sprintf("%s%d", strings.ToLower(child.Name), i)
						*stringData += fmt.Sprintf("\t\t%s%s %s%s = new boolean%s;\n", child.Value, arrayBraces, tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

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

						*stringData += fmt.Sprintf("\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)

						*stringData += fmt.Sprintf("\t\t%s%s = new boolean%s;\n", tempArrObjName, squares, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int index%d%d = 0; index%d%d < %sArrLen%d; index%d%d++) {\n", *currentIterate, i, *currentIterate, i, child.Name, i, *currentIterate, i)

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += fmt.Sprintf("\n\t\t%s%s = bb.getBool();\n", tempArrObjName, squares)

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += fmt.Sprintf("\t}\n")
				}

				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(%s);\n", parentName, child.Name, tempArrObjName)
			} else {
				*stringData += parentArrObj
				parentArrObj = ""
				*stringData += fmt.Sprintf("\n\t\t%sset%s(bb.getBool());\n", parentName, child.Name)
			}
		default:
			classObjName := ""
			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				tempArrObjName := ""
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

						*stringData += fmt.Sprintf("\n\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)
						tempArrObjName = fmt.Sprintf("%sArr%d", strings.ToLower(child.Name), i)
						*stringData += fmt.Sprintf("\t\t%s%s %s%s = new %s%s;\n", child.Value, arrayBraces, tempArrObjName, squares, child.Value, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int %sIndex%d%d = 0; %sIndex%d%d < %sArrLen%d; %sIndex%d%d++) {\n", strings.ToLower(child.Name), *currentIterate, i, strings.ToLower(child.Name), *currentIterate, i, child.Name, i, strings.ToLower(child.Name), *currentIterate, i)

						squares += `[` + strings.ToLower(child.Name) + `Index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
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

						*stringData += fmt.Sprintf("\t\tint %sArrLen%d = (int)bb.getShort();\n", child.Name, i)

						*stringData += fmt.Sprintf("\t\t%s%s = new %s%s;\n", tempArrObjName, squares, child.Value, firstArrBraces)

						*stringData += fmt.Sprintf("\n\t\tfor(int %sIndex%d%d = 0; %sIndex%d%d < %sArrLen%d; %sIndex%d%d++) {\n", strings.ToLower(child.Name), *currentIterate, i, strings.ToLower(child.Name), *currentIterate, i, child.Name, i, strings.ToLower(child.Name), *currentIterate, i)

						squares += `[` + strings.ToLower(child.Name) + `Index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				classObjName = fmt.Sprintf("%sObj", strings.ToLower(child.Value))
				*stringData += fmt.Sprintf("\t\t%s %s = new %s();\n", child.Value, classObjName, child.Value)

				arrClassObjName += fmt.Sprintf("\n\t\t%s%s = %s;\n", tempArrObjName, squares, classObjName)

				setObjToClass += fmt.Sprintf("\n\t\t%sset%s(%s);\n", parentName, child.Name, tempArrObjName)
			} else {
				classObjName = fmt.Sprintf("%sObj", strings.ToLower(child.Name))
				*stringData += fmt.Sprintf("\n\t\t%s %s = new %s();\n", child.Value, classObjName, child.Value)

				setObjToClass += fmt.Sprintf("\n\t\t%sset%s(%s);\n", parentName, child.Name, classObjName)
			}

			GenerateJavaDecoderCode(currentIterate, stringData, child, classObjName+".", "")

			*stringData += arrClassObjName
			for i := 0; i < child.ArrayCount; i++ {
				*stringData += fmt.Sprintf("\t}\n")
			}
		}
	}

	if setObjToClass != "" {
		*stringData += setObjToClass
	}
}

func (p *JavaParser) WriteJavaDecoderData(dirPath string, encoderFileName string, stringDataEncoder string) {
	// creating encoder file
	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {
		// Directory does not exist, so create it
		err := os.MkdirAll(dirPath, 0755) // 0755 is a common directory permission
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}
	// Create the file (or truncate it if it already exists)
	file, err := os.Create(dirPath + "/" + encoderFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Close the file when we're done

	_, err = file.WriteString(stringDataEncoder)
	if err != nil {
		log.Fatal(err)
	}
}
