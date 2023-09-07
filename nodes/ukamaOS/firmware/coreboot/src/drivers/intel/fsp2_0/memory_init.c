/*
 * This file is part of the coreboot project.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <security/vboot/antirollback.h>
#include <arch/symbols.h>
#include <assert.h>
#include <cbfs.h>
#include <cbmem.h>
#include <cf9_reset.h>
#include <console/console.h>
#include <elog.h>
#include <fsp/api.h>
#include <fsp/util.h>
#include <memrange.h>
#include <mrc_cache.h>
#include <program_loading.h>
#include <romstage_handoff.h>
#include <string.h>
#include <symbols.h>
#include <timestamp.h>
#include <security/vboot/vboot_common.h>
#include <security/tpm/tspi.h>
#include <vb2_api.h>
#include <fsp/memory_init.h>
#include <types.h>

static uint8_t temp_ram[CONFIG_FSP_TEMP_RAM_SIZE] __aligned(sizeof(uint64_t));

/* TPM MRC hash functionality depends on vboot starting before memory init. */
_Static_assert(!CONFIG(FSP2_0_USES_TPM_MRC_HASH) ||
	       CONFIG(VBOOT_STARTS_IN_BOOTBLOCK),
	       "for TPM MRC hash functionality, vboot must start in bootblock");

static void save_memory_training_data(bool s3wake, uint32_t fsp_version)
{
	size_t  mrc_data_size;
	const void *mrc_data;

	if (!CONFIG(CACHE_MRC_SETTINGS) || s3wake)
		return;

	mrc_data = fsp_find_nv_storage_data(&mrc_data_size);
	if (!mrc_data) {
		printk(BIOS_ERR, "Couldn't find memory training data HOB.\n");
		return;
	}

	/*
	 * Save MRC Data to CBMEM. By always saving the data this forces
	 * a retrain after a trip through Chrome OS recovery path. The
	 * code which saves the data to flash doesn't write if the latest
	 * training data matches this one.
	 */
	if (mrc_cache_stash_data(MRC_TRAINING_DATA, fsp_version, mrc_data,
				mrc_data_size) < 0)
		printk(BIOS_ERR, "Failed to stash MRC data\n");

	if (CONFIG(FSP2_0_USES_TPM_MRC_HASH))
		mrc_cache_update_hash(mrc_data, mrc_data_size);
}

static void do_fsp_post_memory_init(bool s3wake, uint32_t fsp_version)
{
	struct range_entry fsp_mem;

	fsp_find_reserved_memory(&fsp_mem);

	/* initialize cbmem by adding FSP reserved memory first thing */
	if (!s3wake) {
		cbmem_initialize_empty_id_size(CBMEM_ID_FSP_RESERVED_MEMORY,
			range_entry_size(&fsp_mem));
	} else if (cbmem_initialize_id_size(CBMEM_ID_FSP_RESERVED_MEMORY,
				range_entry_size(&fsp_mem))) {
		if (CONFIG(HAVE_ACPI_RESUME)) {
			printk(BIOS_ERR,
				"Failed to recover CBMEM in S3 resume.\n");
			/* Failed S3 resume, reset to come up cleanly */
			/* FIXME: A "system" reset is likely enough: */
			full_reset();
		}
	}

	/* make sure FSP memory is reserved in cbmem */
	if (range_entry_base(&fsp_mem) !=
		(uintptr_t)cbmem_find(CBMEM_ID_FSP_RESERVED_MEMORY))
		die("Failed to accommodate FSP reserved memory request!\n");

	save_memory_training_data(s3wake, fsp_version);

	/* Create romstage handof information */
	romstage_handoff_init(s3wake);
}

static void fsp_fill_mrc_cache(FSPM_ARCH_UPD *arch_upd, uint32_t fsp_version)
{
	struct region_device rdev;
	void *data;

	arch_upd->NvsBufferPtr = NULL;

	if (!CONFIG(CACHE_MRC_SETTINGS))
		return;

	/*
	 * In recovery mode, force retraining:
	 * 1. Recovery cache is not supported, or
	 * 2. Memory retrain switch is set.
	 */
	if (vboot_recovery_mode_enabled()) {
		if (!CONFIG(HAS_RECOVERY_MRC_CACHE))
			return;
		if (vboot_recovery_mode_memory_retrain())
			return;
	}

	if (mrc_cache_get_current(MRC_TRAINING_DATA, fsp_version, &rdev) < 0)
		return;

	/* Assume boot device is memory mapped. */
	assert(CONFIG(BOOT_DEVICE_MEMORY_MAPPED));
	data = rdev_mmap_full(&rdev);

	if (data == NULL)
		return;

	if (CONFIG(FSP2_0_USES_TPM_MRC_HASH) &&
	    !mrc_cache_verify_hash(data, region_device_sz(&rdev)))
		return;

	/* MRC cache found */
	arch_upd->NvsBufferPtr = data;

	printk(BIOS_SPEW, "MRC cache found, size %zx\n",
			region_device_sz(&rdev));
}

