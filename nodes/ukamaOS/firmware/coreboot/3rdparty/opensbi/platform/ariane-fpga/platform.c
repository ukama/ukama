/* SPDX-License-Identifier: GPL-2.0 */
/*
 * Copyright (C) 2019 FORTH-ICS/CARV
 *				Panagiotis Peristerakis <perister@ics.forth.gr>
 */

#include <sbi/riscv_encoding.h>
#include <sbi/sbi_const.h>
#include <sbi/sbi_platform.h>
#include <sbi_utils/irqchip/plic.h>
#include <sbi_utils/serial/uart8250.h>
#include <sbi_utils/sys/clint.h>
#include <sbi/sbi_console.h>
#include <sbi/sbi_hart.h>
#include <libfdt.h>
#include <fdt.h>
#include <sbi_utils/irqchip/plic.h>
#include <sbi/riscv_io.h>

#define ARIANE_UART_ADDR			0x10000000
#define ARIANE_UART_FREQ			50000000
#define ARIANE_UART_BAUDRATE			115200
#define ARIANE_UART_REG_SHIFT			2
#define ARIANE_UART_REG_WIDTH			4
#define ARIANE_PLIC_ADDR			0xc000000
#define ARIANE_PLIC_NUM_SOURCES			3
#define ARIANE_HART_COUNT			1
#define ARIANE_CLINT_ADDR 0x2000000
#define PLIC_ENABLE_BASE		0x2000
#define PLIC_ENABLE_STRIDE		0x80
#define PLIC_CONTEXT_BASE		0x200000
#define PLIC_CONTEXT_STRIDE		0x1000

#define SBI_ARIANE_FEATURES	\
	(SBI_PLATFORM_HAS_TIMER_VALUE | \
	 SBI_PLATFORM_HAS_SCOUNTEREN | \
	 SBI_PLATFORM_HAS_MCOUNTEREN | \
	 SBI_PLATFORM_HAS_MFAULTS_DELEGATION)


/*
 * Ariane platform early initialization.
 */
static int ariane_early_init(bool cold_boot)
{
	/* For now nothing to do. */
	return 0;
}

/*
 * Ariane platform final initialization.
 */
static int ariane_final_init(bool cold_boot)
{
	void *fdt;

	if (!cold_boot)
		return 0;
	fdt = sbi_scratch_thishart_arg1_ptr();
	plic_fdt_fixup(fdt, "riscv,plic0");
	return 0;
}

/*
 * Initialize the ariane console.
 */
static int ariane_console_init(void)
{
	return uart8250_init(ARIANE_UART_ADDR,
						 ARIANE_UART_FREQ,
						 ARIANE_UART_BAUDRATE,
						ARIANE_UART_REG_SHIFT,
						ARIANE_UART_REG_WIDTH);
}

static int plic_ariane_warm_irqchip_init(u32 target_hart,
			   int m_cntx_id, int s_cntx_id)
{
	size_t i, ie_words = ARIANE_PLIC_NUM_SOURCES / 32 + 1;

	if (ARIANE_HART_COUNT <= target_hart)
		return -1;
	/* By default, enable all IRQs for M-mode of target HART */
	if (m_cntx_id > -1) {
		for (i = 0; i < ie_words; i++)
			plic_set_ie(m_cntx_id, i, 1);
	}
	/* Enable all IRQs for S-mode of target HART */
	if (s_cntx_id > -1) {
		for (i = 0; i < ie_words; i++)
			plic_set_ie(s_cntx_id, i, 1);
	}
	/* By default, enable M-mode threshold */
	if (m_cntx_id > -1)
		plic_set_thresh(m_cntx_id, 1);
	/* By default, disable S-mode threshold */
	if (s_cntx_id > -1)
		plic_set_thresh(s_cntx_id, 0);

	return 0;
}

/*
 * Initialize the ariane interrupt controller for current HART.
 */
static int ariane_irqchip_init(bool cold_boot)
{
	u32 hartid = sbi_current_hartid();
	int ret;

	if (cold_boot) {
		ret = plic_cold_irqchip_init(ARIANE_PLIC_ADDR,
					     ARIANE_PLIC_NUM_SOURCES,
					     ARIANE_HART_COUNT);
		if (ret)
			return ret;
	}
	return plic_ariane_warm_irqchip_init(hartid,
					2 * hartid, 2 * hartid + 1);
}

/*
 * Initialize IPI for current HART.
 */
static int ariane_ipi_init(bool cold_boot)
{
	int ret;

	if (cold_boot) {
		ret = clint_cold_ipi_init(ARIANE_CLINT_ADDR,
					  ARIANE_HART_COUNT);
		if (ret)
			return ret;
	}

	return clint_warm_ipi_init();
}

/*
 * Initialize ariane timer for current HART.
 */
static int ariane_timer_init(bool cold_boot)
{
	int ret;

	if (cold_boot) {
		ret = clint_cold_timer_init(ARIANE_CLINT_ADDR,
							ARIANE_HART_COUNT);
		if (ret)
			return ret;
	}

	return clint_warm_timer_init();
}

/*
 * Reboot the ariane.
 */
static int ariane_system_reboot(u32 type)
{
	/* For now nothing to do. */
	sbi_printf("System reboot\n");
	return 0;
}

/*
 * Shutdown or poweroff the ariane.
 */
static int ariane_system_shutdown(u32 type)
{
	/* For now nothing to do. */
	sbi_printf("System shutdown\n");
	return 0;
}

/*
 * Platform descriptor.
 */
const struct sbi_platform_operations platform_ops = {
	.early_init = ariane_early_init,
	.final_init = ariane_final_init,
	.console_init = ariane_console_init,
	.console_putc = uart8250_putc,
	.console_getc = uart8250_getc,
	.irqchip_init = ariane_irqchip_init,
	.ipi_init = ariane_ipi_init,
	.ipi_send = clint_ipi_send,
	.ipi_clear = clint_ipi_clear,
	.timer_init = ariane_timer_init,
	.timer_value = clint_timer_value,
	.timer_event_start = clint_timer_event_start,
	.timer_event_stop = clint_timer_event_stop,
	.system_reboot = ariane_system_reboot,
	.system_shutdown = ariane_system_shutdown
};

const struct sbi_platform platform = {
	.opensbi_version = OPENSBI_VERSION,
	.platform_version = SBI_PLATFORM_VERSION(0x0, 0x01),
	.name = "ARIANE RISC-V",
	.features = SBI_ARIANE_FEATURES,
	.hart_count = ARIANE_HART_COUNT,
	.hart_stack_size = 4096,
	.disabled_hart_mask = 0,
	.platform_ops_addr = (unsigned long)&platform_ops
};
