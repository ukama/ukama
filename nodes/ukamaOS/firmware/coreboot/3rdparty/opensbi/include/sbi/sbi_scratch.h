/*
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2019 Western Digital Corporation or its affiliates.
 *
 * Authors:
 *   Anup Patel <anup.patel@wdc.com>
 */

#ifndef __SBI_SCRATCH_H__
#define __SBI_SCRATCH_H__

#include <sbi/riscv_asm.h>

/* clang-format off */

/** Offset of fw_start member in sbi_scratch */
#define SBI_SCRATCH_FW_START_OFFSET		(0 * __SIZEOF_POINTER__)
/** Offset of fw_size member in sbi_scratch */
#define SBI_SCRATCH_FW_SIZE_OFFSET		(1 * __SIZEOF_POINTER__)
/** Offset of next_arg1 member in sbi_scratch */
#define SBI_SCRATCH_NEXT_ARG1_OFFSET		(2 * __SIZEOF_POINTER__)
/** Offset of next_addr member in sbi_scratch */
#define SBI_SCRATCH_NEXT_ADDR_OFFSET		(3 * __SIZEOF_POINTER__)
/** Offset of next_mode member in sbi_scratch */
#define SBI_SCRATCH_NEXT_MODE_OFFSET		(4 * __SIZEOF_POINTER__)
/** Offset of warmboot_addr member in sbi_scratch */
#define SBI_SCRATCH_WARMBOOT_ADDR_OFFSET	(5 * __SIZEOF_POINTER__)
/** Offset of platform_addr member in sbi_scratch */
#define SBI_SCRATCH_PLATFORM_ADDR_OFFSET	(6 * __SIZEOF_POINTER__)
/** Offset of hartid_to_scratch member in sbi_scratch */
#define SBI_SCRATCH_HARTID_TO_SCRATCH_OFFSET	(7 * __SIZEOF_POINTER__)
/** Offset of tmp0 member in sbi_scratch */
#define SBI_SCRATCH_TMP0_OFFSET			(8 * __SIZEOF_POINTER__)
/** Offset of options member in sbi_scratch */
#define SBI_SCRATCH_OPTIONS_OFFSET		(9 * __SIZEOF_POINTER__)
/** Offset of extra space in sbi_scratch */
#define SBI_SCRATCH_EXTRA_SPACE_OFFSET		(10 * __SIZEOF_POINTER__)
/** Maximum size of sbi_scratch and sbi_ipi_data */
#define SBI_SCRATCH_SIZE			(64 * __SIZEOF_POINTER__)

/* clang-format on */

#ifndef __ASSEMBLY__

#include <sbi/sbi_types.h>
#include <sbi/sbi_ipi.h>

/** Representation of per-HART scratch space */
struct sbi_scratch {
	/** Start (or base) address of firmware linked to OpenSBI library */
	unsigned long fw_start;
	/** Size (in bytes) of firmware linked to OpenSBI library */
	unsigned long fw_size;
	/** Arg1 (or 'a1' register) of next booting stage for this HART */
	unsigned long next_arg1;
	/** Address of next booting stage for this HART */
	unsigned long next_addr;
	/** Priviledge mode of next booting stage for this HART */
	unsigned long next_mode;
	/** Warm boot entry point address for this HART */
	unsigned long warmboot_addr;
	/** Address of sbi_platform */
	unsigned long platform_addr;
	/** Address of HART ID to sbi_scratch conversion function */
	unsigned long hartid_to_scratch;
	/** Temporary storage */
	unsigned long tmp0;
	/** Options for OpenSBI library */
	unsigned long options;
} __packed;

/** Possible options for OpenSBI library */
enum sbi_scratch_options {
	/** Disable prints during boot */
	SBI_SCRATCH_NO_BOOT_PRINTS = (1 << 0),
	/** Enable runtime debug prints */
	SBI_SCRATCH_DEBUG_PRINTS = (1 << 1),
};

/** Get pointer to sbi_scratch for current HART */
#define sbi_scratch_thishart_ptr() \
	((struct sbi_scratch *)csr_read(CSR_MSCRATCH))

/** Get Arg1 of next booting stage for current HART */
#define sbi_scratch_thishart_arg1_ptr() \
	((void *)(sbi_scratch_thishart_ptr()->next_arg1))

/** Allocate from extra space in sbi_scratch
 *
 * @return zero on failure and non-zero (>= SBI_SCRATCH_EXTRA_SPACE_OFFSET)
 * on success
 */
unsigned long sbi_scratch_alloc_offset(unsigned long size, const char *owner);

/** Free-up extra space in sbi_scratch */
void sbi_scratch_free_offset(unsigned long offset);

/** Get pointer from offset in sbi_scratch */
#define sbi_scratch_offset_ptr(scratch, offset)	((void *)scratch + (offset))

/** Get pointer from offset in sbi_scratch for current HART */
#define sbi_scratch_thishart_offset_ptr(offset)	\
	((void *)sbi_scratch_thishart_ptr() + (offset))

#endif

#endif
