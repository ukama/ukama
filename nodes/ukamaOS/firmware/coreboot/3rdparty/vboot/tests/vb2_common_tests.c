/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Tests for firmware 2common.c
 */

#include "2common.h"
#include "2sysincludes.h"
#include "test_common.h"
#include "vboot_struct.h"  /* For old struct sizes */

/* Mock data */
static int counter_calls_left = 0;

/* Mock functions */
static int counter(void)
{
	counter_calls_left--;
	return 0;
}

/*
 * Test arithmetic-related macros and operators.
 */
static void test_arithmetic(void)
{
	int64_t a = -10, b = -20;
	uint64_t u = (0xabcd00000000ULL);
	uint64_t v = (0xabcd000000ULL);

	TEST_EQ(VB2_MIN(1, 2), 1, "MIN 1");
	TEST_EQ(VB2_MIN(4, 3), 3, "MIN 3");
	TEST_EQ(VB2_MIN(5, 5), 5, "MIN 5");
	TEST_EQ(VB2_MIN(a, b), b, "MIN uint64 1");
	TEST_EQ(VB2_MIN(b, a), b, "MIN uint64 2");
	TEST_EQ(VB2_MIN(b, b), b, "MIN uint64 same");

	counter_calls_left = 2;
	VB2_MIN(counter(), counter());
	TEST_EQ(counter_calls_left, 0, "MIN double-evaluation");

	TEST_EQ(VB2_MAX(1, 2), 2, "MAX 2");
	TEST_EQ(VB2_MAX(4, 3), 4, "MAX 4");
	TEST_EQ(VB2_MAX(5, 5), 5, "MAX 5");
	TEST_EQ(VB2_MAX(a, b), a, "MAX uint64 1");
	TEST_EQ(VB2_MAX(b, a), a, "MAX uint64 2");
	TEST_EQ(VB2_MAX(b, b), b, "MAX uint64 same");

	counter_calls_left = 2;
	VB2_MAX(counter(), counter());
	TEST_EQ(counter_calls_left, 0, "MAX double-evaluation");

	TEST_EQ(u >> 8, v, "uint64_t >> 8");
	TEST_EQ(u >> 0, u, "uint64_t >> 0");
	TEST_EQ(u >> 36, (uint64_t)0xabc, "uint64_t >> 36");

	TEST_EQ(v * (uint32_t)0, 0, "uint64_t * uint32_t 0");
	TEST_EQ(v * (uint32_t)1, v, "uint64_t * uint32_t 1");
	TEST_EQ(v * (uint32_t)256, u, "uint64_t * uint32_t 256");
}

/*
 * Test array size macro.
 */
static void test_array_size(void)
{
	uint8_t arr1[12];
	uint32_t arr2[7];
	uint64_t arr3[9];

	TEST_EQ(ARRAY_SIZE(arr1), 12, "ARRAYSIZE(uint8_t)");
	TEST_EQ(ARRAY_SIZE(arr2), 7, "ARRAYSIZE(uint32_t)");
	TEST_EQ(ARRAY_SIZE(arr3), 9, "ARRAYSIZE(uint64_t)");
}

/*
 * Test struct packing for vboot_struct.h structs which are passed between
 * firmware and OS, or passed between different phases of firmware.
 */
static void test_struct_packing(void)
{
	TEST_EQ(EXPECTED_VB2_PACKED_KEY_SIZE,
		sizeof(struct vb2_packed_key),
		"sizeof(vb2_packed_key)");
	TEST_EQ(EXPECTED_VB2_GBB_HEADER_SIZE,
		sizeof(struct vb2_gbb_header),
		"sizeof(vb2_gbb_header)");
	TEST_EQ(EXPECTED_VB2_SIGNATURE_SIZE,
		sizeof(struct vb2_signature),
		"sizeof(vb2_signature)");
	TEST_EQ(EXPECTED_VB2_KEYBLOCK_SIZE,
		sizeof(struct vb2_keyblock),
		"sizeof(vb2_keyblock)");
}

/**
 * Test memory compare functions
 */
static void test_memcmp(void)
{
	TEST_EQ(vb2_safe_memcmp("foo", "foo", 3), 0, "memcmp equal");
	TEST_NEQ(vb2_safe_memcmp("foo1", "foo2", 4), 0, "memcmp different");
	TEST_EQ(vb2_safe_memcmp("foo1", "foo2", 0), 0, "memcmp 0-size");
}

/**
 * Test alignment functions
 */
