
/**
 * This file was generated by TOP HO FATTO IL .H, a noi due! version 0.1.
 */

#include <string.h>

#include "test.h"



static inline uint8_t pack_left_shift_u8(
    uint8_t value,
    uint8_t shift,
    uint8_t mask)
{
    return (uint8_t)((uint8_t)(value << shift) & mask);
}

static inline uint8_t pack_right_shift_u8(
    uint8_t value,
    uint8_t shift,
    uint8_t mask)
{
    return (uint8_t)((uint8_t)(value >> shift) & mask);
}



static inline uint8_t pack_left_shift_u16(
    uint16_t value,
    uint8_t shift,
    uint8_t mask)
{
    return (uint8_t)((uint8_t)(value << shift) & mask);
}

static inline uint8_t pack_right_shift_u16(
    uint16_t value,
    uint8_t shift,
    uint8_t mask)
{
    return (uint8_t)((uint8_t)(value >> shift) & mask);
}





static inline uint8_t unpack_left_shift_u8(
    uint8_t value,
    uint8_t shift,
    uint8_t mask)
{
    return (uint8_t)((uint8_t)(value & mask) << shift);
}

static inline uint8_t unpack_right_shift_u8(
    uint8_t value,
    uint8_t shift,
    uint8_t mask)
{
    return (uint8_t)((uint8_t)(value & mask) >> shift);
}



static inline uint16_t unpack_left_shift_u16(
    uint8_t value,
    uint8_t shift,
    uint8_t mask)
{
    return (uint16_t)((uint16_t)(value & mask) << shift);
}

static inline uint16_t unpack_right_shift_u16(
    uint8_t value,
    uint8_t shift,
    uint8_t mask)
{
    return (uint16_t)((uint16_t)(value & mask) >> shift);
}




int expected_msg_0_pack(
    uint8_t *dst_p,
    const struct expected_msg_0_t *src_p,
    size_t size)
{
    

    if (size < 8u) {
        return (-EINVAL);
    }

    memset(&dst_p[0], 0, 8);

    // startBit = 0
    // size = 4
    dst_p[0] |= pack_left_shift_u8(src_p->std_sig_0, 0u, 0x0fu);
    // startBit = 4
    // size = 18
    dst_p[0] |= pack_left_shift_u32(src_p->mux_sig_0, 0u, 0xf0u);
    dst_p[1] |= pack_left_shift_u32(src_p->mux_sig_0, 8u, 0xffu);
    dst_p[2] |= pack_left_shift_u32(src_p->mux_sig_0, 16u, 0x3fu);
    

    return (8);
}

int expected_msg_0_unpack(
    struct expected_msg_0_t *dst_p,
    const uint8_t *src_p,
    size_t size)
{
    
    
    if (size < 8u) {
        return (-EINVAL);
    }

    dst_p->std_sig_0 = unpack_left_shift_u8(src_p[0], 0u, 0x0fu);
    dst_p->mux_sig_0 = unpack_left_shift_u32(src_p[0], 0u, 0xf0u);
    dst_p->mux_sig_0 |= unpack_left_shift_u32(src_p[1], 8u, 0xffu);
    dst_p->mux_sig_0 |= unpack_left_shift_u32(src_p[2], 16u, 0x3fu);
    

    return (0);
}


uint8_t expected_msg_0_std_sig_0_encode(double value)
{
    return (uint8_t)(value);
}
double expected_msg_0_std_sig_0_decode(uint8_t value)
{
    return ((float)value);
}
bool expected_msg_0_std_sig_0_is_in_range(uint8_t value)
{
    // 0 <= value <= 15
    return (value <= 15u);
    
}

32_t expected_msg_0_mux_sig_0_encode(double value)
{
    return (32_t)();
}
double expected_msg_0_mux_sig_0_decode(32_t value)
{
    return ();
}
bool expected_msg_0_mux_sig_0_is_in_range(32_t value)
{
    return (true);
}

