/*
 * This file is part of the coreboot project.
 *
 * Copyright 2013 Google Inc.
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

#include <program_loading.h>
#include <vm.h>
#include <arch/boot.h>
#include <arch/encoding.h>
#include <arch/smp/smp.h>
#include <mcall.h>
#include <commonlib/cbfs_serialized.h>
#include <console/console.h>

struct arch_prog_run_args {
	struct prog *prog;
	struct prog *opensbi;
};

/*
 * A pointer to the Flattened Device Tree passed to coreboot by the boot ROM.
 * Presumably this FDT is also in ROM.
 *
 * This pointer is only used in ramstage!
 */

static void do_arch_prog_run(struct arch_prog_run_args *args)
{
	int hart_id = HLS()->hart_id;
	struct prog *prog = args->prog;
	void *fdt = HLS()->fdt;

	if (prog_cbfs_type(prog) == CBFS_TYPE_FIT)
		fdt = prog_entry_arg(prog);

	if (ENV_RAMSTAGE && prog_type(prog) == PROG_PAYLOAD) {
		if (CONFIG(RISCV_OPENSBI))
			run_payload_opensbi(prog, fdt, args->opensbi, RISCV_PAYLOAD_MODE_S);
		else
			run_payload(prog, fdt, RISCV_PAYLOAD_MODE_S);
	} else {
		void (*doit)(int hart_id, void *fdt, void *arg) = prog_entry(prog);
		doit(hart_id, fdt, prog_entry_arg(prog));
	}

	die("Failed to run stage");
}

void arch_prog_run(struct prog *prog)
{
	struct arch_prog_run_args args = {};

	args.prog = prog;

	/* In case of OpenSBI we have to load it before resuming all HARTs */
	if (ENV_RAMSTAGE && CONFIG(RISCV_OPENSBI)) {
		struct prog sbi = PROG_INIT(PROG_OPENSBI, CONFIG_CBFS_PREFIX"/opensbi");

		if (prog_locate(&sbi))
			die("OpenSBI not found");

		if (!selfload_check(&sbi, BM_MEM_OPENSBI))
			die("OpenSBI load failed");

		args.opensbi = &sbi;
	}

	smp_resume((void (*)(void *))do_arch_prog_run, &args);
}