static void test_align(void)
{
	uint64_t buf[4];
	uint8_t *p0, *ptr;
	uint32_t size;

	/* Already aligned */
	p0 = (uint8_t *)buf;
	ptr = p0;
	size = 16;
	TEST_SUCC(vb2_align(&ptr, &size, 4, 16), "vb2_align() aligned");
	TEST_EQ(vb2_offset_of(p0, ptr), 0, "ptr");
	TEST_EQ(size, 16, "  size");
	TEST_EQ(vb2_align(&ptr, &size, 4, 17),
		VB2_ERROR_ALIGN_SIZE, "vb2_align() small");

	/* Offset */
	ptr = p0 + 1;
	size = 15;
	TEST_SUCC(vb2_align(&ptr, &size, 4, 12), "vb2_align() offset");
	TEST_EQ(vb2_offset_of(p0, ptr), 4, "ptr");
	TEST_EQ(size, 12, "  size");

	/* Offset, now too small */
	ptr = p0 + 1;
	size = 15;
	TEST_EQ(vb2_align(&ptr, &size, 4, 15),
		VB2_ERROR_ALIGN_SIZE, "vb2_align() offset small");

	/* Offset, too small even to align */
	ptr = p0 + 1;
	size = 1;
	TEST_EQ(vb2_align(&ptr, &size, 4, 1),
		VB2_ERROR_ALIGN_BIGGER_THAN_SIZE, "vb2_align() offset tiny");
}

/**
 * Test work buffer functions
 */
static void test_workbuf(void)
{
	uint64_t buf[8] __attribute__ ((aligned (VB2_WORKBUF_ALIGN)));
	uint8_t *p0 = (uint8_t *)buf, *ptr;
	struct vb2_workbuf wb;

	/* NOTE: There are several magic numbers below which assume that
	 * VB2_WORKBUF_ALIGN == 16 */

	/* Init */
	vb2_workbuf_init(&wb, p0, 64);
	TEST_EQ(vb2_offset_of(p0, wb.buf), 0, "Workbuf init aligned");
	TEST_EQ(wb.size, 64, "  size");

	vb2_workbuf_init(&wb, p0 + 4, 64);
	TEST_EQ(vb2_offset_of(p0, wb.buf), VB2_WORKBUF_ALIGN,
		"Workbuf init unaligned");
	TEST_EQ(wb.size, 64 - VB2_WORKBUF_ALIGN + 4, "  size");

	vb2_workbuf_init(&wb, p0 + 2, 5);
	TEST_EQ(wb.size, 0, "Workbuf init tiny unaligned size");

	/* Alloc rounds up */
	vb2_workbuf_init(&wb, p0, 64);
	ptr = vb2_workbuf_alloc(&wb, 22);
	TEST_EQ(vb2_offset_of(p0, ptr), 0, "Workbuf alloc");
	TEST_EQ(vb2_offset_of(p0, wb.buf), 32, "  buf");
	TEST_EQ(wb.size, 32, "  size");

	vb2_workbuf_init(&wb, p0, 32);
	TEST_PTR_EQ(vb2_workbuf_alloc(&wb, 33), NULL, "Workbuf alloc too big");

	/* Free reverses alloc */
	vb2_workbuf_init(&wb, p0, 32);
	vb2_workbuf_alloc(&wb, 22);
	vb2_workbuf_free(&wb, 22);
	TEST_EQ(vb2_offset_of(p0, wb.buf), 0, "Workbuf free buf");
	TEST_EQ(wb.size, 32, "  size");

	/* Realloc keeps same pointer as alloc */
	vb2_workbuf_init(&wb, p0, 64);
	vb2_workbuf_alloc(&wb, 6);
	ptr = vb2_workbuf_realloc(&wb, 6, 21);
	TEST_EQ(vb2_offset_of(p0, ptr), 0, "Workbuf realloc");
	TEST_EQ(vb2_offset_of(p0, wb.buf), 32, "  buf");
	TEST_EQ(wb.size, 32, "  size");
}

/**
 * Helper functions not dependent on specific key sizes
 */
