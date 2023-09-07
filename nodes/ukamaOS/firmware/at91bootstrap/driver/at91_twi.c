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
#include "board.h"
#include "arch/at91_twi.h"
#include "div.h"
#include "debug.h"
#include "pmc.h"

#if defined(CONFIG_AT91SAM9X5)
#define TWI_CLK_OFFSET (4)
#elif defined(CONFIG_SAMA5D2)
#define TWI_CLK_OFFSET (3) /* TODO: handle GCK case (offset=0) */
#elif defined(CONFIG_SAMA5D3X)
#define TWI_CLK_OFFSET (4)
#elif defined(CONFIG_SAMA5D4)
#define TWI_CLK_OFFSET (4)
#else
#define TWI_CLK_OFFSET (4)
#endif

#define TWI_CLOCK	400000

unsigned int twi_init_done;

unsigned char hdmi_twi_bus;
unsigned char wm8904_twi_bus;
unsigned char act8865_twi_bus;
unsigned char at24xx_twi_bus;

#define AT91_MAX_TWI_SUPPORTED		16

static unsigned int at91_twi_base[AT91_MAX_TWI_SUPPORTED];
static unsigned int at91_twi_cur_num = 0;

int at91_twi_register_bus(unsigned int twi_base)
{
	if (at91_twi_cur_num >= AT91_MAX_TWI_SUPPORTED)
		return -1;

	at91_twi_base[at91_twi_cur_num] = twi_base;

	at91_twi_cur_num++;

	return at91_twi_cur_num - 1;
}

static inline unsigned int twi_reg_read(unsigned int twi_base,
					unsigned int offset)
{
	return readl(twi_base + offset);
}

static inline void twi_reg_write(unsigned int twi_base,
				unsigned int offset, unsigned int value)
{
	writel(value, twi_base + offset);
}

static unsigned int get_twi_base(unsigned int bus)
{
	if (bus >= at91_twi_cur_num)
		return 0;

	return at91_twi_base[bus];
}

static int twi_configure_master_mode(unsigned int bus,
			unsigned int bus_clock, unsigned int twi_clock)
{
	unsigned int loop = 1;
	unsigned int clkdiv, ckdiv = 0;
	unsigned int clock =  bus_clock;
	unsigned int twi_base;
	unsigned int reg, version;

	twi_base = get_twi_base(bus);
	if (!twi_base)
		return -1;

	version = twi_reg_read(twi_base, TWI_VERSION);

	twi_reg_write(twi_base, TWI_CR, TWI_CR_SWRST);
	twi_reg_read(twi_base, TWI_RHR);

	twi_reg_write(twi_base, TWI_CR, TWI_CR_SVDIS);
	twi_reg_write(twi_base, TWI_CR, TWI_CR_MSEN);

	while (loop) {
		clkdiv = (div(clock, (2 * twi_clock)) - TWI_CLK_OFFSET) >> ckdiv;
		if (clkdiv <= 255)
			loop = 0;
		else
			ckdiv++;
	}

	reg = (ckdiv << 16) | (clkdiv << 8) | clkdiv;

	if (version >= 0x704)
		reg |= TWI_CWGR_HOLD_(25);

	twi_reg_write(twi_base, TWI_CWGR, reg);

	return 0;
}

static void twi_stop(unsigned int twi_base)
{
	twi_reg_write(twi_base, TWI_CR, TWI_CR_STOP);
}

static unsigned char twi_readbyte(unsigned int twi_base)
{
	return twi_reg_read(twi_base, TWI_RHR);
}

static void twi_writebyte(unsigned int twi_base, unsigned char byte)
{
	twi_reg_write(twi_base, TWI_THR, byte);
}

static unsigned char twi_check_rxrdy(unsigned int twi_base)
{
	return ((twi_reg_read(twi_base, TWI_SR) & TWI_SR_RXRDY)
							== TWI_SR_RXRDY);
}

static unsigned char twi_check_txrdy(unsigned int twi_base)
{
	return ((twi_reg_read(twi_base, TWI_SR) & TWI_SR_TXRDY)
							== TWI_SR_TXRDY);
}