int expected_msg_1_pack(
    uint8_t *dst_p,
    const struct expected_msg_1_t *src_p,
    size_t size)
{
    

    if (size < 8u) {
        return (-EINVAL);
    }

    memset(&dst_p[0], 0, 8);

    // startBit = 0
    // size = 18
    dst_p[0] |= pack_left_shift_u32(src_p->mux_sig_1, 0u, 0xffu);
    dst_p[1] |= pack_left_shift_u32(src_p->mux_sig_1, 8u, 0xffu);
    dst_p[2] |= pack_left_shift_u32(src_p->mux_sig_1, 16u, 0x03u);
    

    return (8);
}

int expected_msg_1_unpack(
    struct expected_msg_1_t *dst_p,
    const uint8_t *src_p,
    size_t size)
{
    
    
    if (size < 8u) {
        return (-EINVAL);
    }

    dst_p->mux_sig_1 = unpack_left_shift_u32(src_p[0], 0u, 0xffu);
    dst_p->mux_sig_1 |= unpack_left_shift_u32(src_p[1], 8u, 0xffu);
    dst_p->mux_sig_1 |= unpack_left_shift_u32(src_p[2], 16u, 0x03u);
    

    return (0);
}


32_t expected_msg_1_mux_sig_1_encode(double value)
{
    return (32_t)();
}
double expected_msg_1_mux_sig_1_decode(32_t value)
{
    return ();
}
bool expected_msg_1_mux_sig_1_is_in_range(32_t value)
{
    return (true);
}

int expected_msg_2_pack(
    uint8_t *dst_p,
    const struct expected_msg_2_t *src_p,
    size_t size)
{
    

    if (size < 8u) {
        return (-EINVAL);
    }

    memset(&dst_p[0], 0, 8);

    // startBit = 0
    // size = 4
    dst_p[0] |= pack_left_shift_u8(src_p->enum_sig_0, 0u, 0x0fu);
    

    return (8);
}

int expected_msg_2_unpack(
    struct expected_msg_2_t *dst_p,
    const uint8_t *src_p,
    size_t size)
{
    
    
    if (size < 8u) {
        return (-EINVAL);
    }

    dst_p->enum_sig_0 = unpack_left_shift_u8(src_p[0], 0u, 0x0fu);
    

    return (0);
}


uint8_t expected_msg_2_enum_sig_0_encode(double value)
{
    return (uint8_t)(value);
}
double expected_msg_2_enum_sig_0_decode(uint8_t value)
{
    return ((double)value);
}
bool expected_msg_2_enum_sig_0_is_in_range(uint8_t value)
{
    return (true);
}

int expected_msg_3_pack(
    uint8_t *dst_p,
    const struct expected_msg_3_t *src_p,
    size_t size)
{
    

    if (size < 1u) {
        return (-EINVAL);
    }

    memset(&dst_p[0], 0, 1);

    // startBit = 0
    // size = 4
    dst_p[0] |= pack_left_shift_u8(src_p->std_sig_1, 0u, 0x0fu);
    // startBit = 4
    // size = 4
    dst_p[0] |= pack_left_shift_u8(src_p->std_sig_2, 0u, 0xf0u);
    

    return (1);
}

int expected_msg_3_unpack(
    struct expected_msg_3_t *dst_p,
    const uint8_t *src_p,
    size_t size)
{
    
    
    if (size < 1u) {
        return (-EINVAL);
    }

    dst_p->std_sig_1 = unpack_left_shift_u8(src_p[0], 0u, 0x0fu);
    dst_p->std_sig_2 = unpack_left_shift_u8(src_p[0], 0u, 0xf0u);
    

    return (0);
}


uint8_t expected_msg_3_std_sig_1_encode(double value)
{
    return (uint8_t)(value);
}
double expected_msg_3_std_sig_1_decode(uint8_t value)
{
    return ((float)value);
}
bool expected_msg_3_std_sig_1_is_in_range(uint8_t value)
{
    // 0 <= value <= 15
    return (value <= 15u);
    
}

uint8_t expected_msg_3_std_sig_2_encode(double value)
{
    return (uint8_t)(value);
}
double expected_msg_3_std_sig_2_decode(uint8_t value)
{
    return ((float)value);
}
bool expected_msg_3_std_sig_2_is_in_range(uint8_t value)
{
    // 0 <= value <= 15
    return (value <= 15u);
    
}