static void test_helper_functions(void)
{
	{
		struct vb2_packed_key k = {.key_offset = sizeof(k)};
		TEST_EQ((int)vb2_offset_of(&k, vb2_packed_key_data(&k)),
			sizeof(k), "vb2_packed_key_data() adjacent");
	}

	{
		struct vb2_packed_key k = {.key_offset = 123};
		TEST_EQ((int)vb2_offset_of(&k, vb2_packed_key_data(&k)), 123,
			"vb2_packed_key_data() spaced");
	}
	{
		struct vb2_signature s = {.sig_offset = sizeof(s)};
		TEST_EQ((int)vb2_offset_of(&s, vb2_signature_data(&s)),
			sizeof(s), "vb2_signature_data() adjacent");
	}

	{
		struct vb2_signature s = {.sig_offset = 123};
		TEST_EQ((int)vb2_offset_of(&s, vb2_signature_data(&s)), 123,
			"vb2_signature_data() spaced");
	}

	{
		uint8_t *p = (uint8_t *)test_helper_functions;
		TEST_EQ((int)vb2_offset_of(p, p), 0, "vb2_offset_of() equal");
		TEST_EQ((int)vb2_offset_of(p, p+10), 10,
			"vb2_offset_of() positive");
		TEST_EQ((int)vb2_offset_of(p, p+0x12345678), 0x12345678,
			"vb2_offset_of() large");
	}

	{
		uint8_t *p = (uint8_t *)test_helper_functions;
		TEST_SUCC(vb2_verify_member_inside(p, 20, p, 6, 11, 3),
			  "vb2_verify_member_inside() ok 1");
		TEST_SUCC(vb2_verify_member_inside(p, 20, p+4, 4, 8, 4),
			  "vb2_verify_member_inside() ok 2");
		TEST_EQ(vb2_verify_member_inside(p, 20, p-4, 4, 8, 4),
			VB2_ERROR_INSIDE_MEMBER_OUTSIDE,
			"vb2_verify_member_inside() member before parent");
		TEST_EQ(vb2_verify_member_inside(p, 20, p+20, 4, 8, 4),
			VB2_ERROR_INSIDE_MEMBER_OUTSIDE,
			"vb2_verify_member_inside() member after parent");
		TEST_EQ(vb2_verify_member_inside(p, 20, p, 21, 0, 0),
			VB2_ERROR_INSIDE_MEMBER_OUTSIDE,
			"vb2_verify_member_inside() member too big");
		TEST_EQ(vb2_verify_member_inside(p, 20, p, 4, 21, 0),
			VB2_ERROR_INSIDE_DATA_OUTSIDE,
			"vb2_verify_member_inside() data after parent");
		TEST_EQ(vb2_verify_member_inside(p, 20, p, 4, SIZE_MAX, 0),
			VB2_ERROR_INSIDE_DATA_OUTSIDE,
			"vb2_verify_member_inside() data before parent");
		TEST_EQ(vb2_verify_member_inside(p, 20, p, 4, 4, 17),
			VB2_ERROR_INSIDE_DATA_OUTSIDE,
			"vb2_verify_member_inside() data too big");
		TEST_EQ(vb2_verify_member_inside(p, 20, p, 8, 4, 8),
			VB2_ERROR_INSIDE_DATA_OVERLAP,
			"vb2_verify_member_inside() data overlaps member");
		TEST_EQ(vb2_verify_member_inside(p, -8, p, 12, 0, 0),
			VB2_ERROR_INSIDE_PARENT_WRAPS,
			"vb2_verify_member_inside() wraparound 1");
		TEST_EQ(vb2_verify_member_inside(p, 20, p, -8, 0, 0),
			VB2_ERROR_INSIDE_MEMBER_WRAPS,
			"vb2_verify_member_inside() wraparound 2");
		TEST_EQ(vb2_verify_member_inside(p, 20, p, 4, 4, -12),
			VB2_ERROR_INSIDE_DATA_WRAPS,
			"vb2_verify_member_inside() wraparound 3");
	}

	{
		struct vb2_packed_key k = {.key_offset = sizeof(k),
					   .key_size = 128};
		TEST_SUCC(vb2_verify_packed_key_inside(&k, sizeof(k)+128, &k),
			  "vb2_packed_key_inside() ok 1");
		TEST_SUCC(vb2_verify_packed_key_inside(&k - 1,
						       2*sizeof(k)+128, &k),
			  "vb2_packed_key_inside() ok 2");
		TEST_EQ(vb2_verify_packed_key_inside(&k, 128, &k),
			VB2_ERROR_INSIDE_DATA_OUTSIDE,
			"vb2_packed_key_inside() key too big");
	}

	{
		struct vb2_packed_key k = {.key_offset = 100,
					   .key_size = 4};
		TEST_EQ(vb2_verify_packed_key_inside(&k, 99, &k),
			VB2_ERROR_INSIDE_DATA_OUTSIDE,
			"vb2_packed_key_inside() offset too big");
	}

	{
		struct vb2_signature s = {.sig_offset = sizeof(s),
					  .sig_size = 128};
		TEST_SUCC(vb2_verify_signature_inside(&s, sizeof(s)+128, &s),
			"vb2_verify_signature_inside() ok 1");
		TEST_SUCC(vb2_verify_signature_inside(&s - 1,
						      2*sizeof(s)+128, &s),
			  "vb2_verify_signature_inside() ok 2");
		TEST_EQ(vb2_verify_signature_inside(&s, 128, &s),
			VB2_ERROR_INSIDE_DATA_OUTSIDE,
			"vb2_verify_signature_inside() sig too big");
	}

	{
		struct vb2_signature s = {.sig_offset = 100,
					  .sig_size = 4};
		TEST_EQ(vb2_verify_signature_inside(&s, 99, &s),
			VB2_ERROR_INSIDE_DATA_OUTSIDE,
			"vb2_verify_signature_inside() offset too big");
	}
}

/* Helper for test_assert_die() below */
static int _true_assertion_helper(void)
{
	VB2_ASSERT(2 + 2 == 4);
	return 1;
}

/**
 * Test VB2_ASSERT and VB2_DIE macros
 */
static void test_assert_die(void)
{
	TEST_ABORT(VB2_DIE("die"), "DIE should abort");
	TEST_ABORT(VB2_ASSERT(2 + 2 == 5), "ASSERT false should abort");
	TEST_TRUE(_true_assertion_helper(), "ASSERT true should continue");
}

int main(int argc, char* argv[])
{
	test_arithmetic();
	test_array_size();
	test_struct_packing();
	test_memcmp();
	test_align();
	test_workbuf();
	test_helper_functions();
	test_assert_die();

	return gTestSuccess ? 0 : 255;
}
