/*
 * This file is part of the coreboot project.
 *
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#ifndef _RULES_H
#define _RULES_H

/* Useful helpers to tell whether the code is executing in bootblock,
 * romstage, ramstage or SMM.
 */

#if defined(__DECOMPRESSOR__)
#define ENV_DECOMPRESSOR 1
#define ENV_BOOTBLOCK 0
#define ENV_ROMSTAGE 0
#define ENV_RAMSTAGE 0
#define ENV_SMM 0
#define ENV_VERSTAGE 0
#define ENV_RMODULE 0
#define ENV_POSTCAR 0
#define ENV_LIBAGESA 0
#define ENV_STRING "decompressor"

#elif defined(__BOOTBLOCK__)
#define ENV_DECOMPRESSOR 0
#define ENV_BOOTBLOCK 1
#define ENV_ROMSTAGE 0
#define ENV_RAMSTAGE 0
#define ENV_SMM 0
#define ENV_VERSTAGE 0
#define ENV_RMODULE 0
#define ENV_POSTCAR 0
#define ENV_LIBAGESA 0
#define ENV_STRING "bootblock"

#elif defined(__ROMSTAGE__)
#define ENV_DECOMPRESSOR 0
#define ENV_BOOTBLOCK 0
#define ENV_ROMSTAGE 1
#define ENV_RAMSTAGE 0
#define ENV_SMM 0
#define ENV_VERSTAGE 0
#define ENV_RMODULE 0
#define ENV_POSTCAR 0
#define ENV_LIBAGESA 0
#define ENV_STRING "romstage"

#elif defined(__SMM__)
#define ENV_DECOMPRESSOR 0
#define ENV_BOOTBLOCK 0
#define ENV_ROMSTAGE 0
#define ENV_RAMSTAGE 0
#define ENV_SMM 1
#define ENV_VERSTAGE 0
#define ENV_RMODULE 0
#define ENV_POSTCAR 0
#define ENV_LIBAGESA 0
#define ENV_STRING "smm"

#elif defined(__VERSTAGE__)
#define ENV_DECOMPRESSOR 0
#define ENV_BOOTBLOCK 0
#define ENV_ROMSTAGE 0
#define ENV_RAMSTAGE 0
#define ENV_SMM 0
#define ENV_VERSTAGE 1
#define ENV_RMODULE 0
#define ENV_POSTCAR 0
#define ENV_LIBAGESA 0
#define ENV_STRING "verstage"

#elif defined(__RAMSTAGE__)
#define ENV_DECOMPRESSOR 0
#define ENV_BOOTBLOCK 0
#define ENV_ROMSTAGE 0
#define ENV_RAMSTAGE 1
#define ENV_SMM 0
#define ENV_VERSTAGE 0
#define ENV_RMODULE 0
#define ENV_POSTCAR 0
#define ENV_LIBAGESA 0
#define ENV_STRING "ramstage"

#elif defined(__RMODULE__)
#define ENV_DECOMPRESSOR 0
#define ENV_BOOTBLOCK 0
#define ENV_ROMSTAGE 0
#define ENV_RAMSTAGE 0
#define ENV_SMM 0
#define ENV_VERSTAGE 0
#define ENV_RMODULE 1
#define ENV_POSTCAR 0
#define ENV_LIBAGESA 0
#define ENV_STRING "rmodule"

#elif defined(__POSTCAR__)
#define ENV_DECOMPRESSOR 0
#define ENV_BOOTBLOCK 0
#define ENV_ROMSTAGE 0
#define ENV_RAMSTAGE 0
#define ENV_SMM 0
#define ENV_VERSTAGE 0
#define ENV_RMODULE 0
#define ENV_POSTCAR 1
#define ENV_LIBAGESA 0
#define ENV_STRING "postcar"

#elif defined(__LIBAGESA__)
#define ENV_DECOMPRESSOR 0
#define ENV_BOOTBLOCK 0
#define ENV_ROMSTAGE 0
#define ENV_RAMSTAGE 0
#define ENV_SMM 0
#define ENV_VERSTAGE 0
#define ENV_RMODULE 0
#define ENV_POSTCAR 0
#define ENV_LIBAGESA 1
#define ENV_STRING "libagesa"

#else
/*
 * Default case of nothing set for random blob generation using
 * create_class_compiler that isn't bound to a stage.
 */
#define ENV_DECOMPRESSOR 0
#define ENV_BOOTBLOCK 0
#define ENV_ROMSTAGE 0
#define ENV_RAMSTAGE 0
#define ENV_SMM 0
#define ENV_VERSTAGE 0
#define ENV_RMODULE 0
#define ENV_POSTCAR 0
#define ENV_LIBAGESA 0
#define ENV_STRING "UNKNOWN"
#endif

/* Define helpers about the current architecture, based on toolchain.inc. */

