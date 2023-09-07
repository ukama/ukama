/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018 Intel Corporation.
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

#include <arch/acpi.h>
#include <cbmem.h>
#include <commonlib/storage/sd_mmc.h>
#include <commonlib/sd_mmc_ctrlr.h>
#include <commonlib/sdhci.h>
#include <compiler.h>
#include <console/console.h>
#include <device/pci.h>
#include <intelblocks/mmc.h>
#include <soc/iomap.h>
#include <soc/pci_devs.h>
#include <string.h>

void soc_sd_mmc_controller_quirks(struct sd_mmc_ctrlr *ctrlr)
{
	uint32_t f_min, f_max;

	if (soc_get_mmc_frequencies(&f_min, &f_max) < 0) {
		printk(BIOS_ERR,
			"MMC early init: failed to get mmc frequencies\n");
		return;
	}

	ctrlr->f_min = f_min;
	ctrlr->f_max = f_max;
}

static void enable_mmc_controller_bar(void)
{
	pci_write_config32(PCH_DEV_EMMC, PCI_BASE_ADDRESS_0,
				PRERAM_MMC_BASE_ADDRESS);
	pci_write_config32(PCH_DEV_EMMC, PCI_COMMAND,
				PCI_COMMAND_MASTER | PCI_COMMAND_MEMORY);
}

static void disable_mmc_controller_bar(void)
{
	pci_write_config32(PCH_DEV_EMMC, PCI_BASE_ADDRESS_0, 0);
	pci_write_config32(PCH_DEV_EMMC, PCI_COMMAND,
				~(PCI_COMMAND_MASTER | PCI_COMMAND_MEMORY));
}

static void set_early_mmc_wake_status(int32_t status)
{
	int32_t *ms_cbmem;

	ms_cbmem = cbmem_add(CBMEM_ID_MMC_STATUS, sizeof(int));

	if (ms_cbmem == NULL) {
		printk(BIOS_ERR,
			"%s: Failed to add early mmc wake status to cbmem!\n",
			__func__);
		return;
	}

	*ms_cbmem = status;
}

int early_mmc_wake_hw(void)
{
	struct storage_media media;
	struct sd_mmc_ctrlr *mmc_ctrlr;
	struct sdhci_ctrlr *sdhci_ctrlr;
	int err;

	if (acpi_is_wakeup_s3())
		return -1;

	/* Configure mmc gpios */
	if (soc_configure_mmc_gpios() < 0) {
		printk(BIOS_ERR,
			"%s: MMC early init: failed to configure mmc gpios\n",
			__func__);
		return -1;
	}
	/* Setup pci bar */
	enable_mmc_controller_bar();

	/* Initialize sdhci */
	mmc_ctrlr = new_pci_sdhci_controller(PCH_DEV_EMMC);
	if (mmc_ctrlr == NULL)
		goto out_err;

	sdhci_ctrlr = container_of(mmc_ctrlr, struct sdhci_ctrlr, sd_mmc_ctrlr);

	/* set emmc DLL tuning parameters */
	if (set_mmc_dll(sdhci_ctrlr->ioaddr) < 0)
		goto out_err;

	memset(&media, 0, sizeof(media));
	media.ctrlr = mmc_ctrlr;
	SET_BUS_WIDTH(mmc_ctrlr, 1);
	/*
	 * Set clock to 1 so that the driver can choose minimum frequency
	 * possible
	 */
	SET_CLOCK(mmc_ctrlr, 1);

	/* Reset emmc, send CMD0 */
	if (sd_mmc_go_idle(&media))
		goto out_err;

	/* Send CMD1 */
	err = mmc_send_op_cond(&media);
	if (err != 0 && err != CARD_IN_PROGRESS)
		goto out_err;

	disable_mmc_controller_bar();

	set_early_mmc_wake_status(1);
	return 0;

out_err:

	disable_mmc_controller_bar();
	return -1;
}
