#ifndef __CCOVA_TYPES_H__
#define __CCOVA_TYPES_H__

#include "ccore/c_debug.h"

namespace ncore
{
    enum eopcode_t : u8
    {
        OpPush = 1,
        OpArithmetic,
        OpConvert,
        OpAddr,
        OpOffset,
        OpDereference,
        OpAssign,
        OpCompare,
        OpJumpIfFalse,
        OpJump,
        OpCall,
        OpCallExtern,
        OpRet,
        OpcodeCount,
    };

    enum earithmeticop_t : u8
    {
        ArithmeticInvalid = 0,
        ArithmeticAdd,
        ArithmeticSub,
        ArithmeticMul,
        ArithmeticDiv,
    };

    enum ecompareop_t : u8
    {
        CompareInvalid = 0,
        CompareEqual,
        CompareNotEqual,
        CompareLess,
        CompareLessEqual,
        CompareGreater,
        CompareGreaterEqual,
    };

    enum evaluekind_t : u8
    {
        KindNone = 0,
        KindVoid,
        KindBool,
        KindByte,
        KindInt8,
        KindInt16,
        KindInt32,
        KindInt64,
        KindUint8,
        KindUint16,
        KindUint32,
        KindUint64,
        KindFloat32,
        KindFloat64,
        KindAddress,
        KindCount,
    };

    enum ememorysegment_t : u8
    {
        SegmentInvalid = 0,
        SegmentFrame,
        SegmentBSS,
        SegmentExtern,
        SegmentConst,
        SegmentData,
        SegmentStack,
        SegmentReserved0,
        SegmentReserved1,
        SegmentCount,
    };

    typedef u16 instruction_t;
    typedef u32 address_t;

    static const u32 AddressIndexMask = 0x00ffffffU;

    u32 value_kind_size(evaluekind_t kind);

    instruction_t make_instruction(eopcode_t opcode, evaluekind_t kind);
    instruction_t make_arithmetic_instruction(evaluekind_t kind, earithmeticop_t operation);
    instruction_t make_address_instruction(ememorysegment_t segment);
    instruction_t make_compare_instruction(evaluekind_t kind, ecompareop_t operation);
    instruction_t make_convert_instruction(evaluekind_t from, evaluekind_t to);

    eopcode_t        instruction_opcode(instruction_t instruction);
    evaluekind_t     instruction_kind(instruction_t instruction);
    earithmeticop_t  instruction_arithmetic_op(instruction_t instruction);
    ememorysegment_t instruction_address_segment(instruction_t instruction);
    ecompareop_t     instruction_compare_op(instruction_t instruction);
    evaluekind_t     instruction_convert_from_kind(instruction_t instruction);

    address_t        make_address(ememorysegment_t segment, u32 index);
    ememorysegment_t address_segment(address_t address);
    u32              address_index(address_t address);

    ASSERTCTS(sizeof(instruction_t) == 2, "instruction ABI must be 16-bit");
    ASSERTCTS(sizeof(address_t) == 4, "address ABI must be 32-bit");
    ASSERTCTS(OpcodeCount < 32, "opcode count must fit the instruction encoding");
    ASSERTCTS(KindCount <= 16, "value kind must fit four bits");
} // namespace ncore

#endif