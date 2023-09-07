/*
 * This file is part of the coreboot project.
 *
 * Copyright 2014 Google Inc.
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

/* This file contains macro definitions for memlayout.ld linker scripts. */

#ifndef __MEMLAYOUT_H
#define __MEMLAYOUT_H

#include <arch/memlayout.h>
#include <vb2_constants.h>

#include "fmap_config.h"

/* Macros that the architecture can override. */
#ifndef ARCH_POINTER_ALIGN_SIZE
#define ARCH_POINTER_ALIGN_SIZE 8
#endif

#ifndef ARCH_CACHELINE_ALIGN_SIZE
#define ARCH_CACHELINE_ALIGN_SIZE 64
#endif

#define STR(x) XSTR(x)
#define XSTR(x) #x

#define ALIGN_COUNTER(align) \
	. = ALIGN(align);

#define SET_COUNTER(name, addr) \
	_ = ASSERT(. <= addr, STR(name overlaps the previous region!)); \
	. = addr;

#define SYMBOL(name, addr) \
	SET_COUNTER(name, addr) \
	_##name = .;

#define REGION(name, addr, size, expected_align) \
	SYMBOL(name, addr) \
	_ = ASSERT(. == ALIGN(expected_align), \
		STR(name must be aligned to expected_align!)); \
	SYMBOL(e##name, addr + size)

#define ALIAS_REGION(name, alias) \
	_##alias = _##name; \
	_e##alias = _e##name;

/* Declare according to SRAM/DRAM ranges in SoC hardware-defined address map. */
#define SRAM_START(addr) SYMBOL(sram, addr)

#define SRAM_END(addr) SYMBOL(esram, addr)

#define DRAM_START(addr) SYMBOL(dram, addr)

#define TIMESTAMP(addr, size) \
	REGION(timestamp, addr, size, 8) \
	_ = ASSERT(size >= 212, "Timestamp region must fit timestamp_cache!");

#define PRERAM_CBMEM_CONSOLE(addr, size) \
	REGION(preram_cbmem_console, addr, size, 4)

/* Use either CBFS_CACHE (unified) or both (PRERAM|POSTRAM)_CBFS_CACHE */
#define CBFS_CACHE(addr, size) \
	REGION(cbfs_cache, addr, size, 4) \
	ALIAS_REGION(cbfs_cache, preram_cbfs_cache) \
	ALIAS_REGION(cbfs_cache, postram_cbfs_cache)

#define FMAP_CACHE(addr, sz) \
	REGION(fmap_cache, addr, sz, 4) \
	_ = ASSERT(sz == 0 || sz >= FMAP_SIZE, \
		   STR(FMAP does not fit in FMAP_CACHE! (sz < FMAP_SIZE)));

#if ENV_ROMSTAGE_OR_BEFORE
	#define PRERAM_CBFS_CACHE(addr, size) \
		REGION(preram_cbfs_cache, addr, size, 4) \
		ALIAS_REGION(preram_cbfs_cache, cbfs_cache)
	#define POSTRAM_CBFS_CACHE(addr, size) \
		REGION(postram_cbfs_cache, addr, size, 4)
#else
	#define PRERAM_CBFS_CACHE(addr, size) \
		REGION(preram_cbfs_cache, addr, size, 4)
	#define POSTRAM_CBFS_CACHE(addr, size) \
		REGION(postram_cbfs_cache, addr, size, 4) \
		ALIAS_REGION(postram_cbfs_cache, cbfs_cache)
#endif

/* Careful: 'INCLUDE <filename>' must always be at the end of the output line */
#if ENV_DECOMPRESSOR
	#define DECOMPRESSOR(addr, sz) \
		SYMBOL(decompressor, addr) \
		_edecompressor = _decompressor + sz; \
		_ = ASSERT(_eprogram - _program <= sz, \
			STR(decompressor exceeded its allotted size! (sz))); \
		INCLUDE "decompressor/lib/program.ld"

	#define OVERLAP_DECOMPRESSOR_ROMSTAGE(addr, sz) DECOMPRESSOR(addr, sz)
	#define OVERLAP_DECOMPRESSOR_VERSTAGE_ROMSTAGE(addr, sz) \
		DECOMPRESSOR(addr, sz)
#else
	#define DECOMPRESSOR(addr, sz) \
		REGION(decompressor, addr, sz, 1)

	#define OVERLAP_DECOMPRESSOR_ROMSTAGE(addr, sz) ROMSTAGE(addr, sz)
	#define OVERLAP_DECOMPRESSOR_VERSTAGE_ROMSTAGE(addr, sz) \
		OVERLAP_VERSTAGE_ROMSTAGE(addr, sz)
#endif

#if ENV_BOOTBLOCK
	#define BOOTBLOCK(addr, sz) \
		SYMBOL(bootblock, addr) \
		_ebootblock = _bootblock + sz; \
		_ = ASSERT(_eprogram - _program <= sz, \
			STR(Bootblock exceeded its allotted size! (sz))); \
		INCLUDE "bootblock/lib/program.ld"
#else
	#define BOOTBLOCK(addr, sz) \
		REGION(bootblock, addr, sz, 1)
#endif

#if ENV_ROMSTAGE
	#define ROMSTAGE(addr, sz) \
		SYMBOL(romstage, addr) \
		_eromstage = _romstage + sz; \
		_ = ASSERT(_eprogram - _program <= sz, \
			STR(Romstage exceeded its allotted size! (sz))); \
		INCLUDE "romstage/lib/program.ld"
#else
	#define ROMSTAGE(addr, sz) \
		REGION(romstage, addr, sz, 1)
#endif

#if ENV_RAMSTAGE
	#define RAMSTAGE(addr, sz) \
		SYMBOL(ramstage, addr) \
		_eramstage = _ramstage + sz; \
		_ = ASSERT(_eprogram - _program <= sz, \
			STR(Ramstage exceeded its allotted size! (sz))); \
		INCLUDE "ramstage/lib/program.ld"
#else
	#define RAMSTAGE(addr, sz) \
		REGION(ramstage, addr, sz, 1)
#endif

/* VBOOT2_WORK must always use VB2_FIRMWARE_WORKBUF_RECOMMENDED_SIZE for its
 * size argument.  The constant is imported via vb2_workbuf_size.h.  */
#define VBOOT2_WORK(addr, sz) \
	REGION(vboot2_work, addr, sz, 16) \
	_ = ASSERT(sz == VB2_FIRMWARE_WORKBUF_RECOMMENDED_SIZE, \
		STR(vboot2 work buffer size must be equivalent to \
			VB2_FIRMWARE_WORKBUF_RECOMMENDED_SIZE! (sz)));

#define VBOOT2_TPM_LOG(addr, size) \
	REGION(vboot2_tpm_log, addr, size, 16) \
	_ = ASSERT(size >= 2K, "vboot2 tpm log buffer must be at least 2K!");

#if ENV_VERSTAGE
	#define VERSTAGE(addr, sz) \
		SYMBOL(verstage, addr) \
		_everstage = _verstage + sz; \
		_ = ASSERT(_eprogram - _program <= sz, \
			STR(Verstage exceeded its allotted size! (sz))); \
		INCLUDE "verstage/lib/program.ld"

	#define OVERLAP_VERSTAGE_ROMSTAGE(addr, size) \
		_ = ASSERT(CONFIG(VBOOT_RETURN_FROM_VERSTAGE) == 1, \
			"Must set RETURN_FROM_VERSTAGE to overlap romstage."); \
		VERSTAGE(addr, size)
#else
	#define VERSTAGE(addr, sz) \
		REGION(verstage, addr, sz, 1)

	#define OVERLAP_VERSTAGE_ROMSTAGE(addr, size) ROMSTAGE(addr, size)
#endif

#if ENV_POSTCAR
	#define POSTCAR(addr, sz) \
		SYMBOL(postcar, addr) \
		_epostcar = _postcar + sz; \
		_ = ASSERT(_eprogram - _program <= sz, \
			STR(Aftercar exceeded its allotted size! (sz))); \
		INCLUDE "postcar/lib/program.ld"
#else
	#define POSTCAR(addr, sz) \
		REGION(postcar, addr, sz, 1)
#endif

#define WATCHDOG_TOMBSTONE(addr, size) \
	REGION(watchdog_tombstone, addr, size, 4) \
	_ = ASSERT(size == 4, "watchdog tombstones should be exactly 4 byte!");

#endif /* __MEMLAYOUT_H */
