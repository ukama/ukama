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

#include <cbmem.h>
#include <arch/cache.h>
#include <arch/lib_helpers.h>
#include <arch/stages.h>
#include <arch/transition.h>
#include <bl31.h>
#include <program_loading.h>

static void run_payload(struct prog *prog)
{
	void (*doit)(void *);
	void *arg;

	doit = prog_entry(prog);
	arg = prog_entry_arg(prog);
	u64 payload_spsr = get_eret_el(EL2, SPSR_USE_L);

	if (CONFIG(ARM64_USE_ARM_TRUSTED_FIRMWARE))
		run_bl31((u64)doit, (u64)arg, payload_spsr);
	else
		transition_to_el2(doit, arg, payload_spsr);
}

void arch_prog_run(struct prog *prog)
{
	void (*doit)(void *);

	if (ENV_RAMSTAGE && prog_type(prog) == PROG_PAYLOAD) {
		run_payload(prog);
		return;
	}

	doit = prog_entry(prog);

	doit(prog_entry_arg(prog));
}

/* Generic stage entry point. Can be overridden by board/SoC if needed. */
__weak void stage_entry(uintptr_t stage_arg)
{
	if (!ENV_ROMSTAGE_OR_BEFORE)
		_cbmem_top_ptr = stage_arg;

	main();
}
