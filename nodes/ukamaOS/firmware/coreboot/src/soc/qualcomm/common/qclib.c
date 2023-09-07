/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018, The Linux Foundation.  All rights reserved.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2 and
 * only version 2 as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <console/cbmem_console.h>
#include <cbmem.h>
#include <boardid.h>
#include <string.h>
#include <fmap.h>
#include <assert.h>
#include <arch/mmu.h>
#include <cbfs.h>
#include <console/console.h>
#include <soc/mmu.h>
#include <soc/mmu_common.h>
#include <soc/qclib_common.h>
#include <soc/symbols_common.h>

struct qclib_cb_if_table qclib_cb_if_table = {
	.magic = QCLIB_MAGIC_NUMBER,
	.version = QCLIB_INTERFACE_VERSION,
	.num_entries = 0,
	.max_entries = QCLIB_MAX_NUMBER_OF_ENTRIES,
	.global_attributes = 0,
	.reserved = 0,
};

void qclib_add_if_table_entry(const char *name, void *base,
				uint32_t size, uint32_t attrs)
{
	struct qclib_cb_if_table_entry *te =
		&qclib_cb_if_table.te[qclib_cb_if_table.num_entries++];
	assert(qclib_cb_if_table.num_entries <= qclib_cb_if_table.max_entries);
	strncpy(te->name, name, sizeof(te->name) - 1);
	te->blob_address = (uintptr_t)base;
	te->size = size;
	te->blob_attributes = attrs;
}

static void write_ddr_information(struct qclib_cb_if_table_entry *te)
{
	uint64_t ddr_size;

	/* Save DDR info in SRAM region to share with ramstage */
	ddr_region->offset = te->blob_address;
	ddr_size = te->size;
	ddr_region->size = ddr_size * MiB;

	/* Use DDR info to configure MMU */
	qc_mmu_dram_config_post_dram_init((void *)ddr_region->offset,
		(size_t)ddr_region->size);
}

static void write_qclib_log_to_cbmemc(struct qclib_cb_if_table_entry *te)
{
	int i;
	char *ptr = (char *)te->blob_address;

	for (i = 0; i < te->size; i++)
		__cbmemc_tx_byte(*ptr++);
}

static void write_table_entry(struct qclib_cb_if_table_entry *te)
{

	if (!strncmp(QCLIB_TE_DDR_INFORMATION, te->name,
			sizeof(te->name))) {

		write_ddr_information(te);

	} else if (!strncmp(QCLIB_TE_DDR_TRAINING_DATA, te->name,
			sizeof(te->name))) {

		assert(fmap_overwrite_area(QCLIB_FR_DDR_TRAINING_DATA,
			(const void *)te->blob_address, te->size));

	} else if (!strncmp(QCLIB_TE_LIMITS_CFG_DATA, te->name,
			sizeof(te->name))) {

		assert(fmap_overwrite_area(QCLIB_FR_LIMITS_CFG_DATA,
			(const void *)te->blob_address, te->size));

	} else if (!strncmp(QCLIB_TE_QCLIB_LOG_BUFFER, te->name,
			sizeof(te->name))) {

		write_qclib_log_to_cbmemc(te);

	} else {

		printk(BIOS_WARNING, "%s write not implemented\n", te->name);
		printk(BIOS_WARNING, "  blob_address[%llx]..size[%x]\n",
			te->blob_address, te->size);

	}
}

static void dump_te_table(void)
{
	struct qclib_cb_if_table_entry *te;
	int i;

	for (i = 0; i < qclib_cb_if_table.num_entries; i++) {
		te = &qclib_cb_if_table.te[i];
		printk(BIOS_DEBUG, "[%s][%llx][%x][%x]\n",
			te->name, te->blob_address,
			te->size, te->blob_attributes);
	}
}

__weak int qclib_soc_blob_load(void) { return 0; }

