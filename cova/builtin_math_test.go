package cova

import (
	"math"
	"strings"
	"testing"
)

func TestCompileAndRunMathBuiltInsFloat64(t *testing.T) {
	script := `
float64 script_main() {
	return abs(-2) + sin(0) + cos(0) + tan(0) + asin(0) + acos(1) + atan(0) + pow(2, 3) + sqrt(9);
}
`
	linked := mustLinkProgram(t, script, 0, 0)
	vm := NewVM(testFrameCapacityBytes)
	if status := vm.Run(linked); status != VMStatusOK {
		t.Fatalf("Run failed: %v", status)
	}
	if got, status := vm.PopFloat64(); status != VMStatusOK || got != 14 {
		t.Fatalf("math built-ins returned %v, want 14 (status %v)", got, status)
	}
}

func TestCompileAndRunMathBuiltInsFloat32(t *testing.T) {
	script := `
float32 script_main() {
	return sin(0.0f) + sqrt(9.0f) + pow(2.0f, 3.0f);
}
`
	linked := mustLinkProgram(t, script, 0, 0)
	vm := NewVM(testFrameCapacityBytes)
	if status := vm.Run(linked); status != VMStatusOK {
		t.Fatalf("Run failed: %v", status)
	}
	if got, status := vm.PopFloat32(); status != VMStatusOK || got != 11 {
		t.Fatalf("float32 math built-ins returned %v, want 11 (status %v)", got, status)
	}
}

func TestBuiltInAbsPreservesIntegerKindAndWrapsMinimum(t *testing.T) {
	var code CodeMemory
	appendOpcodeValue(&code, KindInt32, uint64(uint32(0x80000000)))
	code.AppendInstruction(makeBuiltInInstruction(makeBuiltInFunction(BuiltInAbs, KindInt32)))
	if got := runOpcodeResult(t, code, KindInt32); got != 0x80000000 {
		t.Fatalf("abs(int32 minimum) bits = %#x, want %#x", got, uint64(0x80000000))
	}

	var floatCode CodeMemory
	appendOpcodeValue(&floatCode, KindFloat64, math.Float64bits(-3.5))
	floatCode.AppendInstruction(makeBuiltInInstruction(makeBuiltInFunction(BuiltInAbs, KindFloat64)))
	if got := math.Float64frombits(runOpcodeResult(t, floatCode, KindFloat64)); got != 3.5 {
		t.Fatalf("abs(-3.5) = %v, want 3.5", got)
	}
}

func TestMathBuiltInDiagnostics(t *testing.T) {
	tests := []struct {
		name      string
		script    string
		wantError string
	}{
		{"arity", `float64 script_main() { return pow(2); }`, "expects 2 arguments"},
		{"bool", `float64 script_main() { return sin(true); }`, "requires numeric arguments"},
		{"reserved", `float64 sin(float64 value) { return value; } void script_main() { return; }`, "reserved built-in name"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := Tokenize(test.script)
			if err != nil {
				t.Fatalf("Tokenize failed: %v", err)
			}
			program, err := Parse(tokens)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			_, err = NewCompiler().Compile(program)
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("Compile error = %v, want substring %q", err, test.wantError)
			}
		})
	}
}