static enum cb_err check_region_overlap(const struct memranges *ranges,
					const char *description,
					uintptr_t begin, uintptr_t end)
{
	const struct range_entry *r;

	memranges_each_entry(r, ranges) {
		if (end <= range_entry_base(r))
			continue;
		if (begin >= range_entry_end(r))
			continue;
		printk(BIOS_CRIT, "'%s' overlaps currently running program: "
			"[%p, %p)\n", description, (void *)begin, (void *)end);
		return CB_ERR;
	}

	return CB_SUCCESS;
}

static enum cb_err setup_fsp_stack_frame(FSPM_ARCH_UPD *arch_upd,
		const struct memranges *memmap)
{
	uintptr_t stack_begin;
	uintptr_t stack_end;

	/*
	 * FSPM_UPD passed here is populated with default values
	 * provided by the blob itself. We let FSPM use top of CAR
	 * region of the size it requests.
	 */
	stack_end = (uintptr_t)_car_region_end;
	stack_begin = stack_end - arch_upd->StackSize;
	if (check_region_overlap(memmap, "FSPM stack", stack_begin,
				stack_end) != CB_SUCCESS)
		return CB_ERR;

	arch_upd->StackBase = (void *)stack_begin;
	return CB_SUCCESS;
}

static enum cb_err fsp_fill_common_arch_params(FSPM_ARCH_UPD *arch_upd,
					bool s3wake, uint32_t fsp_version,
					const struct memranges *memmap)
{
	/*
	 * FSP 2.1 version would use same stack as coreboot instead of
	 * setting up separate stack frame. FSP 2.1 would not relocate stack
	 * top and does not reinitialize stack pointer. The parameters passed
	 * as StackBase and StackSize are actually for temporary RAM and HOBs
	 * and are not related to FSP stack at all.
	 */
	if (CONFIG(FSP_USES_CB_STACK)) {
		arch_upd->StackBase = temp_ram;
		arch_upd->StackSize = sizeof(temp_ram);
	} else if (setup_fsp_stack_frame(arch_upd, memmap)) {
		return CB_ERR;
	}

	fsp_fill_mrc_cache(arch_upd, fsp_version);

	/* Configure bootmode */
	if (s3wake) {
		/*
		 * For S3 resume case, if valid mrc cache data is not found or
		 * RECOVERY_MRC_CACHE hash verification fails, the S3 data
		 * pointer would be null and S3 resume fails with fsp-m
		 * returning error. Invoking a reset here saves time.
		 */
		if (!arch_upd->NvsBufferPtr)
			/* FIXME: A "system" reset is likely enough: */
			full_reset();
		arch_upd->BootMode = FSP_BOOT_ON_S3_RESUME;
	} else {
		if (arch_upd->NvsBufferPtr)
			arch_upd->BootMode =
				FSP_BOOT_ASSUMING_NO_CONFIGURATION_CHANGES;
		else
			arch_upd->BootMode = FSP_BOOT_WITH_FULL_CONFIGURATION;
	}

	printk(BIOS_SPEW, "bootmode is set to: %d\n", arch_upd->BootMode);

	return CB_SUCCESS;
}

__weak
uint8_t fsp_memory_mainboard_version(void)
{
	return 0;
}

__weak
uint8_t fsp_memory_soc_version(void)
{
	return 0;
}

/*
 * Allow SoC and/or mainboard to bump the revision of the FSP setting
 * number. The FSP spec uses the low 8 bits as the build number. Take over
 * bits 3:0 for the SoC setting and bits 7:4 for the mainboard. That way
 * a tweak in the settings will bump the version used to track the cached
 * setting which triggers retraining when the FSP version hasn't changed, but
 * the SoC or mainboard settings have.
 */
static uint32_t fsp_memory_settings_version(const struct fsp_header *hdr)
{
	/* Use the full FSP version by default. */
	uint32_t ver = hdr->fsp_revision;

	if (!CONFIG(FSP_PLATFORM_MEMORY_SETTINGS_VERSIONS))
		return ver;

	ver &= ~0xff;
	ver |= (0xf & fsp_memory_mainboard_version()) << 4;
	ver |= (0xf & fsp_memory_soc_version()) << 0;

	return ver;
}

