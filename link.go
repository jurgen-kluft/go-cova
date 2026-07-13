package cova

import "fmt"

type Linker struct {
	VariableCapacity int
	FunctionCapacity int
}

func NewLinker(variableCapacity, functionCapacity int) *Linker {
	return &Linker{VariableCapacity: variableCapacity, FunctionCapacity: functionCapacity}
}

func (linker *Linker) Link(program *ProgramNode, compiled *RelocatableProgram) (*LinkedProgram, error) {
	if program == nil {
		return nil, fmt.Errorf("link error: program is nil")
	}
	if compiled == nil {
		return nil, fmt.Errorf("link error: compiled program is nil")
	}
	if linker == nil {
		return nil, fmt.Errorf("link error: linker is nil")
	}

	for _, binding := range compiled.ProgramSymbols.ExternSymbols {
		if binding.ByteOffset < 0 || binding.ByteOffset+binding.ByteSize > linker.VariableCapacity {
			return nil, fmt.Errorf("link error: extern variable %q requests byte range [%d,%d), but extern memory capacity is %d", binding.Name, binding.ByteOffset, binding.ByteOffset+binding.ByteSize, linker.VariableCapacity)
		}
		if binding.ByteAlignment > 1 && binding.ByteOffset%binding.ByteAlignment != 0 {
			return nil, fmt.Errorf("link error: extern variable %q byte offset %d is not aligned to %d", binding.Name, binding.ByteOffset, binding.ByteAlignment)
		}
	}

	tempToFunction := make(map[int]int, len(compiled.Functions))
	functions := make([]ScriptFunctionDescriptor, 0, len(compiled.Functions))
	paramKinds := make([]ValueKind, 0)
	paramOffsets := make([]int, 0)
	for _, binding := range compiled.Functions {
		switch binding.Scope {
		case ScopeExtern:
			if binding.SlotIndex < 0 || binding.SlotIndex >= linker.FunctionCapacity {
				return nil, fmt.Errorf("link error: host-linked function %q requests slot %d, but function capacity is %d", binding.Name, binding.SlotIndex, linker.FunctionCapacity)
			}
		case ScopeBSS:
			if binding.ParamCount != len(binding.ParamTypes) || binding.ParamCount != len(binding.ParamOffsets) {
				return nil, fmt.Errorf("link error: function %q has inconsistent parameter metadata", binding.Name)
			}
			paramStart := len(paramKinds)
			for index, typ := range binding.ParamTypes {
				kind := valueKindFromType(typ)
				if kind == KindNone || kind == KindVoid {
					return nil, fmt.Errorf("link error: function %q parameter %d has unsupported kind %d", binding.Name, index, kind)
				}
				paramKinds = append(paramKinds, kind)
				paramOffsets = append(paramOffsets, binding.ParamOffsets[index])
			}
			tempToFunction[binding.TempFuncID] = len(functions)
			functions = append(functions, ScriptFunctionDescriptor{
				BodyAddress:   binding.ScriptAddress,
				ParamStart:    paramStart,
				ParamCount:    binding.ParamCount,
				FrameByteSize: binding.FrameByteSize,
				ReturnKind:    valueKindFromType(binding.Type),
			})
		default:
			return nil, fmt.Errorf("link error: function %q has invalid scope %d", binding.Name, binding.Scope)
		}
	}

	linkedText := compiled.Text.Clone()
	for _, patch := range compiled.CallPatches {
		functionIndex, ok := tempToFunction[patch.TempFuncID]
		if !ok {
			return nil, fmt.Errorf("link error on line %d: unresolved function id %d", patch.Line, patch.TempFuncID)
		}
		linkedText.PatchInt(patch.OperandPos, functionIndex)
	}

	entryPoint, ok := tempToFunction[compiled.EntryFunction]
	if !ok {
		return nil, fmt.Errorf("link error: entry function id %d was not finalized", compiled.EntryFunction)
	}

	linked := &LinkedProgram{
		Text:          linkedText,
		EntryPoint:    entryPoint,
		Functions:     functions,
		ParamKinds:    paramKinds,
		ParamOffsets:  paramOffsets,
		FrameSize:     compiled.FrameSize,
		FrameByteSize: compiled.FrameByteSize,
		ConstByteSize: compiled.ConstByteSize,
		ConstData:     append([]byte(nil), compiled.ConstData...),
		DataByteSize:  compiled.DataByteSize,
		DataData:      append([]byte(nil), compiled.DataData...),
		BSSSize:       compiled.BSSSize,
		BSSByteSize:   compiled.BSSByteSize,
		DebugSymbols:  CopyProgramSymbols(compiled.ProgramSymbols),
	}
	return linked, nil
}
