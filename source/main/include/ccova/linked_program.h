#ifndef __COVA_LINKED_PROGRAM_H__
#define __COVA_LINKED_PROGRAM_H__

#include "ccova/code_memory.h"

namespace ncore
{
    template <typename T> struct array_view_t
    {
        const T* m_data;
        u32      m_size;
    };

    struct script_function_t
    {
        u32          m_body_address;
        u32          m_param_start;
        u32          m_param_count;
        u32          m_frame_byte_size;
        evaluekind_t m_return_kind;
    };

    struct linked_program_t
    {
        code_memory_t                 m_text;
        u32                           m_entry_point;
        array_view_t<script_function_t> m_functions;
        array_view_t<evaluekind_t>      m_param_kinds;
        array_view_t<u32>               m_param_offsets;
        u32                           m_frame_size;
        u32                           m_frame_byte_size;
        array_view_t<byte>              m_const_data;
        array_view_t<byte>              m_data_data;
        u32                           m_bss_byte_size;
    };

    void validate_linked_program(const linked_program_t* program);
}

#endif