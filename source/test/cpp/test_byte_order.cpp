#include "ccore/c_target.h"

#include "ccova/byte_order.h"

#include "cunittest/cunittest.h"

using namespace ncore;

UNITTEST_SUITE_BEGIN(cova_byte_order)
{
    UNITTEST_FIXTURE(main)
    {
        UNITTEST_FIXTURE_SETUP() {}
        UNITTEST_FIXTURE_TEARDOWN() {}

        UNITTEST_TEST(read_write)
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
    }
}
UNITTEST_SUITE_END
