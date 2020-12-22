/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2006, Atmel Corporation
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
#ifndef __SPI_H__
#define __SPI_H__

/* SPI mode flags */
#define	SPI_MODE0	0
#define	SPI_MODE1	1
#define	SPI_MODE2	2	
#define	SPI_MODE3	3

/* Controller-specific definitions: */
#define AT91C_SPI_PCS0_DATAFLASH	0
#define AT91C_SPI_PCS1_DATAFLASH	1
#define AT91C_SPI_PCS2_DATAFLASH	2
#define AT91C_SPI_PCS3_DATAFLASH	3

/* functions */
extern void at91_spi_cs_activate(void);
extern void at91_spi_cs_deactivate(void);

extern void at91_spi_enable(void);
extern void at91_spi_disable(void);
extern int at91_spi_init(unsigned int pcs,
			unsigned int clock,
			unsigned int mode);
extern void at91_spi_write_data(unsigned short data);
extern unsigned int at91_spi_read_spi(void);
extern unsigned int at91_spi_read_sr(void);

#endif	/* #ifndef __SPI_H__ */
