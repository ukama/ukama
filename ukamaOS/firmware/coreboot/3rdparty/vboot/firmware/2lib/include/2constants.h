/* Copyright 2019 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * General vboot-related constants.
 *
 * Constants that need to be exposed to assembly files or linker scripts
 * may be placed here and imported via vb2_constants.h.
 */

#ifndef VBOOT_REFERENCE_2CONSTANTS_H_
#define VBOOT_REFERENCE_2CONSTANTS_H_

/*
 * Size of non-volatile data used by vboot.
 *
 * If you only support non-volatile data format V1, then use VB2_NVDATA_SIZE.
 * If you support V2, use VB2_NVDATA_SIZE_V2 and set context flag
 * VB2_CONTEXT_NVDATA_V2.
 */
#define VB2_NVDATA_SIZE 16
#define VB2_NVDATA_SIZE_V2 64

/* Size of secure data spaces used by vboot */
#define VB2_SECDATA_FIRMWARE_SIZE 10
#define VB2_SECDATA_KERNEL_SIZE 13
#define VB2_SECDATA_FWMP_MIN_SIZE 40
#define VB2_SECDATA_FWMP_MAX_SIZE 64

/* TODO(chromium:972956): Remove once coreboot is using updated names */
#define VB2_SECDATA_SIZE 10
#define VB2_SECDATAK_SIZE 13

/*
 * Recommended size of work buffer for firmware verification stage.
 *
 * TODO: The recommended size really depends on which key algorithms are
 * used.  Should have a better / more accurate recommendation than this.
 */
#define VB2_FIRMWARE_WORKBUF_RECOMMENDED_SIZE (12 * 1024)

/*
 * Recommended size of work buffer for kernel verification stage.
 *
 * This is bigger because vboot 2.0 kernel preambles are usually padded to
 * 64 KB.
 *
 * TODO: The recommended size really depends on which key algorithms are
 * used.  Should have a better / more accurate recommendation than this.
 */
#define VB2_KERNEL_WORKBUF_RECOMMENDED_SIZE (80 * 1024)

/* Recommended buffer size for vb2api_get_pcr_digest. */
#define VB2_PCR_DIGEST_RECOMMENDED_SIZE 32

/*
 * Alignment for work buffer pointers/allocations should be useful for any
 * data type. When declaring workbuf buffers on the stack, the caller should
 * use explicit alignment to avoid run-time errors. For example:
 *
 *    int foo(void)
 *    {
 *        struct vb2_workbuf wb;
 *        uint8_t buf[NUM] __attribute__ ((aligned (VB2_WORKBUF_ALIGN)));
 *        wb.buf = buf;
 *        wb.size = sizeof(buf);
 */

/* We might get away with using __alignof__(void *), but since GCC defines a
 * macro for us we'll be safe and use that. */
#define VB2_WORKBUF_ALIGN __BIGGEST_ALIGNMENT__

/* Maximum length of a HWID in bytes, counting terminating null. */
#define VB2_GBB_HWID_MAX_SIZE 256

/* Type and offset of flags member in vb2_gbb_header struct. */
#define VB2_GBB_FLAGS_OFFSET 12
#ifndef __ASSEMBLER__
#include <stdint.h>
typedef uint32_t vb2_gbb_flags_t;
#endif

#endif  /* VBOOT_REFERENCE_2CONSTANTS_H_ */
