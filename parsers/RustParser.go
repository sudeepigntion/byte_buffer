package parsers

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GenerateRustStruct(classDefinitions string, rootClassName RootClass, globalMap *map[string][]string) (string, RootClass) {

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

				fieldType = strings.ReplaceAll(fieldType, "long", "i64")
				fieldType = strings.ReplaceAll(fieldType, "short", "i16")
				fieldType = strings.ReplaceAll(fieldType, "float", "f64")
				fieldType = strings.ReplaceAll(fieldType, "int", "i32")
				fieldType = strings.ReplaceAll(fieldType, "string", "String")

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
		goCode += fmt.Sprintf("\t #[derive(Clone, Debug, Default)] \n")
		goCode += fmt.Sprintf("\t pub struct %s {\n", structName)
		for _, field := range fields {

			count := strings.Count(field.Type, "[]")
			typeStr := ""
			for i := 0; i < count; i++ {
				typeStr += "Vec<"
			}

			typeStr += field.Type

			for i := 0; i < count; i++ {
				typeStr += ">"
			}

			if typeStr != "" {
				typeStr = strings.ReplaceAll(typeStr, "[]", "")
				goCode += fmt.Sprintf("\t\t pub %s: %s, \n", field.Name, typeStr)
			} else {
				field.Type = field.Type + typeStr
				goCode += fmt.Sprintf("\t\t pub %s: %s, \n", field.Name, field.Type)
			}
		}
		goCode += "\t }\n\n"
	}

	return goCode, rootClassName
}

