/* Copyright 2012 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Registers map and definitions for Cortex-MLM4x processor
 */

#ifndef __CROS_EC_CPU_H
#define __CROS_EC_CPU_H

#include <stdint.h>
#include "compile_time_macros.h"

/* Macro to access 32-bit registers */
#define CPUREG(addr) (*(volatile uint32_t*)(addr))

#define CPU_NVIC_ST_CTRL       CPUREG(0xE000E010)
#define ST_ENABLE              BIT(0)
#define ST_TICKINT             BIT(1)
#define ST_CLKSOURCE           BIT(2)
#define ST_COUNTFLAG           BIT(16)

/* Nested Vectored Interrupt Controller */
#define CPU_NVIC_EN(x)         CPUREG(0xe000e100 + 4 * (x))
#define CPU_NVIC_DIS(x)        CPUREG(0xe000e180 + 4 * (x))
#define CPU_NVIC_UNPEND(x)     CPUREG(0xe000e280 + 4 * (x))
#define CPU_NVIC_PRI(x)        CPUREG(0xe000e400 + 4 * (x))
#define CPU_NVIC_APINT         CPUREG(0xe000ed0c)
#define CPU_NVIC_SWTRIG        CPUREG(0xe000ef00)

#define CPU_SCB_SYSCTRL        CPUREG(0xe000ed10)

#define CPU_NVIC_CCR           CPUREG(0xe000ed14)
#define CPU_NVIC_SHCSR         CPUREG(0xe000ed24)
#define CPU_NVIC_MMFS          CPUREG(0xe000ed28)
#define CPU_NVIC_HFSR          CPUREG(0xe000ed2c)
#define CPU_NVIC_DFSR          CPUREG(0xe000ed30)
#define CPU_NVIC_MFAR          CPUREG(0xe000ed34)
#define CPU_NVIC_BFAR          CPUREG(0xe000ed38)

enum {
	CPU_NVIC_MMFS_BFARVALID		= BIT(15),
	CPU_NVIC_MMFS_MFARVALID		= BIT(7),

	CPU_NVIC_CCR_ICACHE		= BIT(17),
	CPU_NVIC_CCR_DCACHE		= BIT(16),
	CPU_NVIC_CCR_DIV_0_TRAP		= BIT(4),
	CPU_NVIC_CCR_UNALIGN_TRAP	= BIT(3),

	CPU_NVIC_HFSR_DEBUGEVT		= 1UL << 31,
	CPU_NVIC_HFSR_FORCED		= BIT(30),
	CPU_NVIC_HFSR_VECTTBL		= BIT(1),

	CPU_NVIC_SHCSR_MEMFAULTENA	= BIT(16),
	CPU_NVIC_SHCSR_BUSFAULTENA	= BIT(17),
	CPU_NVIC_SHCSR_USGFAULTENA	= BIT(18),
};

/* System Control Block: cache registers */
#define CPU_SCB_CCSIDR         CPUREG(0xe000ed80)
#define CPU_SCB_CCSELR         CPUREG(0xe000ed84)
#define CPU_SCB_ICIALLU        CPUREG(0xe000ef50)
#define CPU_SCB_DCISW          CPUREG(0xe000ef60)
#define CPU_SCB_DCCISW         CPUREG(0xe000ef74)

/* Set up the cpu to detect faults */
void cpu_init(void);
/* Enable the CPU I-cache and D-cache if they are not already enabled */
void cpu_enable_caches(void);
/* Invalidate the D-cache */
void cpu_invalidate_dcache(void);
/* Clean and Invalidate the D-cache to the Point of Coherency */
void cpu_clean_invalidate_dcache(void);

/* Invalidate a single range of the D-cache */
void cpu_invalidate_dcache_range(uintptr_t base, unsigned int length);
/* Clean and Invalidate a single range of the D-cache */
void cpu_clean_invalidate_dcache_range(uintptr_t base, unsigned int length);

#endif /* __CROS_EC_CPU_H */
