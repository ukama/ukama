/*
 * This file is part of the coreboot project.
 *
 * Copyright 2015 MediaTek Inc.
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

#include <assert.h>
#include <console/console.h>
#include <delay.h>
#include <device/i2c_simple.h>
#include <symbols.h>
#include <timer.h>
#include <device/mmio.h>
#include <soc/addressmap.h>
#include <soc/i2c.h>
#include <device/mmio.h>
#include <soc/pll.h>
#include <soc/i2c.h>

#define I2C_CLK_HZ (AXI_HZ / 16)

struct mtk_i2c mtk_i2c_bus_controller[7] = {
	/* i2c0 setting */
	{
		.i2c_regs = (void *)I2C_BASE,
		.i2c_dma_regs = (void *)(I2C_DMA_BASE + 0x80),
	},

	/* i2c1 setting */
	{
		.i2c_regs = (void *)(I2C_BASE + 0x1000),
		.i2c_dma_regs = (void *)(I2C_DMA_BASE + 0x100),
	},

	/* i2c2 setting */
	{
		.i2c_regs = (void *)(I2C_BASE + 0x2000),
		.i2c_dma_regs = (void *)(I2C_DMA_BASE + 0x180),
	},

	/* i2c3 setting */
	{
		.i2c_regs = (void *)(I2C_BASE + 0x9000),
		.i2c_dma_regs = (void *)(I2C_DMA_BASE + 0x200),
	},

	/* i2c4 setting */
	{
		.i2c_regs = (void *)(I2C_BASE + 0xa000),
		.i2c_dma_regs = (void *)(I2C_DMA_BASE + 0x280),
	},

	/* i2c5 is reserved for internal use. */
	{
	},

	/* i2c6 setting */
	{
		.i2c_regs = (void *)(I2C_BASE + 0xc000),
		.i2c_dma_regs = (void *)I2C_DMA_BASE,
	}
};

#define I2CTAG                "[I2C][PL] "

#if CONFIG(DEBUG_I2C)
#define I2CLOG(fmt, arg...)   printk(BIOS_INFO, I2CTAG fmt, ##arg)
#else
#define I2CLOG(fmt, arg...)
#endif /* CONFIG_DEBUG_I2C */

#define I2CERR(fmt, arg...)   printk(BIOS_ERR, I2CTAG fmt, ##arg)

void mtk_i2c_bus_init(uint8_t bus)
{
	uint8_t step_div;
	uint32_t i2c_freq;
	const uint8_t sample_div = 1;

	assert(bus < ARRAY_SIZE(mtk_i2c_bus_controller));

	/* Calculate i2c frequency */
	step_div = DIV_ROUND_UP(I2C_CLK_HZ, (400 * KHz * sample_div * 2));
	i2c_freq = I2C_CLK_HZ / (step_div * sample_div * 2);
	assert(sample_div < 8 && step_div < 64 && i2c_freq < 400 * KHz &&
	       i2c_freq >= 380 * KHz);

	/* Init i2c bus Timing register */
	write32(&mtk_i2c_bus_controller[bus].i2c_regs->timing,
		(sample_div - 1) << 8 | (step_div - 1));
}
