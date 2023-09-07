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

#include <commonlib/helpers.h>
#include <console/console.h>
#include <program_loading.h>
#include <ip_checksum.h>
#include <symbols.h>

int payload_arch_usable_ram_quirk(uint64_t start, uint64_t size)
{
	if (start < 1 * MiB && (start + size) <= 1 * MiB) {
		printk(BIOS_DEBUG,
			"Payload being loaded at below 1MiB without region being marked as RAM usable.\n");
		return 1;
	}

	return 0;
}

void arch_prog_run(struct prog *prog)
{
#ifdef __x86_64__
	void (*doit)(void *arg);
#else
	/* Ensure the argument is pushed on the stack. */
	asmlinkage void (*doit)(void *arg);
#endif
	doit = prog_entry(prog);
	doit(prog_entry_arg(prog));
}