static unsigned char twi_check_txcompleted(unsigned int twi_base)
{
	return ((twi_reg_read(twi_base, TWI_SR) & TWI_SR_TXCOMP)
							== TWI_SR_TXCOMP);
}

static void twi_startread(unsigned int twi_base, unsigned char device_addr,
			unsigned int internal_addr, unsigned char iaddr_size)
{
	twi_reg_write(twi_base, TWI_MMR, TWI_MMR_IADRSZ(iaddr_size)
				| TWI_MMR_MREAD_RD
				| TWI_MMR_DADR(device_addr));

	twi_reg_write(twi_base, TWI_IADR, internal_addr);

	twi_reg_write(twi_base, TWI_CR, TWI_CR_START);
}

static void twi_startwrite(unsigned int twi_base, unsigned char device_addr,
			unsigned int internal_addr, unsigned char iaddr_size,
			unsigned char byte)
{
	twi_reg_write(twi_base, TWI_MMR, TWI_MMR_IADRSZ(iaddr_size)
				| TWI_MMR_MREAD_WR
				| TWI_MMR_DADR(device_addr));

	twi_reg_write(twi_base, TWI_IADR, internal_addr);

	twi_reg_write(twi_base, TWI_THR, byte);
};

int twi_read(unsigned int bus, unsigned char device_addr,
		unsigned int internal_addr, unsigned char iaddr_size,
		unsigned char *data, unsigned int bytes)
{
	int timeout;
	unsigned int twi_base;

	twi_base = get_twi_base(bus);
	if (!twi_base) {
		dbg_loud("%s: the base address is NULL\n", __func__);
		return -1;
	}

	twi_startread(twi_base, device_addr, internal_addr, iaddr_size);

	while (bytes > 0) {
		if (bytes == 1)
			twi_stop(twi_base);

		timeout = 10000;
		while (!twi_check_rxrdy(twi_base) && (--timeout))
			;

		if (!timeout) {
			dbg_loud("twi read: timeout to wait RXRDY bit on bus %u\n", bus);
			return -1;
		}

		*data++ = twi_readbyte(twi_base);
		bytes--;
	}

	timeout = 10000;
	while (!twi_check_txcompleted(twi_base) && (--timeout))
		;
	if (!timeout) {
		dbg_loud("twi read: timeout to wait TXCOMP bit\n");
		return -1;
	}

	return 0;
}

int twi_write(unsigned int bus, unsigned char device_addr,
		unsigned int internal_addr, unsigned char iaddr_size,
		unsigned char *data, unsigned int bytes)
{
	int timeout;
	unsigned int twi_base;

	twi_base = get_twi_base(bus);
	if (!twi_base)
		return -1;

	twi_startwrite(twi_base, device_addr,
				internal_addr, iaddr_size, *data++);
	bytes--;

	while (bytes > 0) {
		timeout = 10000;
		while (!twi_check_txrdy(twi_base) && (--timeout))
			;

		if (!timeout) {
			dbg_loud("twi write: timeout to wait TXRDY bit\n");
			return -1;
		}

		twi_writebyte(twi_base, *data++);
		bytes--;
	}

	twi_stop(twi_base);


	timeout = 10000;
	while (!twi_check_txcompleted(twi_base) && (--timeout))
		;
	if (!timeout) {
		dbg_loud("twi write: timeout to wait TXCOMP bit\n");
		return -1;
	}

	return 0;
}

int twi_bus_init(unsigned int (*at91_twi_hw_init)(void))
{
	unsigned int bus_clock = at91_get_ahb_clock();
	unsigned int base = at91_twi_hw_init();
	int bus;

	bus = at91_twi_register_bus(base);
	if (bus < 0)
		return bus;

	twi_configure_master_mode(bus, bus_clock, TWI_CLOCK);

	hdmi_twi_bus	= 0xff;
	wm8904_twi_bus	= 0xff;
	act8865_twi_bus	= 0xff;
	at24xx_twi_bus	= 0xff;

	twi_init_done = 1;

	return bus;
}
