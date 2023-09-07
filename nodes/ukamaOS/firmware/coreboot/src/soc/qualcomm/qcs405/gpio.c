/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018-2019 Qualcomm Inc.
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

#include <device/mmio.h>
#include <types.h>
#include <delay.h>
#include <gpio.h>

void gpio_configure(gpio_t gpio, uint32_t func, uint32_t pull,
				uint32_t drive_str, uint32_t enable)
{
	uint32_t reg_val;
	struct tlmm_gpio *regs = (void *)(uintptr_t)gpio.addr;

	reg_val = ((enable & GPIO_CFG_OE_BMSK) << GPIO_CFG_OE_SHFT) |
		  ((drive_str & GPIO_CFG_DRV_BMSK) << GPIO_CFG_DRV_SHFT) |
		  ((func & GPIO_CFG_FUNC_BMSK) << GPIO_CFG_FUNC_SHFT) |
		  ((pull & GPIO_CFG_PULL_BMSK) << GPIO_CFG_PULL_SHFT);

	write32(&regs->cfg, reg_val);
}

void gpio_set(gpio_t gpio, int value)
{
	struct tlmm_gpio *regs = (void *)(uintptr_t)gpio.addr;
	write32(&regs->in_out, (!!value) << GPIO_IO_OUT_SHFT);
}

int gpio_get(gpio_t gpio)
{
	struct tlmm_gpio *regs = (void *)(uintptr_t)gpio.addr;

	return ((read32(&regs->in_out) >> GPIO_IO_IN_SHFT) &
		GPIO_IO_IN_BMSK);
}

void gpio_input_pulldown(gpio_t gpio)
{
	gpio_configure(gpio, GPIO_FUNC_GPIO,
				GPIO_PULL_DOWN, GPIO_2MA, GPIO_DISABLE);
}

void gpio_input_pullup(gpio_t gpio)
{
	gpio_configure(gpio, GPIO_FUNC_GPIO,
				GPIO_PULL_UP, GPIO_2MA, GPIO_DISABLE);
}

void gpio_input(gpio_t gpio)
{
	gpio_configure(gpio, GPIO_FUNC_GPIO,
				GPIO_NO_PULL, GPIO_2MA, GPIO_DISABLE);
}

void gpio_output(gpio_t gpio, int value)
{
	gpio_set(gpio, value);
	gpio_configure(gpio, GPIO_FUNC_GPIO,
				GPIO_NO_PULL, GPIO_2MA, GPIO_ENABLE);
}

void gpio_input_irq(gpio_t gpio, enum gpio_irq_type type, uint32_t pull)
{
	struct tlmm_gpio *regs = (void *)(uintptr_t)gpio.addr;

	gpio_configure(gpio, GPIO_FUNC_GPIO,
				pull, GPIO_2MA, GPIO_DISABLE);

	clrsetbits_le32(&regs->intr_cfg, GPIO_INTR_DECT_CTL_MASK <<
		GPIO_INTR_DECT_CTL_SHIFT, type << GPIO_INTR_DECT_CTL_SHIFT);
	clrsetbits_le32(&regs->intr_cfg, GPIO_INTR_RAW_STATUS_ENABLE
		<< GPIO_INTR_RAW_STATUS_EN_SHIFT, GPIO_INTR_RAW_STATUS_ENABLE
		<< GPIO_INTR_RAW_STATUS_EN_SHIFT);
}

int gpio_irq_status(gpio_t gpio)
{
	struct tlmm_gpio *regs = (void *)(uintptr_t)gpio.addr;

	if (!(read32(&regs->intr_status) & GPIO_INTR_STATUS_MASK))
		return 0;

	write32(&regs->intr_status, GPIO_INTR_STATUS_DISABLE);
	return 1;
}