#if defined(__ARCH_arm__)
#define ENV_ARM 1
#define ENV_ARM64 0
#if __COREBOOT_ARM_ARCH__ == 4
#define ENV_ARMV4 1
#define ENV_ARMV7 0
#elif __COREBOOT_ARM_ARCH__ == 7
#define ENV_ARMV4 0
#define ENV_ARMV7 1
#if defined(__COREBOOT_ARM_V7_A__)
#define ENV_ARMV7_A 1
#define ENV_ARMV7_M 0
#define ENV_ARMV7_R 0
#elif defined(__COREBOOT_ARM_V7_M__)
#define ENV_ARMV7_A 0
#define ENV_ARMV7_M 1
#define ENV_ARMV7_R 0
#elif defined(__COREBOOT_ARM_V7_R__)
#define ENV_ARMV7_A 0
#define ENV_ARMV7_M 0
#define ENV_ARMV7_R 1
#endif
#else
#define ENV_ARMV4 0
#define ENV_ARMV7 0
#endif
#define ENV_ARMV8 0
#define ENV_MIPS 0
#define ENV_RISCV 0
#define ENV_X86 0
#define ENV_X86_32 0
#define ENV_X86_64 0

#elif defined(__ARCH_arm64__)
#define ENV_ARM 0
#define ENV_ARM64 1
#define ENV_ARMV4 0
#define ENV_ARMV7 0
#if __COREBOOT_ARM_ARCH__ == 8
#define ENV_ARMV8 1
#else
#define ENV_ARMV8 0
#endif
#define ENV_MIPS 0
#define ENV_RISCV 0
#define ENV_X86 0
#define ENV_X86_32 0
#define ENV_X86_64 0

#elif defined(__ARCH_mips__)
#define ENV_ARM 0
#define ENV_ARM64 0
#define ENV_ARMV4 0
#define ENV_ARMV7 0
#define ENV_ARMV8 0
#define ENV_MIPS 1
#define ENV_RISCV 0
#define ENV_X86 0
#define ENV_X86_32 0
#define ENV_X86_64 0

#elif defined(__ARCH_riscv__)
#define ENV_ARM 0
#define ENV_ARM64 0
#define ENV_ARMV4 0
#define ENV_ARMV7 0
#define ENV_ARMV8 0
#define ENV_MIPS 0
#define ENV_RISCV 1
#define ENV_X86 0
#define ENV_X86_32 0
#define ENV_X86_64 0

#elif defined(__ARCH_x86_32__)
#define ENV_ARM 0
#define ENV_ARM64 0
#define ENV_ARMV4 0
#define ENV_ARMV7 0
#define ENV_ARMV8 0
#define ENV_MIPS 0
#define ENV_RISCV 0
#define ENV_X86 1
#define ENV_X86_32 1
#define ENV_X86_64 0

#elif defined(__ARCH_x86_64__)
#define ENV_ARM 0
#define ENV_ARM64 0
#define ENV_ARMV4 0
#define ENV_ARMV7 0
#define ENV_ARMV8 0
#define ENV_MIPS 0
#define ENV_RISCV 0
#define ENV_X86 1
#define ENV_X86_32 0
#define ENV_X86_64 1

#else
#define ENV_ARM 0
#define ENV_ARM64 0
#define ENV_ARMV4 0
#define ENV_ARMV7 0
#define ENV_ARMV8 0
#define ENV_MIPS 0
#define ENV_RISCV 0
#define ENV_X86 0
#define ENV_X86_32 0
#define ENV_X86_64 0

#endif

#if CONFIG(RAMPAYLOAD)
/* ENV_PAYLOAD_LOADER is set to ENV_POSTCAR when CONFIG_RAMPAYLOAD is enabled */
#define ENV_PAYLOAD_LOADER ENV_POSTCAR
#else
/* ENV_PAYLOAD_LOADER is set when you are in a stage that loads the payload.
 * For now, that is the ramstage. */
#define ENV_PAYLOAD_LOADER ENV_RAMSTAGE
#endif

#define ENV_ROMSTAGE_OR_BEFORE \
	(ENV_DECOMPRESSOR || ENV_BOOTBLOCK || ENV_ROMSTAGE || \
	(ENV_VERSTAGE && CONFIG(VBOOT_STARTS_IN_BOOTBLOCK)))

#if CONFIG(ARCH_X86)
/* Indicates memory layout is determined with arch/x86/car.ld. */
#define ENV_CACHE_AS_RAM		(ENV_ROMSTAGE_OR_BEFORE && !CONFIG(RESET_VECTOR_IN_RAM))
/* No .data sections with execute-in-place from ROM.  */
#define ENV_STAGE_HAS_DATA_SECTION	!ENV_CACHE_AS_RAM
/* No .bss sections for stage with CAR teardown. */
#define ENV_STAGE_HAS_BSS_SECTION	!(ENV_ROMSTAGE && CONFIG(CAR_GLOBAL_MIGRATION))
#else
/* Both .data and .bss, sometimes SRAM not DRAM. */
#define ENV_STAGE_HAS_DATA_SECTION	1
#define ENV_STAGE_HAS_BSS_SECTION	1
#define ENV_CACHE_AS_RAM		0
#endif

/* Currently rmodules, ramstage and smm have heap. */
#define ENV_STAGE_HAS_HEAP_SECTION	(ENV_RMODULE || ENV_RAMSTAGE || ENV_SMM)

/**
 * For pre-DRAM stages and post-CAR always build with simple device model, ie.
 * PCI, PNP and CPU functions operate without use of devicetree. The reason
 * post-CAR utilizes __SIMPLE_DEVICE__ is for simplicity. Currently there's
 * no known requirement that devicetree would be needed during that stage.
 *
 * For ramstage individual source file may define __SIMPLE_DEVICE__
 * before including any header files to force that particular source
 * be built with simple device model.
 */

#if !ENV_RAMSTAGE
#define __SIMPLE_DEVICE__
#endif

#endif /* _RULES_H */
