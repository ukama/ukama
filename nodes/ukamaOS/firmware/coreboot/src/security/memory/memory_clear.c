/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2019 9elements Agency GmbH
 * Copyright (C) 2019 Facebook Inc.
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

#if CONFIG(ARCH_X86)
#include <cpu/x86/pae.h>
#else
#define memset_pae(a, b, c, d, e) 0
#define MEMSET_PAE_PGTL_ALIGN 0
#define MEMSET_PAE_PGTL_SIZE 0
#define MEMSET_PAE_PGTL_SIZE 0
#define MEMSET_PAE_VMEM_ALIGN 0
#endif

#include <memrange.h>
#include <bootmem.h>
#include <bootstate.h>
#include <symbols.h>
#include <console/console.h>
#include <arch/memory_clear.h>
#include <string.h>
#include <security/memory/memory.h>
#include <cbmem.h>
#include <arch/acpi.h>

/* Helper to find free space for memset_pae. */
static uintptr_t get_free_memory_range(struct memranges *mem,
				       const resource_t align,
				       const resource_t size)
{
	const struct range_entry *r;

	/* Find a spot for virtual memory address */
	memranges_each_entry(r, mem) {
		if (range_entry_tag(r) != BM_MEM_RAM)
			continue;

		if (ALIGN_UP(range_entry_base(r) + size, align) + size >
		    range_entry_end(r))
			continue;

		return ALIGN_UP(range_entry_base(r) + size, align);
	}
	printk(BIOS_ERR, "%s: Couldn't find free memory range\n", __func__);

	return 0;
}

/*
 * Clears all memory regions marked as BM_MEM_RAM.
 * Uses memset_pae if the memory region can't be accessed by memset and
 * architecture is x86.
 *
 * @return 0 on success, 1 on error
 */
static void clear_memory(void *unused)
{
	const struct range_entry *r;
	struct memranges mem;
	uintptr_t pgtbl, vmem_addr;

	if (acpi_is_wakeup_s3())
		return;

	if (!security_clear_dram_request())
		return;

	/* FSP1.0 is marked as MMIO and won't appear here */

	memranges_init(&mem, IORESOURCE_MEM | IORESOURCE_FIXED |
			IORESOURCE_STORED | IORESOURCE_ASSIGNED |
			IORESOURCE_CACHEABLE,
			IORESOURCE_MEM | IORESOURCE_FIXED |
			IORESOURCE_STORED | IORESOURCE_ASSIGNED |
			IORESOURCE_CACHEABLE,
			BM_MEM_RAM);

	/* Add reserved entries */
	void *baseptr = NULL;
	size_t size = 0;

	/* Only skip CBMEM, as RELOCATABLE_RAMSTAGE is a requirement, no need
	 * to separately protect stack or heap */

	cbmem_get_region(&baseptr, &size);
	memranges_insert(&mem, (uintptr_t)baseptr, size, BM_MEM_TABLE);

	if (CONFIG(PLATFORM_USES_FSP1_0)) {
		/* Protect CBMEM pointer */
		memranges_insert(&mem, CBMEM_FSP_HOB_PTR, sizeof(void *),
				 BM_MEM_TABLE);
	}

	if (CONFIG(ARCH_X86)) {
		/* Find space for PAE enabled memset */
		pgtbl = get_free_memory_range(&mem, MEMSET_PAE_PGTL_ALIGN,
					MEMSET_PAE_PGTL_SIZE);

		/* Don't touch page tables while clearing */
		memranges_insert(&mem, pgtbl, MEMSET_PAE_PGTL_SIZE,
					BM_MEM_TABLE);

		vmem_addr = get_free_memory_range(&mem, MEMSET_PAE_VMEM_ALIGN,
						MEMSET_PAE_PGTL_SIZE);

		printk(BIOS_SPEW, "%s: pgtbl at %p, virt memory at %p\n",
		__func__, (void *)pgtbl, (void *)vmem_addr);
	}

	/* Now clear all useable DRAM */
	memranges_each_entry(r, &mem) {
		if (range_entry_tag(r) != BM_MEM_RAM)
			continue;
		printk(BIOS_DEBUG, "%s: Clearing DRAM %016llx-%016llx\n",
		       __func__, range_entry_base(r), range_entry_end(r));

		/* Does regular memset work? */
		if (sizeof(resource_t) == sizeof(void *) ||
		    !(range_entry_end(r) >> (sizeof(void *) * 8))) {
			/* fastpath */
			memset((void *)(uintptr_t)range_entry_base(r), 0,
			       range_entry_size(r));
		}
		/* Use PAE if available */
		else if (CONFIG(ARCH_X86)) {
			if (memset_pae(range_entry_base(r), 0,
			    range_entry_size(r), (void *)pgtbl,
			    (void *)vmem_addr))
				printk(BIOS_ERR, "%s: Failed to memset "
				       "memory\n", __func__);
		} else {
			printk(BIOS_ERR, "%s: Failed to memset memory\n",
			       __func__);
		}
	}

	if (CONFIG(ARCH_X86)) {
		/* Clear previously skipped memory reserved for pagetables */
		printk(BIOS_DEBUG, "%s: Clearing DRAM %016lx-%016lx\n",
		__func__, pgtbl, pgtbl + MEMSET_PAE_PGTL_SIZE);

		memset((void *)pgtbl, 0, MEMSET_PAE_PGTL_SIZE);
	}

	memranges_teardown(&mem);
}

/* After DEV_INIT as MTRRs needs to be configured on x86 */
BOOT_STATE_INIT_ENTRY(BS_DEV_INIT, BS_ON_EXIT, clear_memory, NULL);
