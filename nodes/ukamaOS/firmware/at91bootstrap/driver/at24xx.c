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
#include "twi.h"
#include "board_hw_info.h"
#include "debug.h"

#ifdef CONFIG_SAMA5D2_ICP
/* ICP EK has different EEPROM */
#define	EK_AT24XX_ADDR		0x50
#define	MAX_AT24XX_BYTES	128
#elif defined(CONFIG_SAM9X60EK) || defined(CONFIG_SAM9X60_DDR2_SIP_EB)
#define	EK_AT24XX_ADDR		0x53
#define	MAX_AT24XX_BYTES	128
#else
#define EK_AT24XX_ADDR		0x54
#define	MAX_AT24XX_BYTES	256
#endif

#define EK_INFO_OFFSET		(MAX_AT24XX_BYTES - HW_INFO_TOTAL_SIZE)

static unsigned int at24_get_twi_bus(void)
{
	unsigned int bus = 0;

	if (at24xx_twi_bus != 0xff) {
		bus = at24xx_twi_bus;
	} else {
#if defined(CONFIG_EEPROM_ON_TWI0)
		bus = 0;
#elif defined(CONFIG_EEPROM_ON_TWI1)
		bus = 1;
#elif defined(CONFIG_EEPROM_ON_TWI2)
		bus = 2;
#elif defined(CONFIG_EEPROM_ON_TWI3)
		bus = 3;
#endif
	}

	return bus;
}
static int at24_read(unsigned char device_addr, unsigned char internal_addr,
			unsigned char *buff, unsigned int length)
{
	unsigned char iaddr_size = 1;
	unsigned int bus;
	int ret = 0;

	bus = at24_get_twi_bus();

	ret = twi_read(bus, device_addr, internal_addr, iaddr_size, buff, length);
	if (ret)
		dbg_loud("EEPROM: Failed to read\n");

	return ret;
}

int load_ek_at24xx(unsigned char *buff, unsigned int length)
{
	unsigned char dev_addr = EK_AT24XX_ADDR;
	unsigned char offset = EK_INFO_OFFSET;
	int ret = 0;

	if (!twi_init_done)
		twi_init();

	ret = at24_read(dev_addr, offset, buff, length);

	if (!ret) {
		dbg_loud("EEPROM Buff:\n");
#if (BOOTSTRAP_DEBUG_LEVEL >= DEBUG_LOUD)
		dbg_hexdump(buff, length, DUMP_WIDTH_BIT_8);
#endif
	}

	return ret;
}
