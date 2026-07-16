package cova

func Parse(tokens []Token) (*AstProgramNode, error) {
	core := newParserCore(tokens)
	core.expr = newExpressionParser(&core)
	return core.parseProgram()
}
