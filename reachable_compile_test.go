package cova

import (
	"strings"
	"testing"
)

func compileBlockSource(t *testing.T, script string) *RelocatableProgram {
	t.Helper()
	tokens, err := Tokenize(script)
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}
	program, err := Parse(tokens)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	compiled, err := NewCompiler().Compile(program)
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}
	return compiled
}

func TestCompileRequiresScriptMain(t *testing.T) {
	tokens, err := Tokenize("int helper() { return 1; }")
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}
	program, err := Parse(tokens)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	_, err = NewCompiler().Compile(program)
	if err == nil || !strings.Contains(err.Error(), "required entry function \"script_main\" not found") {
		t.Fatalf("expected missing script_main error, got %v", err)
	}
}

func TestCompileAssemblesReachableBlocksInSourceOrder(t *testing.T) {
	compiled := compileBlockSource(t, `
int dead() { return 99; }
int leaf() { return 3; }
int helper() { return leaf(); }
int script_main() { return helper(); }
int also_dead() { return 100; }
`)

	var names []string
	for _, function := range compiled.Functions {
		if function.Scope == ScopeBSS {
			names = append(names, function.Name)
		}
	}
	want := []string{"leaf", "helper", "script_main"}
	if strings.Join(names, ",") != strings.Join(want, ",") {
		t.Fatalf("local function order = %v, want %v", names, want)
	}
}

func TestCompileDeadBlockDoesNotContributeTextOrFrameMaximum(t *testing.T) {
	withDead := compileBlockSource(t, `
int64 dead(int64 a, int64 b, int64 c) { return a + b + c; }
int script_main() { return 7; }
`)
	withoutDead := compileBlockSource(t, `
int script_main() { return 7; }
`)
	if len(withDead.Text) != len(withoutDead.Text) {
		t.Fatalf("text with dead block = %d bytes, want %d", len(withDead.Text), len(withoutDead.Text))
	}
	if withDead.FrameByteSize != withoutDead.FrameByteSize {
		t.Fatalf("frame size with dead block = %d, want %d", withDead.FrameByteSize, withoutDead.FrameByteSize)
	}
}

func TestCompileStillValidatesDeadFunctionBodies(t *testing.T) {
	tokens, err := Tokenize(`
int dead() { return missing(); }
int script_main() { return 1; }
`)
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}
	program, err := Parse(tokens)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if _, err := NewCompiler().Compile(program); err == nil || !strings.Contains(err.Error(), "unknown function \"missing\"") {
		t.Fatalf("expected dead body compile error, got %v", err)
	}
}

func TestCompileIgnoresUnreachableRecursion(t *testing.T) {
	compiled := compileBlockSource(t, `
int dead_a() { return dead_b(); }
int dead_b() { return dead_a(); }
int script_main() { return 1; }
`)
	localFunctions := 0
	for _, function := range compiled.Functions {
		if function.Scope == ScopeBSS {
			localFunctions++
		}
	}
	if localFunctions != 1 {
		t.Fatalf("local function count = %d, want 1", localFunctions)
	}
}

func TestRunRelocatesReachableFunctionBlocks(t *testing.T) {
	linked := mustLinkProgram(t, `
int dead() {
	int value = 0;
	while (value < 20) { value = value + 1; }
	return value;
}

int helper(int value) {
	if (value > 2) { return value + 4; }
	return 0;
}

int script_main() {
	int total = 0;
	int index = 0;
	for (index = 0; index < 5; index = index + 1) {
		if (index == 1) { continue; }
		if (index == 4) { break; }
		total = total + index;
	}
	switch (total) {
	case 5: return helper(total);
	default: return 0;
	}
}
`, 0, 0)
	vm := NewVM(testFrameCapacityBytes)
	if status := vm.Run(linked); status != VMStatusOK {
		t.Fatalf("Run failed: %s", status)
	}
	value, status := vm.PopInt32()
	if status != VMStatusOK || value != 9 {
		t.Fatalf("result = %d, status = %s, want 9", value, status)
	}
}

func TestReachableFunctionBlocksRoundTripThroughImage(t *testing.T) {
	linked := mustLinkProgram(t, `
int dead() { return 99; }
int helper() { return 12; }
int script_main() { return helper(); }
`, 0, 0)
	imageBytes, status := BuildProgramImage(linked)
	if status != VMStatusOK {
		t.Fatalf("BuildProgramImage failed: %s", status)
	}
	image, status := OpenProgramImage(imageBytes)
	if status != VMStatusOK {
		t.Fatalf("OpenProgramImage failed: %s", status)
	}
	program, status := linkedProgramFromImage(image)
	if status != VMStatusOK {
		t.Fatalf("linkedProgramFromImage failed: %s", status)
	}
	vm := NewVM(testFrameCapacityBytes)
	if status := vm.Run(program); status != VMStatusOK {
		t.Fatalf("Run failed: %s", status)
	}
	value, status := vm.PopInt32()
	if status != VMStatusOK || value != 12 {
		t.Fatalf("result = %d, status = %s, want 12", value, status)
	}
}
