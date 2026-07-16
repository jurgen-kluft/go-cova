#ifndef __CCOVA_FLOAT_BITS_H__
#define __CCOVA_FLOAT_BITS_H__

#include "ccova/types.h"

namespace ncore
{
    union f32_bits_t
    {
        f32 m_value;
        u32 m_bits;
    };

    union f64_bits_t
    {
        f64 m_value;
        u64 m_bits;
    };

    inline f32 bits_to_f32(u32 bits)
    {
        f32_bits_t value;
        value.m_bits = bits;
        return value.m_value;
    }

    inline f64 bits_to_f64(u64 bits)
    {
        f64_bits_t value;
        value.m_bits = bits;
        return value.m_value;
    }

    inline u32 f32_to_bits(f32 value)
    {
        f32_bits_t bits;
        bits.m_value = value;
        return bits.m_bits;
    }

    inline u64 f64_to_bits(f64 value)
    {
        f64_bits_t bits;
        bits.m_value = value;
        return bits.m_bits;
    }

    static_assert(sizeof(f32_bits_t) == sizeof(u32), "f32 bit conversion must be 32-bit");
    static_assert(sizeof(f64_bits_t) == sizeof(u64), "f64 bit conversion must be 64-bit");
} // namespace ncore

#endif