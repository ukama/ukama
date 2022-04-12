/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2013, Atmel Corporation
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
#include "common.h"
#include "hardware.h"
#include "arch/at91_pmc/pmc.h"
#include "pmc.h"
#include "arch/tz_matrix.h"
#include "debug.h"

#define MATRIX_AXIMX	1
#define MATRIX_H64MX	2
#define MATRIX_H32MX	3

#define SECURITY_TYPE_AS	1
#define SECURITY_TYPE_NS	2
#define	SECURITY_TYPE_PS	3

struct peri_security {
	unsigned int	peri_id;
	unsigned int	matrix;
	unsigned int	security_type;
};

static const struct peri_security peri_security_array[] = {
#ifdef CONFIG_SAMA5D4
	/* SAIC */
	{
		.peri_id = AT91C_ID_FIQ,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* SYS */
	{
		.peri_id = AT91C_ID_SYS,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* ARM */
	{
		.peri_id = AT91C_ID_ARM,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* PIT */
	{
		.peri_id = AT91C_ID_PIT,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* WDT */
	{
		.peri_id = AT91C_ID_WDT,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* PIOD */
	{
		.peri_id = AT91C_ID_PIOD,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* USART0 */
	{
		.peri_id = AT91C_ID_USART0,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* USART1 */
	{
		.peri_id = AT91C_ID_USART1,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* XDMAC0 */
	{
		.peri_id = AT91C_ID_XDMAC0,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* ICM */
	{
		.peri_id = AT91C_ID_ICM,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* PKCC */
	{
		.peri_id = AT91C_ID_PKCC,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* SCI */
	{
		.peri_id = AT91C_ID_SCI,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* AES */
	{
		.peri_id = AT91C_ID_AES,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* AESB */
	{
		.peri_id = AT91C_ID_AESB,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* TDES */
	{
		.peri_id = AT91C_ID_TDES,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* SHA */
	{
		.peri_id = AT91C_ID_SHA,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* MPDDRC */
	{
		.peri_id = AT91C_ID_MPDDRC,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* MATRIX1 */
	{
		.peri_id = AT91C_ID_MATRIX1,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* MATRIX0 */
	{
		.peri_id = AT91C_ID_MATRIX0,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* VDEC */
	{
		.peri_id = AT91C_ID_VDEC,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* SECUMOD */
	{
		.peri_id = AT91C_ID_SECUMOD,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* MSADCC */
	{
		.peri_id = AT91C_ID_MSADCC,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* HSMC */
	{
		.peri_id = AT91C_ID_HSMC,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* PIOA*/
	{
		.peri_id = AT91C_ID_PIOA,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* PIOB*/
	{
		.peri_id = AT91C_ID_PIOB,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* PIOC*/
	{
		.peri_id = AT91C_ID_PIOC,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* PIOE*/
	{
		.peri_id = AT91C_ID_PIOE,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* UART0 */
	{
		.peri_id = AT91C_ID_UART0,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* UART1 */
	{
		.peri_id = AT91C_ID_UART1,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* USART2 */
	{
		.peri_id = AT91C_ID_USART2,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* USART3 */
	{
		.peri_id = AT91C_ID_USART3,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* USART4 */
	{
		.peri_id = AT91C_ID_USART4,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* TWI0 */
	{
		.peri_id = AT91C_ID_TWI0,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* TWI1 */
	{
		.peri_id = AT91C_ID_TWI1,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* TWI2 */
	{
		.peri_id = AT91C_ID_TWI2,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* HSMCI0 */
	{
		.peri_id = AT91C_ID_HSMCI0,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* HSMCI1 */
	{
		.peri_id = AT91C_ID_HSMCI1,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* SPI0 */
	{
		.peri_id = AT91C_ID_SPI0,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* SPI1 */
	{
		.peri_id = AT91C_ID_SPI1,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* SPI2 */
	{
		.peri_id = AT91C_ID_SPI2,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* TC0 */
	{
		.peri_id = AT91C_ID_TC0,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* TC1 */
	{
		.peri_id = AT91C_ID_TC1,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* TC2 */
	{
		.peri_id = AT91C_ID_TC2,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* PWM */
	{
		.peri_id = AT91C_ID_PWM,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* ADC */
	{
		.peri_id = AT91C_ID_ADC,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* DBGU */
	{
		.peri_id = AT91C_ID_DBGU,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* UHPHS */
	{
		.peri_id = AT91C_ID_UHPHS,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* UDPHS */
	{
		.peri_id = AT91C_ID_UDPHS,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* SSC0 */
	{
		.peri_id = AT91C_ID_SSC0,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* SSC1 */
	{
		.peri_id = AT91C_ID_SSC1,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* XDMAC1 */
	{
		.peri_id = AT91C_ID_XDMAC1,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_NS,
	},
	/* LCDC */
	{
		.peri_id = AT91C_ID_LCDC,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* ISI */
	{
		.peri_id = AT91C_ID_ISI,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* TRNG */
	{
		.peri_id = AT91C_ID_TRNG,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* GMAC */
	{
		.peri_id = AT91C_ID_GMAC,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* GMAC1 */
	{
		.peri_id = AT91C_ID_GMAC1,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* IRQ */
	{
		.peri_id = AT91C_ID_IRQ,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_NS,
	},
	/* SFC */
	{
		.peri_id = AT91C_ID_SFC,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* SECURAM */
	{
		.peri_id = AT91C_ID_SECURAM,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* CTB */
	{
		.peri_id = AT91C_ID_CTB,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* SMD */
	{
		.peri_id = AT91C_ID_SMD,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* TWI3 */
	{
		.peri_id = AT91C_ID_TWI3,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_PS,
	},
	/* CATB */
	{
		.peri_id = AT91C_ID_CATB,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* SFR */
	{
		.peri_id = AT91C_ID_SFR,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* AIC */
	{
		.peri_id = AT91C_ID_AIC,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_NS,
	},
	/* SAIC */
	{
		.peri_id = AT91C_ID_SAIC,
		.matrix = MATRIX_H32MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* L2CC */
	{
		.peri_id = AT91C_ID_L2CC,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_AS,
	},
	/* PMC */
	{
		.peri_id = AT91C_ID_PMC,
		.matrix = MATRIX_H64MX,
		.security_type = SECURITY_TYPE_AS,
	},
#endif /* CONFIG_SAMA5D4 */
};


static inline void matrix_write(unsigned int base,
				unsigned int offset,
				const unsigned int value)
{
	writel(value, offset + base);
}

static inline unsigned int matrix_read(int base, unsigned int offset)
{
	return readl(offset + base);
}

void matrix_write_protect_enable(unsigned int matrix_base)
{
	matrix_write(matrix_base, MATRIX_WPMR,
		(MATRIX_WPMR_WPKEY_PASSWD | MATRIX_WPMR_WPEN_ENABLE));
}

void matrix_write_protect_disable(unsigned int matrix_base)
{
	matrix_write(matrix_base, MATRIX_WPMR, MATRIX_WPMR_WPKEY_PASSWD);
}

void matrix_configure_slave_security(unsigned int matrix_base,
				unsigned int slave,
				unsigned int srtop_setting,
				unsigned int srsplit_setting,
				unsigned int ssr_setting)
{
	matrix_write(matrix_base, MATRIX_SSR(slave), ssr_setting);
	matrix_write(matrix_base, MATRIX_SRTSR(slave), srtop_setting);
	matrix_write(matrix_base, MATRIX_SASSR(slave), srsplit_setting);
}

void matrix_read_slave_security(void)
{
	unsigned int matrix_base;
	unsigned int slave;

	dbg_very_loud("\n\nMATRIX64:\n");
	matrix_base = AT91C_BASE_MATRIX64;
	for (slave = 0; slave < 13; slave++) {
		dbg_very_loud("MATRIX_SRTSR%d: %x, MATRIX_SASSR%d: %x, MATRIX_SSR%d: %x\n",
			slave, matrix_read(matrix_base, MATRIX_SRTSR(slave)),
			slave, matrix_read(matrix_base, MATRIX_SASSR(slave)),
			slave, matrix_read(matrix_base, MATRIX_SSR(slave)));
	}

	dbg_very_loud("\n\nMATRIX32:\n");
	matrix_base = AT91C_BASE_MATRIX32;
	for (slave = 0; slave < 7; slave++) {
		dbg_very_loud("MATRIX_SRTSR%d: %x, MATRIX_SASSR%d: %x, MATRIX_SSR%d: %x\n",
			slave, matrix_read(matrix_base, MATRIX_SRTSR(slave)),
			slave, matrix_read(matrix_base, MATRIX_SASSR(slave)),
			slave, matrix_read(matrix_base, MATRIX_SSR(slave)));
	}
}

void matrix_read_periperal_security(void)
{
	unsigned int i;

	unsigned int matrix_base = AT91C_BASE_MATRIX32;
	dbg_very_loud("\n\nMATRIX32\n");
	for (i = 0; i < 3; i++) {
		dbg_very_loud("MATRIX_SPSELR(%d): %x\n",
				i, matrix_read(matrix_base, MATRIX_SPSELR(i)));
	}

	matrix_base = AT91C_BASE_MATRIX64;
	dbg_very_loud("\n\n_MATRIX64\n");
	for (i = 0; i < 3; i++) {
		dbg_very_loud("MATRIX_SPSELR(%d): %x\n",
				i, matrix_read(matrix_base, MATRIX_SPSELR(i)));
	}
}

static struct peri_security *get_peri_security(unsigned int peri_id)
{
	unsigned int i;

	for (i = 0; i < ARRAY_SIZE(peri_security_array); i++) {
		if (peri_id == peri_security_array[i].peri_id)
			break;
	}

	return ((i == ARRAY_SIZE(peri_security_array))) ?
			NULL : (struct peri_security *)&peri_security_array[i];
}

static int matrix_set_peri_security(unsigned int matrix, unsigned peri_id)
{
	unsigned int base;
	unsigned int spselr;
	unsigned int idx;
	unsigned int bit;

	idx = peri_id / 32;
	if (idx > 3)
		return -1;

	bit = (0x01 << (peri_id % 32));
	if (matrix == MATRIX_H32MX)
		base = AT91C_BASE_MATRIX32;
	else if (matrix == MATRIX_H64MX)
		base = AT91C_BASE_MATRIX64;
	else
		return -1;

	spselr = matrix_read(base, MATRIX_SPSELR(idx));
	spselr |= bit;
	matrix_write(base, MATRIX_SPSELR(idx), spselr);

	return 0;
}

int matrix_configure_peri_security(unsigned int *peri_id_array,
					unsigned int size)
{
	unsigned int i;
	unsigned int *peri_id_p;
	unsigned int matrix;
	unsigned int peri_id;
	struct peri_security *periperal_sec;
	int ret;

	if ((peri_id_array == NULL) || (size == 0))
		return -1;

	peri_id_p = peri_id_array;
	for (i = 0; i < size; i++) {
		periperal_sec = get_peri_security(*peri_id_p);
		if (periperal_sec == NULL)
			return -1;

		if (periperal_sec->security_type != SECURITY_TYPE_PS)
			return -1;

		matrix = periperal_sec->matrix;
		peri_id = *peri_id_p;
		ret = matrix_set_peri_security(matrix, peri_id);
		if (ret)
			return -1;

		peri_id_p++;
	}

	return 0;
}

#ifdef CONFIG_SAMA5D4
/*
 * is_peripheral_secure - tell if the peripheral is in secure mode
 * @periph_id: the peripheral id that is checked
 *
 * Check security of a particular peripheral by providing its ID.
 * Note that a wrong preripheral ID leads to the "true" return code.
 */
int is_peripheral_secure(unsigned int periph_id)
{
	struct peri_security *periperal_sec;
	unsigned int matrix;
	unsigned int base;
	unsigned int mask;

	if ((periph_id > AT91C_ID_FIQ) && (periph_id < AT91C_ID_COUNTS)) {
		/* special cases here */
		if ((periph_id == AT91C_ID_IRQ)
		 || (periph_id == AT91C_ID_AIC)
		 || (periph_id == AT91C_ID_XDMAC1))
			return 0;

		periperal_sec = get_peri_security(periph_id);
		if (periperal_sec == NULL)
			return -1;

		matrix = periperal_sec->matrix;
		if (matrix == MATRIX_H32MX)
			base = AT91C_BASE_MATRIX32;
		else if (matrix == MATRIX_H64MX)
			base = AT91C_BASE_MATRIX64;
		else
			return -1;

		mask = 1 << (periph_id % 32);
		if (matrix_read(base, MATRIX_SPSELR(periph_id / 32)) & mask)
			return 0;

	}
	return 1;
}

int is_sys_clk_secure(unsigned int sys_mask)
{
	unsigned int periph_id = sys_mask_to_per_id(sys_mask);

	return is_peripheral_secure(periph_id);
}

int is_pck_clk_secure(unsigned int pck_mask)
{
	switch (pck_mask) {
	case AT91C_PMC_PCK0:
	case AT91C_PMC_PCK1:
	case AT91C_PMC_PCK2:
		return 0;
	default:
		return 1;
	}
}

int is_usb_hs_secure(void)
{
	return is_peripheral_secure(AT91C_ID_UHPHS)
		|| is_peripheral_secure(AT91C_ID_UDPHS);
}

int is_usb_host_secure(void)
{
	return is_peripheral_secure(AT91C_ID_UHPHS);
}

int is_switching_clock_forbiden(unsigned int periph_id, unsigned int is_on, unsigned int *silent)
{
	/* disable console clock : forbiden */
	if ((periph_id == AT91C_ID_USART3) && (is_on == 0)) {
		/* keep it silent */
		*silent = 1;
		return 1;
	} else {
		return 0;
	}
}
#endif /* CONFIG_SAMA5D4 */
