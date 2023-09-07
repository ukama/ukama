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

#include <stdint.h>
#include <security/intel/txt/txt.h>
#include "memory.h"

/**
 * To be called after DRAM init.
 * Tells the caller if DRAM must be cleared as requested by the user,
 * firmware or security framework.
 */
bool security_clear_dram_request(void)
{
	if (CONFIG(SECURITY_CLEAR_DRAM_ON_REGULAR_BOOT))
		return true;

	if (CONFIG(INTEL_TXT) && intel_txt_memory_has_secrets())
		return true;

	/* TODO: Add TEE environments here */

	return false;
}
