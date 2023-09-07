/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2019 Asami Doi <d0iasm.pub@gmail.com>
 *
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

#include <cbmem.h>
#include <ramdetect.h>
#include <symbols.h>

void *cbmem_top_chipset(void)
{
	return _dram + (probe_ramsize((uintptr_t)_dram, CONFIG_DRAM_SIZE_MB) * MiB);
}