static void do_fsp_memory_init(struct fsp_header *hdr, bool s3wake,
					const struct memranges *memmap)
{
	uint32_t status;
	fsp_memory_init_fn fsp_raminit;
	FSPM_UPD fspm_upd, *upd;
	FSPM_ARCH_UPD *arch_upd;
	uint32_t fsp_version;

	post_code(POST_MEM_PREINIT_PREP_START);

	fsp_version = fsp_memory_settings_version(hdr);

	upd = (FSPM_UPD *)(hdr->cfg_region_offset + hdr->image_base);

	if (upd->FspUpdHeader.Signature != FSPM_UPD_SIGNATURE)
		die_with_post_code(POST_INVALID_VENDOR_BINARY,
			"Invalid FSPM signature!\n");

	/* Copy the default values from the UPD area */
	memcpy(&fspm_upd, upd, sizeof(fspm_upd));

	arch_upd = &fspm_upd.FspmArchUpd;

	/* Reserve enough memory under TOLUD to save CBMEM header */
	arch_upd->BootLoaderTolumSize = cbmem_overhead_size();

	/* Fill common settings on behalf of chipset. */
	if (fsp_fill_common_arch_params(arch_upd, s3wake, fsp_version,
					memmap) != CB_SUCCESS)
		die_with_post_code(POST_INVALID_VENDOR_BINARY,
			"FSPM_ARCH_UPD not found!\n");

	/* Give SoC and mainboard a chance to update the UPD */
	platform_fsp_memory_init_params_cb(&fspm_upd, fsp_version);

	if (CONFIG(MMA))
		setup_mma(&fspm_upd.FspmConfig);

	post_code(POST_MEM_PREINIT_PREP_END);

	/* Call FspMemoryInit */
	fsp_raminit = (void *)(hdr->image_base + hdr->memory_init_entry_offset);
	fsp_debug_before_memory_init(fsp_raminit, upd, &fspm_upd);

	post_code(POST_FSP_MEMORY_INIT);
	timestamp_add_now(TS_FSP_MEMORY_INIT_START);
	status = fsp_raminit(&fspm_upd, fsp_get_hob_list_ptr());
	post_code(POST_FSP_MEMORY_EXIT);
	timestamp_add_now(TS_FSP_MEMORY_INIT_END);

	/* Handle any errors returned by FspMemoryInit */
	fsp_handle_reset(status);
	if (status != FSP_SUCCESS) {
		printk(BIOS_CRIT, "FspMemoryInit returned 0x%08x\n", status);
		die_with_post_code(POST_RAM_FAILURE,
			"FspMemoryInit returned an error!\n");
	}

	do_fsp_post_memory_init(s3wake, fsp_version);

	/*
	 * fsp_debug_after_memory_init() checks whether the end of the tolum
	 * region is the same as the top of cbmem, so must be called here
	 * after cbmem has been initialised in do_fsp_post_memory_init().
	 */
	fsp_debug_after_memory_init(status);
}

/* Load the binary into the memory specified by the info header. */
static enum cb_err load_fspm_mem(struct fsp_header *hdr,
					const struct region_device *rdev,
					const struct memranges *memmap)
{
	uintptr_t fspm_begin;
	uintptr_t fspm_end;

	if (fsp_validate_component(hdr, rdev) != CB_SUCCESS)
		return CB_ERR;

	fspm_begin = hdr->image_base;
	fspm_end = fspm_begin + hdr->image_size;

	if (check_region_overlap(memmap, "FSPM", fspm_begin, fspm_end) !=
		CB_SUCCESS)
		return CB_ERR;

	/* Load binary into memory at provided address. */
	if (rdev_readat(rdev, (void *)fspm_begin, 0, fspm_end - fspm_begin) < 0)
		return CB_ERR;

	return CB_SUCCESS;
}

/* Handle the case when FSPM is running XIP. */
static enum cb_err load_fspm_xip(struct fsp_header *hdr,
					const struct region_device *rdev)
{
	void *base;

	if (fsp_validate_component(hdr, rdev) != CB_SUCCESS)
		return CB_ERR;

	base = rdev_mmap_full(rdev);
	if ((uintptr_t)base != hdr->image_base) {
		printk(BIOS_CRIT, "FSPM XIP base does not match: %p vs %p\n",
			(void *)(uintptr_t)hdr->image_base, base);
		return CB_ERR;
	}

	/*
	 * Since the component is XIP it's already in the address space. Thus,
	 * there's no need to rdev_munmap().
	 */
	return CB_SUCCESS;
}

void fsp_memory_init(bool s3wake)
{
	struct fsp_header hdr;
	enum cb_err status;
	struct cbfsf file_desc;
	struct region_device file_data;
	const char *name = CONFIG_FSP_M_CBFS;
	struct memranges memmap;
	struct range_entry prog_ranges[2];

	elog_boot_notify(s3wake);

	if (cbfs_boot_locate(&file_desc, name, NULL)) {
		printk(BIOS_CRIT, "Could not locate %s in CBFS\n", name);
		die("FSPM not available!\n");
	}

	cbfs_file_data(&file_data, &file_desc);

	/* Build up memory map of romstage address space including CAR. */
	memranges_init_empty(&memmap, &prog_ranges[0], ARRAY_SIZE(prog_ranges));
	if (ENV_CACHE_AS_RAM)
		memranges_insert(&memmap, (uintptr_t)_car_region_start,
			_car_unallocated_start - _car_region_start, 0);
	memranges_insert(&memmap, (uintptr_t)_program, REGION_SIZE(program), 0);

	if (!CONFIG(FSP_M_XIP))
		status = load_fspm_mem(&hdr, &file_data, &memmap);
	else
		status = load_fspm_xip(&hdr, &file_data);

	if (status != CB_SUCCESS)
		die("Loading FSPM failed!\n");

	/* Signal that FSP component has been loaded. */
	prog_segment_loaded(hdr.image_base, hdr.image_size, SEG_FINAL);

	timestamp_add_now(TS_BEFORE_INITRAM);

	do_fsp_memory_init(&hdr, s3wake, &memmap);

	timestamp_add_now(TS_AFTER_INITRAM);
}
