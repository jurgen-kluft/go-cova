package cova

import "fmt"

func (core *parserCore) parseProgram() (*AstProgramNode, error) {
	program := &AstProgramNode{}
	for !core.isEOF() {
		if core.peek().Kind == TokExtern {
			decl, err := core.parseExternDecl()
			if err != nil {
				return nil, err
			}
			program.Decls = append(program.Decls, decl)
			continue
		}

		decl, function, err := core.parseTopLevelDeclOrFunction()
		if err != nil {
			return nil, err
		}
		if decl != nil {
			program.Decls = append(program.Decls, decl)
			continue
		}
		program.Functions = append(program.Functions, function)
	}

	return program, nil
}

func (core *parserCore) parseExternDecl() (*AstTopLevelDeclNode, error) {
	line := core.peek().Line
	if _, err := core.expect(TokExtern); err != nil {
		return nil, err
	}
	if _, err := core.expect(TokLParen); err != nil {
		return nil, err
	}
	indexToken, err := core.expect(TokInteger)
	if err != nil {
		return nil, err
	}
	if _, err := core.expect(TokRParen); err != nil {
		return nil, err
	}

	index := int(indexToken.IntValue)
	if core.peek().Kind == TokConst {
		return nil, fmt.Errorf("syntax error on line %d: extern declarations cannot be const", core.peek().Line)
	}

	typ, err := core.parseType()
	if err != nil {
		return nil, err
	}
	nameToken, err := core.expect(TokIdent)
	if err != nil {
		return nil, err
	}

	decl := &AstTopLevelDeclNode{Index: index, Name: nameToken.Text, Type: typ, Scope: ScopeExtern, Line: line}
	if core.match(TokLParen) {
		params, err := core.parseParameters()
		if err != nil {
			return nil, err
		}
		decl.Params = params
		decl.Kind = DeclFunction
		if _, err := core.expect(TokRParen); err != nil {
			return nil, err
		}
	} else {
		decl.Kind = DeclVariable
	}

	if _, err := core.expect(TokSemicolon); err != nil {
		return nil, err
	}
	return decl, nil
}

func (core *parserCore) parseTopLevelDeclOrFunction() (*AstTopLevelDeclNode, *AstFunctionNode, error) {
	line := core.peek().Line
	returnType, err := core.parseType()
	if err != nil {
		return nil, nil, err
	}
	nameToken, err := core.expect(TokIdent)
	if err != nil {
		return nil, nil, err
	}
	if core.match(TokLParen) {
		params, err := core.parseParameters()
		if err != nil {
			return nil, nil, err
		}
		if _, err := core.expect(TokRParen); err != nil {
			return nil, nil, err
		}
		body, err := core.parseBlock()
		if err != nil {
			return nil, nil, err
		}
		return nil, &AstFunctionNode{ReturnType: returnType, Name: nameToken.Text, Params: params, Body: body, Line: line}, nil
	}
	if returnType.Kind == TypeVoid {
		return nil, nil, fmt.Errorf("syntax error on line %d: internal variable %q cannot have type void", line, nameToken.Text)
	}
	var initializer AstExprNode
	scope := ScopeBSS
	if core.match(TokAssign) {
		initializer, err = core.parseExpression()
		if err != nil {
			return nil, nil, err
		}
		scope = ScopeData
	}
	if IsTopLevelConst(returnType) {
		scope = ScopeConst
	}
	if _, err := core.expect(TokSemicolon); err != nil {
		return nil, nil, err
	}
	decl := &AstTopLevelDeclNode{
		Index:       -1,
		Name:        nameToken.Text,
		Type:        returnType,
		Kind:        DeclVariable,
		Scope:       scope,
		Initializer: initializer,
		Line:        line,
	}
	return decl, nil, nil
}

func (core *parserCore) parseParameters() ([]AstParameter, error) {
	if core.peek().Kind == TokRParen {
		return nil, nil
	}

	params := make([]AstParameter, 0, 4)
	for {
		typ, err := core.parseType()
		if err != nil {
			return nil, err
		}
		nameToken, err := core.expect(TokIdent)
		if err != nil {
			return nil, err
		}
		params = append(params, AstParameter{Type: typ, Name: nameToken.Text, Line: nameToken.Line})

		if !core.match(TokComma) {
			break
		}
	}
	return params, nil
}