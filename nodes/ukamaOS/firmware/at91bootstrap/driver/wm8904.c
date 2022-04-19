/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2014, Atmel Corporation
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Atmel's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY ATMEL "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT ARE
 * DISCLAIMED. IN NO EVENT SHALL ATMEL BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
#include "hardware.h"
#include "twi.h"
#include "debug.h"

/*
 * WM8904 Register MAP.
 */
#define WM8904_SW_RESET_AND_ID		(0x00)
#define WM8904_REVISION			(0x01)
#define WM8904_BIAS_CONTROL_0		(0x04)
#define WM8904_VMID_CONTROL_0		(0x05)

/*
 * Field Definitions.
 */

/*
 * R4 (0x04) - Bias Control 0
 */
#define WM8904_BIAS_ENA			(0x0001)	/* BIAS_ENA */

/*
 * R5 (0x05) - VMID Control 0
 */
#define WM8904_VMID_ENA			(0x0001)	/* VMID_ENA */
#define WM8904_VMID_RES			(0x0006)	/* VMID_RES - [2:1] */
#define WM8904_VMID_BUF_ENA		(0x0040)	/* VMID_BUF_ENA */

/*
 * WM8904 Slave Device Address
 */
#define WM8904_ADDR		(0x1a)

static unsigned int wm8904_get_twi_bus(void)
{
	unsigned int bus = 0;

	if (wm8904_twi_bus != 0xff) {
		bus = wm8904_twi_bus;
	} else {
#if defined(CONFIG_CODEC_ON_TWI0)
		bus = 0;
#elif defined(CONFIG_CODEC_ON_TWI1)
		bus = 1;
#elif defined(CONFIG_CODEC_ON_TWI2)
		bus = 2;
#elif defined(CONFIG_CODEC_ON_TWI3)
		bus = 3;
#endif
	}

	return bus;
}

static int wm8904_read(unsigned char reg_addr, unsigned short *data)
{
	unsigned int bus;
	int ret;

	bus = wm8904_get_twi_bus();

	ret = twi_read(bus, WM8904_ADDR, reg_addr, 1, (unsigned char *)data, 2);
	if (ret) {
		dbg_loud("WM8904: Failed to read on TWI bus: %d\n", bus);
		return -1;
	}

	return 0;
}

static int wm8904_write(unsigned char reg_addr, unsigned short data)
{
	unsigned int bus;
	int ret;

	bus = wm8904_get_twi_bus();

	ret = twi_write(bus, WM8904_ADDR,
				reg_addr, 1, (unsigned char *)&data, 2);
	if (ret) {
		dbg_loud("WM8904: Failed to write on TWI bus: %d\n", bus);
		return -1;
	}

	return 0;
}

static int wm8904_chipid(void)
{
	unsigned short reg = 0;
	unsigned short chip_id;
	int ret;

	ret = wm8904_read(WM8904_SW_RESET_AND_ID, &reg);
	if (ret)
		return -1;

	chip_id = ((reg << 8) & 0xff00) | ((reg >> 8) & 0x00ff);

	return (chip_id == 0x8904) ? 0 : -1;
}

static int wm8904_version(void)
{
	unsigned short reg = 0;
	int ret;

	ret = wm8904_read(WM8904_REVISION, &reg);
	if (ret)
		return -1;

	return reg + 'A';
}

static int wm8904_reset(void)
{
	int ret;

	ret = wm8904_write(WM8904_SW_RESET_AND_ID, 0);
	if (ret)
		return -1;

	return 0;
}

static int wm8904_disable_vmid_ref(void)
{
	unsigned short reg = 0;
	int ret;

	ret = wm8904_read(WM8904_VMID_CONTROL_0, &reg);
	if (ret)
		return -1;

	reg &= ~WM8904_VMID_ENA;
	reg &= ~WM8904_VMID_RES;
	ret = wm8904_write(WM8904_VMID_CONTROL_0, reg);
	if (ret)
		return -1;

	return 0;
}

static int wm8904_disable_bias(void)
{
	unsigned short reg = 0;
	int ret;

	ret = wm8904_read(WM8904_BIAS_CONTROL_0, &reg);
	if (ret)
		return -1;

	reg &= ~WM8904_BIAS_ENA;
	ret = wm8904_write(WM8904_BIAS_CONTROL_0, reg);
	if (ret)
		return -1;

	return 0;
}

int wm8904_enter_low_power(void)
{
	int ret;

	ret = wm8904_chipid();
	if (ret == -1) {
		dbg_loud("WM8904: Failed to read Chip ID\n");
		return -1;
	}

	ret = wm8904_version();
	if (ret == -1) {
		dbg_loud("WM8904: Failed to read version\n");
		return -1;
	}

	ret = wm8904_reset();
	if (ret) {
		dbg_loud("WM8904: Failed to issue software reset\n");
		return -1;
	}

	ret = wm8904_disable_vmid_ref();
	if (ret) {
		dbg_loud("WM8904: Failed to disable VMID Control 0\n");
		return -1;
	}

	ret = wm8904_disable_bias();
	if (ret) {
		dbg_loud("WM8904: Failed to disable Bias Control 0\n");
		return -1;
	}

	dbg_loud("WM8904: Enter Low-power mode\n");

	return 0;
}
