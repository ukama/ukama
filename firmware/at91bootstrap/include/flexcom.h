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

#ifndef FLEXCOM_H_
#define FLEXCOM_H_

#ifdef CONFIG_FLEXCOM

#define AT91C_OFFSET_FLEXCOM_USART		0x200
#define AT91C_OFFSET_FLEXCOM_SPI		0x400
#define AT91C_OFFSET_FLEXCOM_TWI		0x600

#define FLEXCOM_USART		0x1
#define FLEXCOM_SPI		0x2
#define FLEXCOM_TWI		0x3

struct at91_flexcom {
	unsigned int id; 		/* SOC Peripheral ID */
	unsigned int mode;		/* Peripheral mode : SPI/USART/TWI */
	unsigned int addr; 		/* Flexcom Address */
	unsigned int regmap;		/* register base for either SPI/USART/TWI */
};

extern void flexcoms_init(struct at91_flexcom *board_specific_flexcoms);
extern int flexcom_set_mode(unsigned int flexcom_num, unsigned int mode);
extern int flexcom_init(unsigned int flexcom_num);
extern unsigned int flexcom_get_base(unsigned int flexcom_num);
extern unsigned int flexcom_get_regmap(unsigned int flexcom_num);

#endif /* CONFIG_FLEXCOM */
#endif /* FLEXCOM_H_ */
