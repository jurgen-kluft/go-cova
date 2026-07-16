#include "ccova/linked_program.h"

namespace ncore
{
    template <typename T> static void validate_view(const array_view_t<T>& view)
    {
        ASSERT(view.m_data != nullptr || view.m_size == 0);
    }

    void validate_linked_program(const linked_program_t* program)
    {
        ASSERT(program != nullptr);
        ASSERT(program->m_text.m_code != nullptr || program->m_text.m_size == 0);
        validate_view(program->m_functions);
        validate_view(program->m_param_kinds);
        validate_view(program->m_param_offsets);
        validate_view(program->m_const_data);
        validate_view(program->m_data_data);
        ASSERT(program->m_param_kinds.m_size == program->m_param_offsets.m_size);
        ASSERT(program->m_functions.m_size != 0);
        ASSERT(program->m_entry_point < program->m_functions.m_size);

        for (u32 index = 0; index < program->m_param_kinds.m_size; ++index)
        {
            const evaluekind_t kind = program->m_param_kinds.m_data[index];
            ASSERT(kind > KindVoid && kind < KindCount);
        }

        for (u32 index = 0; index < program->m_functions.m_size; ++index)
        {
            const script_function_t& function = program->m_functions.m_data[index];
            ASSERT(function.m_body_address < program->m_text.m_size);
            ASSERT(function.m_return_kind > KindNone && function.m_return_kind < KindCount);
            ASSERT(function.m_param_start <= program->m_param_kinds.m_size);
            ASSERT(function.m_param_count <= program->m_param_kinds.m_size - function.m_param_start);
            ASSERT(function.m_param_start + function.m_param_count <= program->m_frame_size);
            ASSERT(function.m_frame_byte_size <= program->m_frame_byte_size);
        }
    }
}