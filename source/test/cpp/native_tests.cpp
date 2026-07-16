#include "ccore/c_target.h"

#include "ccova/byte_order.h"
#include "ccova/code_memory.h"
#include "ccova/linked_program.h"
#include "ccova/segment_memory.h"

#include "cunittest/cunittest.h"

using namespace ncore;

UNITTEST_SUITE_BEGIN(cova_native)
{
    UNITTEST_FIXTURE(main)
    {
        UNITTEST_FIXTURE_SETUP() {}
        UNITTEST_FIXTURE_TEARDOWN() {}

        UNITTEST_TEST(byte_order)
        {
            byte bytes[8] = {};

            write_le_u16(bytes, 0xa5b6U);
            CHECK_EQUAL((u8)0xb6U, bytes[0]);
            CHECK_EQUAL((u8)0xa5U, bytes[1]);
            CHECK_EQUAL((u16)0xa5b6U, read_le_u16(bytes));

            write_le_u32(bytes, 0xa5b6c7d8U);
            CHECK_EQUAL((u8)0xd8U, bytes[0]);
            CHECK_EQUAL((u8)0xc7U, bytes[1]);
            CHECK_EQUAL((u8)0xb6U, bytes[2]);
            CHECK_EQUAL((u8)0xa5U, bytes[3]);
            CHECK_EQUAL((u32)0xa5b6c7d8U, read_le_u32(bytes));

            write_le_u64(bytes, 0x0123456789abcdefULL);
            CHECK_EQUAL((u64)0x0123456789abcdefULL, read_le_u64(bytes));
        }

        UNITTEST_TEST(instructions)
        {
            for (u32 opcode = OpPush; opcode < OpcodeCount; ++opcode)
            {
                const instruction_t instruction = make_instruction((eopcode_t)opcode, KindUint32);
                CHECK_EQUAL(opcode, (u32)instruction_opcode(instruction));
                CHECK_EQUAL((u32)KindUint32, (u32)instruction_kind(instruction));
            }

            for (u32 operation = ArithmeticAdd; operation <= ArithmeticDiv; ++operation)
            {
                const instruction_t instruction = make_arithmetic_instruction(KindInt16, (earithmeticop_t)operation);
                CHECK_EQUAL(operation, (u32)instruction_arithmetic_op(instruction));
            }

            for (u32 operation = CompareEqual; operation <= CompareGreaterEqual; ++operation)
            {
                const instruction_t instruction = make_compare_instruction(KindFloat32, (ecompareop_t)operation);
                CHECK_EQUAL(operation, (u32)instruction_compare_op(instruction));
            }

            const instruction_t conversion = make_convert_instruction(KindInt16, KindFloat64);
            CHECK_EQUAL((u32)KindInt16, (u32)instruction_convert_from_kind(conversion));
            CHECK_EQUAL((u32)KindFloat64, (u32)instruction_kind(conversion));

            const instruction_t address_instruction = make_address_instruction(SegmentExtern);
            CHECK_EQUAL((u32)SegmentExtern, (u32)instruction_address_segment(address_instruction));

            const address_t address = make_address(SegmentData, 0x00abcdefU);
            CHECK_EQUAL((u32)SegmentData, (u32)address_segment(address));
            CHECK_EQUAL((u32)0x00abcdefU, address_index(address));
        }

        UNITTEST_TEST(code_memory)
        {
            byte code[16] = {};
            write_le_u16(code, make_instruction(OpPush, KindUint32));
            write_le_u32(code + 2, 0xa5b6c7d8U);
            write_le_u64(code + 6, 0x0123456789abcdefULL);

            const code_memory_t memory = {code, 14};
            u32                 offset = 0;

            CHECK_EQUAL((u32)OpPush, (u32)instruction_opcode(read_instruction(&memory, &offset)));
            CHECK_EQUAL((u32)2, offset);
            CHECK_EQUAL((u64)0xa5b6c7d8U, read_immediate(&memory, &offset, KindUint32));
            CHECK_EQUAL((u32)6, offset);
            CHECK_EQUAL((u64)0x0123456789abcdefULL, read_immediate(&memory, &offset, KindUint64));
            CHECK_EQUAL((u32)14, offset);
        }

        UNITTEST_TEST(segment_memory_append_and_truncate)
        {
            byte             storage[32] = {};
            segment_memory_t segment     = {storage, 0, 32};

            append_u8(&segment, 0x12U);
            append_u16(&segment, 0x3456U);
            append_u32(&segment, 0x789abcdeU);
            append_u64(&segment, 0x0123456789abcdefULL);

            CHECK_EQUAL((u32)15, segment.m_size);
            CHECK_EQUAL((u64)0x0123456789abcdefULL, truncate_u64(&segment));
            CHECK_EQUAL((u32)0x789abcdeU, truncate_u32(&segment));
            CHECK_EQUAL((u16)0x3456U, truncate_u16(&segment));
            CHECK_EQUAL((u8)0x12U, truncate_u8(&segment));
            CHECK_EQUAL((u32)0, segment.m_size);
        }

        UNITTEST_TEST(segment_memory_read_write_and_copy)
        {
            byte             storage[32] = {};
            segment_memory_t segment     = {storage, 16, 32};

            write_u8(&segment, 0, 0xffU);
            write_u32(&segment, 4, 0x12345678U);
            write_u64(&segment, 8, 0x0102030405060708ULL);

            CHECK_EQUAL((u8)0xffU, read_u8(&segment, 0));
            CHECK_EQUAL((u32)0x12345678U, read_u32(&segment, 4));
            CHECK_EQUAL((u64)0x0102030405060708ULL, read_u64(&segment, 8));

            byte             destination_storage[8] = {};
            segment_memory_t destination            = {destination_storage, 8, 8};

            truncate_to(&segment, &destination, 0, 4);
            CHECK_EQUAL((u32)12, segment.m_size);
            CHECK_EQUAL((u32)0x01020304U, read_u32(&destination, 0));

            append_from(&segment, &destination, 0, 4);
            CHECK_EQUAL((u32)16, segment.m_size);
            CHECK_EQUAL((u32)0x01020304U, read_u32(&segment, 12));
        }

        UNITTEST_TEST(linked_program)
        {
            const byte              text[]       = {OpRet, 0};
            const script_function_t functions[] = {{0, 0, 0, 0, KindVoid}};
            const linked_program_t  program      = {
                {text, 2},
                0,
                {functions, 1},
                {nullptr, 0},
                {nullptr, 0},
                0,
                0,
                {nullptr, 0},
                {nullptr, 0},
                0,
            };

            validate_linked_program(&program);
            CHECK_EQUAL((u32)1, program.m_functions.m_size);
            CHECK_EQUAL((u32)0, program.m_entry_point);
        }
    }
}
UNITTEST_SUITE_END
