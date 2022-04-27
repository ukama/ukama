/* ----------------------------------------------------------------------------
 *         MICROCHIP/ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2018, Microchip Technology, Inc.
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Microchip's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY MICROCHIP "AS IS" AND ANY EXPRESS OR
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
#include "debug.h"
#include "board.h"
#include "flexcom.h"

#ifndef AT91C_NUM_FLEXCOM
#error "Must define AT91C_NUM_FLEXCOM by board definitions"
#endif

struct at91_flexcom *flexcoms;

/* Must be called by board init routine
 * Board must enable the clocks for all flexcoms and configure the PIO
 * accordingly.
 */
void flexcoms_init(struct at91_flexcom *board_specific_flexcoms)
{
	flexcoms = board_specific_flexcoms;
}

unsigned int flexcom_get_base(unsigned int flexcom_num)
{
	if (flexcom_num >= AT91C_NUM_FLEXCOM)
		return 0;

	return flexcoms[flexcom_num].addr;
}

unsigned int flexcom_get_regmap(unsigned int flexcom_num)
{
	if (flexcom_num >= AT91C_NUM_FLEXCOM)
		return 0;

	return flexcoms[flexcom_num].regmap;
}

int flexcom_set_mode(unsigned int flexcom_num, unsigned int mode)
{
	unsigned int base_addr = flexcom_get_base(flexcom_num);

	if (base_addr) {
		writel(mode, base_addr);
		flexcoms[flexcom_num].mode = mode;
		return 0;
	}

	/* failed */
	return -1;
}

int flexcom_init(unsigned int flexcom_num)
{
	int ret;

	ret = flexcom_set_mode(flexcom_num, flexcoms[flexcom_num].mode);
	if (ret)
		return ret;

	switch (flexcoms[flexcom_num].mode) {
	case FLEXCOM_USART:
		flexcoms[flexcom_num].regmap = flexcoms[flexcom_num].addr +
						AT91C_OFFSET_FLEXCOM_USART;
		break;
	case FLEXCOM_TWI:
		flexcoms[flexcom_num].regmap = flexcoms[flexcom_num].addr +
						AT91C_OFFSET_FLEXCOM_TWI;
		break;
	case FLEXCOM_SPI:
		flexcoms[flexcom_num].regmap = flexcoms[flexcom_num].addr +
						AT91C_OFFSET_FLEXCOM_SPI;
		break;
	}

	return 0;
}
