package cova

import (
	"bytes"
	"testing"
)

func TestLinkerPrintCompilationReport(t *testing.T) {
	compiled := &RelocatableProgram{
		Functions: []SymbolBinding{
			{Kind: DeclFunction, Scope: ScopeBSS},
			{Kind: DeclFunction, Scope: ScopeExtern},
			{Kind: DeclFunction, Scope: ScopeExtern},
		},
	}
	linked := &LinkedProgram{
		Text:          make(CodeMemory, 23),
		BSSByteSize:   11,
		DataByteSize:  7,
		ConstByteSize: 13,
		Functions: []ScriptFunctionDescriptor{
			{},
			{},
		},
		DebugSymbols: &ProgramSymbols{
			ExternSymbols: []SymbolBinding{
				{Kind: DeclVariable, ByteOffset: 4, ByteSize: 4},
				{Kind: DeclVariable, ByteOffset: 16, ByteSize: 8},
				{Kind: DeclFunction, ByteOffset: 100, ByteSize: 20},
			},
		},
	}

	var output bytes.Buffer
	if err := NewLinker(24, 2).Report(&output, compiled, linked); err != nil {
		t.Fatalf("PrintCompilationReport failed: %v", err)
	}

	want := "Text size: 23 bytes\n" +
		"BSS size: 11 bytes\n" +
		"Data size: 7 bytes\n" +
		"Const size: 13 bytes\n" +
		"Local Functions: 2 functions\n" +
		"External Functions: 2 functions\n" +
		"External Variables: 2 variables, 24 bytes\n"
	if output.String() != want {
		t.Fatalf("unexpected report:\n%s\nwant:\n%s", output.String(), want)
	}
}

func TestLinkerPrintCompilationReportWithoutExternalVariables(t *testing.T) {
	compiled := &RelocatableProgram{}
	linked := &LinkedProgram{DebugSymbols: NewProgramSymbols()}

	var output bytes.Buffer
	if err := NewLinker(0, 0).Report(&output, compiled, linked); err != nil {
		t.Fatalf("PrintCompilationReport failed: %v", err)
	}

	want := "Text size: 0 bytes\n" +
		"BSS size: 0 bytes\n" +
		"Data size: 0 bytes\n" +
		"Const size: 0 bytes\n" +
		"Local Functions: 0 functions\n" +
		"External Functions: 0 functions\n" +
		"External Variables: 0 variables, 0 bytes\n"
	if output.String() != want {
		t.Fatalf("unexpected report:\n%s\nwant:\n%s", output.String(), want)
	}
}
