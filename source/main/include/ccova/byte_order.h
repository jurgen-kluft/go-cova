#ifndef __CCOVA_BYTE_ORDER_H__
#define __CCOVA_BYTE_ORDER_H__

#include "ccore/c_debug.h"

namespace ncore
{
    inline u16 read_le_u16(const byte* data)
    {
        ASSERT(data != nullptr);
        return (u16)((u16)data[0] | ((u16)data[1] << 8));
    }

    inline u32 read_le_u32(const byte* data)
    {
        ASSERT(data != nullptr);
        return (u32)((u32)data[0] | ((u32)data[1] << 8) | ((u32)data[2] << 16) | ((u32)data[3] << 24));
    }

    inline u64 read_le_u64(const byte* data)
    {
        ASSERT(data != nullptr);
        return (u64)((u64)data[0] | ((u64)data[1] << 8) | ((u64)data[2] << 16) | ((u64)data[3] << 24) | ((u64)data[4] << 32) | ((u64)data[5] << 40) | ((u64)data[6] << 48) | ((u64)data[7] << 56));
    }

    inline void write_le_u16(byte* data, u16 value)
    {
        ASSERT(data != nullptr);
        data[0] = (byte)(value & 0xffU);
        data[1] = (byte)((value >> 8) & 0xffU);
    }

    inline void write_le_u32(byte* data, u32 value)
    {
        ASSERT(data != nullptr);
        data[0] = (byte)(value & 0xffU);
        data[1] = (byte)((value >> 8) & 0xffU);
        data[2] = (byte)((value >> 16) & 0xffU);
        data[3] = (byte)((value >> 24) & 0xffU);
    }

    inline void write_le_u64(byte* data, u64 value)
    {
        ASSERT(data != nullptr);
        data[0] = (byte)(value & 0xffULL);
        data[1] = (byte)((value >> 8) & 0xffULL);
        data[2] = (byte)((value >> 16) & 0xffULL);
        data[3] = (byte)((value >> 24) & 0xffULL);
        data[4] = (byte)((value >> 32) & 0xffULL);
        data[5] = (byte)((value >> 40) & 0xffULL);
        data[6] = (byte)((value >> 48) & 0xffULL);
        data[7] = (byte)((value >> 56) & 0xffULL);
    }
} // namespace ncore

#endif