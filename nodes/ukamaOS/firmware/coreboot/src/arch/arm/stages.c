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

/*
 * This file contains entry/exit functions for each stage during coreboot
 * execution (bootblock entry and ramstage exit will depend on external
 * loading).
 *
 * Entry points should be set in the linker script and honored by CBFS,
 * so text section layout shouldn't matter. Still, it doesn't hurt to put
 * stage_entry first (which XXXstage.ld will do automatically through the
 * .text.stage_entry section created by -ffunction-sections).
 */

#include <cbmem.h>
#include <arch/stages.h>
#include <arch/cache.h>

/**
 * generic stage entry point. override this if board specific code is needed.
 */
__weak void stage_entry(uintptr_t stage_arg)
{
	if (!ENV_ROMSTAGE_OR_BEFORE)
		_cbmem_top_ptr = stage_arg;
	main();
}
