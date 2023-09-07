/*
 * This file is part of the coreboot project.
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

#include <bootmem.h>
#include <boot/tables.h>
#include <boot/coreboot_tables.h>
#include <symbols.h>

DECLARE_OPTIONAL_REGION(bl31);

void arch_write_tables(uintptr_t coreboot_table)
{
}

void bootmem_arch_add_ranges(void)
{
	bootmem_add_range((uintptr_t)_ttb, REGION_SIZE(ttb), BM_MEM_RAMSTAGE);

	if (CONFIG(ARM64_USE_ARM_TRUSTED_FIRMWARE) &&
	    REGION_SIZE(bl31) > 0)
		bootmem_add_range((uintptr_t)_bl31, REGION_SIZE(bl31),
				  BM_MEM_BL31);

	if (!CONFIG(COMMON_CBFS_SPI_WRAPPER))
		return;
	bootmem_add_range((uintptr_t)_postram_cbfs_cache,
			  REGION_SIZE(postram_cbfs_cache), BM_MEM_RAMSTAGE);
}

void lb_arch_add_records(struct lb_header *header)
{
}
