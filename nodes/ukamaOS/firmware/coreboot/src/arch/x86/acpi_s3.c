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

#include <console/console.h>
#include <string.h>
#include <arch/acpi.h>
#include <arch/cpu.h>
#include <cbmem.h>
#include <commonlib/helpers.h>
#include <fallback.h>
#include <timestamp.h>
#include <program_loading.h>
#include <romstage_handoff.h>
#include <symbols.h>
#include <cpu/x86/smm.h>

#if ENV_RAMSTAGE || ENV_POSTCAR

/* This is filled with acpi_is_wakeup() call early in ramstage. */
static int acpi_slp_type = -1;

static void acpi_handoff_wakeup(void)
{
	if (acpi_slp_type < 0) {
		if (romstage_handoff_is_resume()) {
			printk(BIOS_DEBUG, "S3 Resume.\n");
			acpi_slp_type = ACPI_S3;
		} else {
			printk(BIOS_DEBUG, "Normal boot.\n");
			acpi_slp_type = ACPI_S0;
		}
	}
}

int acpi_is_wakeup(void)
{
	acpi_handoff_wakeup();
	/* Both resume from S2 and resume from S3 restart at CPU reset */
	return (acpi_slp_type == ACPI_S3 || acpi_slp_type == ACPI_S2);
}

int acpi_is_wakeup_s3(void)
{
	acpi_handoff_wakeup();
	return (acpi_slp_type == ACPI_S3);
}

int acpi_is_wakeup_s4(void)
{
	acpi_handoff_wakeup();
	return (acpi_slp_type == ACPI_S4);
}
#endif /* ENV_RAMSTAGE */

#define WAKEUP_BASE 0x600

asmlinkage void (*acpi_do_wakeup)(uintptr_t vector) = (void *)WAKEUP_BASE;

extern unsigned char __wakeup;
extern unsigned int __wakeup_size;

static void acpi_jump_to_wakeup(void *vector)
{
	if (!acpi_s3_resume_allowed()) {
		printk(BIOS_WARNING, "ACPI: S3 resume not allowed.\n");
		return;
	}

	/* Copy wakeup trampoline in place. */
	memcpy((void *)WAKEUP_BASE, &__wakeup, __wakeup_size);

	set_boot_successful();

	timestamp_add_now(TS_ACPI_WAKE_JUMP);

	acpi_do_wakeup((uintptr_t)vector);
}

void __weak mainboard_suspend_resume(void)
{
}

void acpi_resume(void *wake_vec)
{
	if (CONFIG(HAVE_SMI_HANDLER)) {
		void *gnvs_address = cbmem_find(CBMEM_ID_ACPI_GNVS);

		/* Restore GNVS pointer in SMM if found */
		if (gnvs_address) {
			printk(BIOS_DEBUG, "Restore GNVS pointer to %p\n",
			       gnvs_address);
			smm_setup_structures(gnvs_address, NULL, NULL);
		}
	}

	/* Call mainboard resume handler first, if defined. */
	mainboard_suspend_resume();

	post_code(POST_OS_RESUME);
	acpi_jump_to_wakeup(wake_vec);
}
