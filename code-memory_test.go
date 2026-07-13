package cova

import (
	"strings"
	"testing"
)

func TestValueKindSizeRejectsInvalidKind(t *testing.T) {
	if got := ValueKind(255).Size(); got != 0 {
		t.Fatalf("expected invalid kind size 0, got %d", got)
	}
}

func TestCheckedCodeReadersRejectTruncatedInput(t *testing.T) {
	tests := []struct {
		name string
		read func(CodeMemory) error
	}{
		{
			name: "instruction",
			read: func(code CodeMemory) error {
				ip := 0
				_, err := code.ReadInstructionChecked(&ip)
				return err
			},
		},
		{
			name: "immediate",
			read: func(code CodeMemory) error {
				ip := 0
				_, err := code.ReadImmediateChecked(&ip, KindUint32)
				return err
			},
		},
		{
			name: "int operand",
			read: func(code CodeMemory) error {
				ip := 0
				_, err := code.ReadIntChecked(&ip)
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.read(CodeMemory{0}); err == nil {
				t.Fatal("expected truncated input error")
			}
		})
	}
}

func TestRunRejectsTruncatedInstructionWithoutPanic(t *testing.T) {
	code := CodeMemory{byte(OpRet)}
	program := &LinkedProgram{
		Text:       code,
		EntryPoint: 0,
		Functions:  []ScriptFunctionDescriptor{{BodyAddress: 0, ReturnKind: KindVoid}},
	}

	if err := NewVM(0).Run(program); err == nil || !strings.Contains(err.Error(), "truncated instruction") {
		t.Fatalf("expected truncated instruction error, got %v", err)
	}
}

func TestRunRejectsTruncatedOperandWithoutPanic(t *testing.T) {
	code := CodeMemory{}
	code.AppendInstruction(makeInstruction(OpJump, KindNone, ModeNone, FlagNone))
	code = append(code, 0)
	program := &LinkedProgram{
		Text:       code,
		EntryPoint: 0,
		Functions:  []ScriptFunctionDescriptor{{BodyAddress: 0, ReturnKind: KindVoid}},
	}

	if err := NewVM(0).Run(program); err == nil || !strings.Contains(err.Error(), "truncated 4-byte operand") {
		t.Fatalf("expected truncated operand error, got %v", err)
	}
}
