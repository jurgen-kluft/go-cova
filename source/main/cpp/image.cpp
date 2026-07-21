#include "ccore/c_debug.h"
#include "ccova/image.h"

namespace ncore
{
    const linked_program_t* open_program_image(const byte* block, u32 block_size)
    {
        ASSERT(block != nullptr);
        ASSERT(block_size >= ProgramImageHeaderSize);
        ASSERT(((uint_t)block & 3U) == 0);

        const linked_program_t* program = (const linked_program_t*)block;
        validate_linked_program(program, block_size);
        return program;
    }

    static void validate()
    {
        CC_STATIC_ASSERTS(sizeof(script_function_t) == ProgramImageFunctionSize, "script function must match the image ABI");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(script_function_t, m_body_address) == 0, "script function body address ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(script_function_t, m_param_start) == 4, "script function parameter start ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(script_function_t, m_param_count) == 8, "script function parameter count ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(script_function_t, m_frame_byte_size) == 12, "script function frame size ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(script_function_t, m_return_kind) == 16, "script function return kind ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_magic) == 0, "program magic ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_entry_point) == 8, "program entry point ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_bss_byte_size) == 12, "program BSS size ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_frame_size) == 16, "program frame size ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_frame_byte_size) == 20, "program frame byte size ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_functions) == 24, "program functions ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_param_kinds) == 32, "program parameter kinds ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_param_offsets) == 40, "program parameter offsets ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_text) == 48, "program text ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_const_data) == 56, "program const data ABI mismatch");
        // CC_STATIC_ASSERTS(CC_OFFSETOF(linked_program_t, m_data_data) == 64, "program data ABI mismatch");
    }
} // namespace ncore
