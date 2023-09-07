/*
 * This file is part of the coreboot project.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation; either version 2 of
 * the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but without any warranty; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <arch/early_variables.h>
#include <commonlib/sdhci.h>
#include <device/pci.h>
#include <device/pci_ops.h>
#include <stdint.h>
#include <string.h>

#include "sd_mmc.h"
#include "storage.h"

/* Initialize an SDHCI port */
int sdhci_controller_init(struct sdhci_ctrlr *sdhci_ctrlr, void *ioaddr)
{
	memset(sdhci_ctrlr, 0, sizeof(*sdhci_ctrlr));
	sdhci_ctrlr->ioaddr = ioaddr;
	return add_sdhci(sdhci_ctrlr);
}

struct sd_mmc_ctrlr *new_mem_sdhci_controller(void *ioaddr)
{
	static bool sdhci_init_done CAR_GLOBAL;
	static struct sdhci_ctrlr sdhci_ctrlr CAR_GLOBAL;

	if (car_get_var(sdhci_init_done) == true) {
		sdhc_error("Error: SDHCI is already initialized.\n");
		return NULL;
	}

	if (sdhci_controller_init(car_get_var_ptr(&sdhci_ctrlr), ioaddr)) {
		sdhc_error("Error: SDHCI initialization failed.\n");
		return NULL;
	}

	car_set_var(sdhci_init_done, true);

	return car_get_var_ptr(&sdhci_ctrlr.sd_mmc_ctrlr);
}

struct sd_mmc_ctrlr *new_pci_sdhci_controller(pci_devfn_t dev)
{
	uint32_t addr;

	addr = pci_s_read_config32(dev, PCI_BASE_ADDRESS_0);
	if (addr == ((uint32_t)~0)) {
		sdhc_error("Error: PCI SDHCI not found\n");
		return NULL;
	}

	addr &= ~0xf;
	return new_mem_sdhci_controller((void *)addr);
}