func GenerateRustEncodeCode(currentIterate *int, stringData *string, node *TreeNode, parentName string) {

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
				    bb.put_short(` + path + `.len() as i16);
				`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.put_short(` + path + loopSquares + `.len() as i16);
						`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + loopSquares + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.put_int(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.put_int(` + path + `);
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
				    bb.put_short(` + path + `.len() as i16);
				`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.put_short(` + path + loopSquares + `.len() as i16);
						`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + loopSquares + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.put_long(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.put_long(` + path + `);
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
				    bb.put_short(` + path + `.len() as i16);
				`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.put_short(` + path + loopSquares + `.len() as i16);
						`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + loopSquares + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.put_short(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.put_short(` + path + `);
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
				    bb.put_short(` + path + `.len() as i16);
				`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.put_short(` + path + loopSquares + `.len() as i16);
						`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + loopSquares + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.put_string(` + path + squares + `.to_string());
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.put_string(` + path + `.to_string());
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
				    bb.put_short(` + path + `.len() as i16);
				`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.put_short(` + path + loopSquares + `.len() as i16);
						`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + loopSquares + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.put_float(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.put_float(` + path + `);
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
				    bb.put_short(` + path + `.len() as i16);
				`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					} else {
						dec += 1
						*stringData += `
							bb.put_short(` + path + loopSquares + `.len() as i16);
						`
						*stringData += `
							for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + loopSquares + `.iter().enumerate() {
						`
						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
						loopSquares = squares
					}
				}

				*stringData += `
				bb.put_bool(` + path + squares + `);
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				bb.put_bool(` + path + `);
			`
			}
		default:

			if child.ArrayCount > 0 {
				squares := ""
				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				    bb.put_short(` + path + `.len() as i16);
				`
					*stringData += `
						for (index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `, _) in ` + path + `.iter().enumerate() {
					`
					squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

					path += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					*currentIterate += 1
				}
			}

			GenerateRustEncodeCode(currentIterate, stringData, child, path+".")
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

func GenerateRustDecoderCode(importPackage *string, currentIterate *int, stringData *string, node *TreeNode, parentName string) {
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
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
				    let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
				`
						*stringData += path + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
						let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
						`

						if i == child.ArrayCount-1 {
							*stringData += path + squares + ` = vec![0; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						} else {
							*stringData += path + squares + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						}

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.get_int();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.get_int();
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
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
				    let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
				`
						*stringData += path + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
						let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
						`

						if i == child.ArrayCount-1 {
							*stringData += path + squares + ` = vec![0; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						} else {
							*stringData += path + squares + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						}

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.get_long();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.get_long();
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
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
				    let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
				`
						*stringData += path + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
						let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
						`

						if i == child.ArrayCount-1 {
							*stringData += path + squares + ` = vec![0; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						} else {
							*stringData += path + squares + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						}

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.get_short();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.get_short();
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
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
				    let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
				`
						*stringData += path + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
						let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
						`

						if i == child.ArrayCount-1 {
							*stringData += path + squares + ` = vec!["".to_string(); ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						} else {
							*stringData += path + squares + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						}

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.get_string();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.get_string();
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
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
				    let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
				`
						*stringData += path + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
						let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
						`

						if i == child.ArrayCount-1 {
							*stringData += path + squares + ` = vec![0.0; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						} else {
							*stringData += path + squares + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						}

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.get_float();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.get_float();
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
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
				    let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
				`
						*stringData += path + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					} else {

						arrayBraces := ""

						for j := 0; j < arrayCount; j++ {
							arrayBraces += "Vec<"
						}

						arrayBraces += child.Value

						for j := 0; j < arrayCount; j++ {
							arrayBraces += ">"
						}

						*stringData += `
						let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
						`

						if i == child.ArrayCount-1 {
							*stringData += path + squares + ` = vec![false; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						} else {
							*stringData += path + squares + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
						}

						*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`

						squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`
					}

					arrayCount -= 1
				}

				*stringData += `
				` + path + squares + ` = bb.get_bool();
				`

				for i := 0; i < child.ArrayCount; i++ {
					*stringData += `
				}
				`
				}
			} else {
				*stringData += `
				` + path + ` = bb.get_bool();
			`
			}
		default:

			*importPackage += child.Value + ","

			if child.ArrayCount > 0 {
				squares := ""
				arrayCount := child.ArrayCount
				for i := 0; i < child.ArrayCount; i++ {

					arrayBraces := ""

					for j := 0; j < arrayCount; j++ {
						arrayBraces += "Vec<"
					}

					arrayBraces += child.Value

					for j := 0; j < arrayCount; j++ {
						arrayBraces += ">"
					}

					*stringData += `
					let ` + child.Name + `ArrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
				`

					if i == child.ArrayCount-1 {
						*stringData += path + ` = vec![` + child.Value + `{..Default::default()}; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
					} else {
						*stringData += path + ` = vec![vec![]; ` + child.Name + `ArrLen` + strconv.Itoa(i) + `];`
					}

					*stringData += `
							for index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + ` in 0..` + child.Name + `ArrLen` + strconv.Itoa(i) + `{
						`
					squares += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

					path += `[index` + strconv.Itoa(*currentIterate) + strconv.Itoa(i) + `]`

					*stringData += path + ` = ` + child.Value + `{..Default::default()};`

					*currentIterate += 1
					arrayCount -= 1
				}
			}

			GenerateRustDecoderCode(importPackage, currentIterate, stringData, child, path+".")
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

func WriteRustEncoderData(encoderFileName string, stringDataEncoder string) {
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

	fmt.Println("Execute 'rustfmt " + encoderFileName + "' to format the code...")
}

func WriteRustDecoderData(decoderFileName string, stringDataDecoder string) {
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

	fmt.Println("Execute 'rustfmt " + decoderFileName + "' to format the code...")
}

func WriteRustStructData(modelFileName string, finalStruct string) {

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

	fmt.Println("Execute 'rustfmt " + modelFileName + "' to format the code...")
}

func RustEncoderCodeGeneration(rootClassName RootClass, stringDataEncoder *string, packageName *string, fileName *string, treeNode *TreeNode) {

	currentIterate := 0

	squareBrackets := ""
	for i := 0; i < rootClassName.ArrayCount; i++ {
		squareBrackets += "Vec<"
	}

	squareBrackets += rootClassName.Name

	for i := 0; i < rootClassName.ArrayCount; i++ {
		squareBrackets += ">"
	}

	totalParentBraces := ""
	*stringDataEncoder = `

		pub mod ` + *packageName + ` {

		use crate::` + *packageName + `::RustBuffer::ByteBuff;
		use crate::` + *packageName + `::` + *fileName + `::` + *packageName + `::{{Packages}};

		pub fn ` + *fileName + `_Encoder(obj: &` + squareBrackets + `) -> Vec<u8>{

			let mut bb = ByteBuff{
				multiplier: 10000.0,
				endian: "big".to_string(),
				..Default::default()
			};

			bb.init("big".to_string());
	`

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			if i == 0 {
				*stringDataEncoder += `
			bb.put_short(obj.len() as i16);
	`
				*stringDataEncoder += `
				for (i` + strconv.Itoa(i) + `, _) in obj.iter().enumerate() {
	`
			} else {
				*stringDataEncoder += `
			bb.put_short(obj` + totalParentBraces + `.len() as i16);
	`
				*stringDataEncoder += `
				for (i` + strconv.Itoa(i) + `, _) in obj` + totalParentBraces + `.iter().enumerate() {
	`
			}

			totalParentBraces += "[i" + strconv.Itoa(i) + "]"
		}

	}

	totalParentBraces = "obj" + totalParentBraces + "."

	newStringEncode := ""

	GenerateRustEncodeCode(&currentIterate, &newStringEncode, treeNode, totalParentBraces)

	*stringDataEncoder = strings.ReplaceAll(*stringDataEncoder, "{Packages}", rootClassName.Name)

	*stringDataEncoder += newStringEncode

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			*stringDataEncoder += `
			}
				`
		}
	}

	*stringDataEncoder += `
			return bb.to_array();
		}
	}
	`
}

func RustDecoderCodeGeneration(rootClassName RootClass, stringDataDecoder *string, packageName *string, fileName *string, treeNode *TreeNode) {
	squareBrackets := ""

	for i := 0; i < rootClassName.ArrayCount; i++ {
		squareBrackets += "Vec<"
	}

	if squareBrackets != "" {
		squareBrackets += rootClassName.Name
	}

	for i := 0; i < rootClassName.ArrayCount; i++ {
		squareBrackets += ">"
	}

	totalParentBraces := ""

	currentIterate := 0
	rootArrayClass := ""

	if squareBrackets != "" {
		rootArrayClass = `
			let arrLen = bb.get_short() as usize;
		`
		rootArrayClass += "let mut obj:" + squareBrackets + " = vec![vec![]; arrLen];"
	} else {
		rootArrayClass = "let mut obj = " + rootClassName.Name + "{..Default::default()};"
	}

	if squareBrackets == "" {
		squareBrackets = rootClassName.Name
	}

	*stringDataDecoder = `

	pub mod ` + *packageName + `{

	use crate::` + *packageName + `::RustBuffer::ByteBuff;
	use crate::` + *packageName + `::`+*fileName+`::` + *packageName + `::{{Packages}};

	pub fn ` + *fileName + `_Decoder(byte_arr: Vec<u8>) -> ` + squareBrackets + `{

		let mut bb = ByteBuff{
			multiplier: 10000.0,
			endian: "big".to_string(),
			..Default::default()
		};

		bb.init("big".to_string());

		bb.wrap(byte_arr);

		` + rootArrayClass + `
`
	if squareBrackets != "" {
		innerbracesCount := rootClassName.ArrayCount - 1
		for i := 0; i < rootClassName.ArrayCount; i++ {
			if i == 0 {
				*stringDataDecoder += `
				for i` + strconv.Itoa(i) + ` in 0..arrLen{
`
			} else {

				nestedSquareBrackets := ""

				for j := 0; j < innerbracesCount; j++ {
					nestedSquareBrackets += "Vec<"
				}

				nestedSquareBrackets += rootClassName.Name

				for j := 0; j < innerbracesCount; j++ {
					nestedSquareBrackets += ">"
				}

				innerbracesCount -= 1

				*stringDataDecoder += `

				let arrLen` + strconv.Itoa(i) + ` = bb.get_short() as usize;
`
				if i == rootClassName.ArrayCount-1 {
					*stringDataDecoder += `
					obj` + totalParentBraces + ` = vec![` + rootClassName.Name + `{..Default::default()}; arrLen` + strconv.Itoa(i) + `];
					`
				} else {
					*stringDataDecoder += `
					obj` + totalParentBraces + ` = vec![vec![]; arrLen` + strconv.Itoa(i) + `];
					`
				}

				*stringDataDecoder += `
				for i` + strconv.Itoa(i) + ` in 0..arrLen` + strconv.Itoa(i) + `{
				`
			}

			totalParentBraces += "[i" + strconv.Itoa(i) + "]"
		}

	}

	if rootClassName.ArrayCount > 0 {
		*stringDataDecoder += `
	obj` + totalParentBraces + ` = ` + treeNode.Value + `{..Default::default()};
	`
	}

	totalParentBraces = "obj" + totalParentBraces + "."

	newStringDecode := ""
	importPackage := rootClassName.Name + ","

	GenerateRustDecoderCode(&importPackage, &currentIterate, &newStringDecode, treeNode, totalParentBraces)

	importPackage = strings.TrimSuffix(importPackage, ",")
	*stringDataDecoder = strings.ReplaceAll(*stringDataDecoder, "{Packages}", importPackage)

	*stringDataDecoder += newStringDecode

	if squareBrackets != "" {
		for i := 0; i < rootClassName.ArrayCount; i++ {
			*stringDataDecoder += `
		}
			`
		}
	}

	*stringDataDecoder += `
		return obj;
	}
}
`
}
