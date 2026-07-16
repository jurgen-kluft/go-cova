#ifndef __CCOVA_PROGRAM_MEMORY_H__
#define __CCOVA_PROGRAM_MEMORY_H__

#include "ccova/segment_memory.h"

namespace ncore
{
    struct program_memory_t
    {
        segment_memory_t m_segments[SegmentCount];
    };

    void initialize_program_memory(program_memory_t* memory, const segment_memory_t& frame, const segment_memory_t& bss, const segment_memory_t& external, const segment_memory_t& constant, const segment_memory_t& data, const segment_memory_t& stack);

    bool      read_bool(const program_memory_t* memory, address_t address);
    s8        read_s8(const program_memory_t* memory, address_t address);
    s16       read_s16(const program_memory_t* memory, address_t address);
    s32       read_s32(const program_memory_t* memory, address_t address);
    s64       read_s64(const program_memory_t* memory, address_t address);
    u8        read_u8(const program_memory_t* memory, address_t address);
    u16       read_u16(const program_memory_t* memory, address_t address);
    u32       read_u32(const program_memory_t* memory, address_t address);
    u64       read_u64(const program_memory_t* memory, address_t address);
    f32       read_f32(const program_memory_t* memory, address_t address);
    f64       read_f64(const program_memory_t* memory, address_t address);
    address_t read_address(const program_memory_t* memory, address_t address);

    void write_bool(program_memory_t* memory, address_t address, bool value);
    void write_s8(program_memory_t* memory, address_t address, s8 value);
    void write_s16(program_memory_t* memory, address_t address, s16 value);
    void write_s32(program_memory_t* memory, address_t address, s32 value);
    void write_s64(program_memory_t* memory, address_t address, s64 value);
    void write_u8(program_memory_t* memory, address_t address, u8 value);
    void write_u16(program_memory_t* memory, address_t address, u16 value);
    void write_u32(program_memory_t* memory, address_t address, u32 value);
    void write_u64(program_memory_t* memory, address_t address, u64 value);
    void write_f32(program_memory_t* memory, address_t address, f32 value);
    void write_f64(program_memory_t* memory, address_t address, f64 value);
    void write_address(program_memory_t* memory, address_t destination, address_t value);
} // namespace ncore

#endif