void qclib_load_and_run(void)
{
	int i;
	ssize_t ssize;
	struct mmu_context pre_qclib_mmu_context;

	/* zero ddr_information SRAM region, needs new data each boot */
	memset(ddr_region, 0, sizeof(struct region));

	/* output area, QCLib copies console log buffer out */
	if (CONFIG(CONSOLE_CBMEM))
		qclib_add_if_table_entry(QCLIB_TE_QCLIB_LOG_BUFFER,
				_qclib_serial_log,
				REGION_SIZE(qclib_serial_log), 0);

	/* output area, QCLib fills in DDR details */
	qclib_add_if_table_entry(QCLIB_TE_DDR_INFORMATION, NULL, 0, 0);

	/* Attempt to load DDR Training Blob */
	ssize = fmap_read_area(QCLIB_FR_DDR_TRAINING_DATA, _ddr_training,
		REGION_SIZE(ddr_training));
	if (ssize < 0)
		goto fail;
	qclib_add_if_table_entry(QCLIB_TE_DDR_TRAINING_DATA,
				 _ddr_training, ssize, 0);

	/* hook for SoC specific binary blob loads */
	if (qclib_soc_blob_load()) {
		printk(BIOS_ERR, "qclib_soc_blob_load failed\n");
		goto fail;
	}

	/* Enable QCLib serial output, based on Kconfig */
	if (CONFIG(CONSOLE_SERIAL))
		qclib_cb_if_table.global_attributes =
			QCLIB_GA_ENABLE_UART_LOGGING;

	if (CONFIG(QC_SDI_ENABLE)) {
		struct prog qcsdi =
			PROG_INIT(PROG_REFCODE, CONFIG_CBFS_PREFIX "/qcsdi");

		/* Attempt to load QCSDI elf */
		if (prog_locate(&qcsdi))
			goto fail;

		if (cbfs_prog_stage_load(&qcsdi))
			goto fail;

		qclib_add_if_table_entry(QCLIB_TE_QCSDI, prog_entry(&qcsdi),
				   prog_size(&qcsdi), 0);
		printk(BIOS_INFO, "qcsdi.entry[%p]\n", qcsdi.entry);
	}

	dump_te_table();

	/* Attempt to load QCLib elf */
	struct prog qclib =
		PROG_INIT(PROG_REFCODE, CONFIG_CBFS_PREFIX "/qclib");

	if (prog_locate(&qclib))
		goto fail;

	if (cbfs_prog_stage_load(&qclib))
		goto fail;

	prog_set_entry(&qclib, prog_entry(&qclib), &qclib_cb_if_table);

	printk(BIOS_DEBUG, "\n\n\nQCLib is about to Initialize DDR\n");
	printk(BIOS_DEBUG, "Global Attributes[%x]..Table Entries Count[%d]\n",
		qclib_cb_if_table.global_attributes,
		qclib_cb_if_table.num_entries);
	printk(BIOS_DEBUG, "Jumping to QCLib code at %p(%p)\n",
		prog_entry(&qclib), prog_entry_arg(&qclib));

	/* back-up mmu context before disabling mmu and executing qclib */
	mmu_save_context(&pre_qclib_mmu_context);
	/* disable mmu before jumping to qclib. mmu_disable also
	   flushes and invalidates caches before disabling mmu. */
	mmu_disable();

	prog_run(&qclib);

	/* Before returning, QCLib flushes cache and disables mmu.
	   Explicitly disable mmu (flush, invalidate and disable mmu)
	   before re-enabling mmu with backed-up mmu context */
	mmu_disable();
	mmu_restore_context(&pre_qclib_mmu_context);
	mmu_enable();

	/* step through I/F table, handling return values */
	for (i = 0; i < qclib_cb_if_table.num_entries; i++)
		if (qclib_cb_if_table.te[i].blob_attributes &
				QCLIB_BA_SAVE_TO_STORAGE)
			write_table_entry(&qclib_cb_if_table.te[i]);

	/* confirm that we received valid ddr information from QCLib */
	assert((uintptr_t)_dram == region_offset(ddr_region) &&
		region_sz(ddr_region) >= (u8 *)cbmem_top() - _dram);

	printk(BIOS_DEBUG, "QCLib completed\n\n\n");

	return;

fail:
	die("Couldn't run QCLib.\n");
}
