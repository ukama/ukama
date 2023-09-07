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

#ifndef DEVICE_PCI_TYPE_H
#define DEVICE_PCI_TYPE_H

#include <stdint.h>

typedef u32 pci_devfn_t;

/* Convert pci_devfn_t to offset in MMCONF space.
 * As it is one-to-one,  nothing needs to be done. */
#define PCI_DEVFN_OFFSET(x) ((x))

#define PCI_DEV(SEGBUS, DEV, FN) ( \
	(((SEGBUS) & 0xFFF) << 20) | \
	(((DEV) & 0x1F) << 15) | \
	(((FN)  & 0x07) << 12))

#define PCI_DEV_INVALID (0xffffffffU)

#if 1
/* FIXME: For most of the time in ramstage, we get valid device pointer
 * from calling the driver entry points. The assert should only be used
 * with searches like pcidev_behind(), and only if caller does not make
 * the check themselves.
 */
#define PCI_BDF(dev) pcidev_assert((dev))
#else
#define PCI_BDF(dev) pcidev_bdf((dev))
#endif

#endif /* DEVICE_PCI_TYPE_H */
