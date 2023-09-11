package parsers

type Parser interface {
	GenerateStruct(classDefinitions string, rootClassName RootClass, globalMap *map[string][]string) (string, RootClass)
	GenerateEncodeCode(currentIterate *int, stringData *string, node *TreeNode, parentName string)
	GenerateDecoderCode(currentIterate *int, stringData *string, node *TreeNode, parentName string)
	WriteEncoderData(encoderFileName string, stringDataEncoder string)
	WriteDecoderData(decoderFileName string, stringDataDecoder string)
	WriteStructData(modelFileName string, finalStruct string)
	EncoderCodeGeneration(rootClassName RootClass, stringDataEncoder *string, packageName *string, fileName *string, treeNode *TreeNode)
	DecoderCodeGeneration(rootClassName RootClass, stringDataDecoder *string, packageName *string, fileName *string, treeNode *TreeNode)
}
