/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2010, Atmel Corporation
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
#include "arch/at91_dbgu.h"

#ifdef CONFIG_USART

#ifndef USART_BASE
#define USART_BASE	AT91C_BASE_DBGU
#endif

static inline void write_usart(unsigned int offset, const unsigned int value)
{
	writel(value, offset + USART_BASE);
}

static inline unsigned int read_usart(unsigned int offset)
{
	return readl(offset + USART_BASE);
}

void usart_init(unsigned int baudrate)
{
	/* Disable interrupts */
	write_usart(DBGU_IDR, -1);

	/* Reset the receiver and transmitter */
	write_usart(DBGU_CR, AT91C_DBGU_RSTRX
				| AT91C_DBGU_RSTTX
				| AT91C_DBGU_RXDIS
				| AT91C_DBGU_TXDIS);

	/* Configure the baudrate */
	write_usart(DBGU_BRGR, baudrate);

	/* Configure USART in Asynchronous mode */
	write_usart(DBGU_MR, AT91C_DBGU_PAR_NONE
				| AT91C_DBGU_CHMODE_NORMAL
				| AT91C_DBGU_CHRL_8BIT
				| AT91C_DBGU_NBSTOP_1BIT);

	/* Enable RX and Tx */
	write_usart(DBGU_CR, AT91C_DBGU_RXEN | AT91C_DBGU_TXEN);
}

static void usart_putc(const char c)
{
	while (!(read_usart(DBGU_CSR) & AT91C_DBGU_TXRDY))
		;

	write_usart(DBGU_THR, c);
}

void usart_puts(const char *ptr)
{
	int i = 0;

	while (ptr[i] != '\0') {
		if (ptr[i] == '\n')
			usart_putc('\r');
		usart_putc(ptr[i]);
		i++;
	}
}

char usart_getc(void)
{
	while (!(read_usart(DBGU_CSR) & AT91C_DBGU_RXRDY))
		;

	return (char)read_usart(DBGU_RHR);
}

#else

void usart_puts(const char *ptr)
{

}

#endif
