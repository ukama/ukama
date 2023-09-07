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
#include "hardware.h"
#include "arch/at91_pio.h"
#include "gpio.h"
#include "debug.h"

static inline int pio_base_addr(unsigned int pio)
{
	if (pio > AT91C_NUM_PIO)
		return -1;

	if (pio == 0)
		return AT91C_BASE_PIOA;
	else if (pio == 1)
		return AT91C_BASE_PIOB;
	else if (pio == 2)
		return AT91C_BASE_PIOC;
	else if (pio == 3)
#ifdef AT91C_BASE_PIOD
		return AT91C_BASE_PIOD;
#else
		return -1;
#endif
	else if (pio == 4)
#ifdef AT91C_BASE_PIOE
		return AT91C_BASE_PIOE;
#else
		return -1;
#endif
	else
		return -1;
}

static inline void write_pio(unsigned int pio,
			unsigned int offset,
			const unsigned int value)
{
	int base = pio_base_addr(pio);

	if (base == -1)
		return;

	writel(value, offset + base);
}

static inline unsigned int read_pio(unsigned int pio, unsigned int offset)
{
	int base = pio_base_addr(pio);

	if (base == -1)
		return 0;

	return readl(offset + base);
}

static inline unsigned pin_to_controller(unsigned pin)
{
	return (pin) / PIO_NUM_IO;
}

static inline unsigned pin_to_mask(unsigned pin)
{
	return 1 << ((pin) % PIO_NUM_IO);
}

#ifdef CONFIG_CPU_HAS_PIO4
static void pio4_set_periph(unsigned pio, unsigned mask,
				int config, unsigned func)
{
	unsigned int value = func;

	write_pio(pio, PIO_MSKR, mask);

	value |= (config & PIO_PULLUP) ? AT91C_PIO_CFGR_PUEN : 0;
	value |= (config & PIO_PULLDOWN) ? AT91C_PIO_CFGR_PDEN : 0;
	value |= (config & PIO_DRVSTR_LO) ? AT91C_PIO_CFGR_DRVSTR_LOW : 0;
	value |= (config & PIO_DRVSTR_ME) ? AT91C_PIO_CFGR_DRVSTR_MEDIUM : 0;
	value |= (config & PIO_DRVSTR_HI) ? AT91C_PIO_CFGR_DRVSTR_HIGH : 0;

	write_pio(pio, PIO_CFGR, value);
}
#endif

#ifdef CONFIG_CPU_HAS_PIO3
static void pio3_set_drvstr(unsigned pin, int config)
{
	unsigned reg_value;
	unsigned mask = pin_to_mask(pin);
	unsigned pio = pin_to_controller(pin);

	if (config & PIO_DRVSTR_HI) {
		reg_value = read_pio(pio, PIO_DRIVER1);
		write_pio(pio, PIO_DRIVER1, reg_value | mask);
	}
	else if (config & PIO_DRVSTR_LO) {
		reg_value = read_pio(pio, PIO_DRIVER1);
		write_pio(pio, PIO_DRIVER1, reg_value & ~mask);
	}
}

static void pio3_set_slewrate(unsigned pin, int config)
{
	unsigned reg_value;
	unsigned mask = pin_to_mask(pin);
	unsigned pio = pin_to_controller(pin);

	if (config & PIO_SLEWR_CTRL) {
		reg_value = read_pio(pio, PIO_SLEWR);
		write_pio(pio, PIO_SLEWR, reg_value | mask);
	}
}
#endif

static int pio_set_a_periph(unsigned pin, int config)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	if (pio >= AT91C_NUM_PIO)
		return -1;

	if (config & PIO_PULLUP && config & PIO_PULLDOWN)
		return -1;

#ifdef CONFIG_CPU_HAS_PIO4
	pio4_set_periph(pio, mask, config, AT91C_PIO_CFGR_FUNC_PERIPH_A);
#else
	write_pio(pio, PIO_IDR, mask);
	write_pio(pio, ((config & PIO_PULLUP) ? PIO_PPUER : PIO_PPUDR), mask);
#ifdef CONFIG_CPU_HAS_PIO3
	write_pio(pio, ((config & PIO_PULLDOWN) ? PIO_PPDER : PIO_PPDDR), mask);

	pio3_set_drvstr(pin, config);
	pio3_set_slewrate(pin, config);

	write_pio(pio, PIO_SP1, read_pio(pio, PIO_SP1) & ~mask);
	write_pio(pio, PIO_SP2, read_pio(pio, PIO_SP2) & ~mask);
#else
	write_pio(pio, PIO_ASR, mask);
#endif
	write_pio(pio, PIO_PDR, mask);
#endif

	return 0;
}

static int pio_set_b_periph(unsigned pin, int config)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	if (pio >= AT91C_NUM_PIO)
		return -1;

	if (config & PIO_PULLUP && config & PIO_PULLDOWN)
		return -1;

#ifdef CONFIG_CPU_HAS_PIO4
	pio4_set_periph(pio, mask, config, AT91C_PIO_CFGR_FUNC_PERIPH_B);
#else
	write_pio(pio, PIO_IDR, mask);
	write_pio(pio, ((config & PIO_PULLUP) ? PIO_PPUER : PIO_PPUDR), mask);
#ifdef CONFIG_CPU_HAS_PIO3
	write_pio(pio, ((config & PIO_PULLDOWN) ? PIO_PPDER : PIO_PPDDR), mask);

	pio3_set_drvstr(pin, config);
	pio3_set_slewrate(pin, config);

	write_pio(pio, PIO_SP1, read_pio(pio, PIO_SP1) | mask);
	write_pio(pio, PIO_SP2, read_pio(pio, PIO_SP2) & ~mask);
#else
	write_pio(pio, PIO_BSR, mask);
#endif
	write_pio(pio, PIO_PDR, mask);
#endif

	return 0;
}

static int pio_set_c_periph(unsigned pin, int config)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	mask = mask;

	if (pio >= AT91C_NUM_PIO)
		return -1;

	if (config & PIO_PULLUP && config & PIO_PULLDOWN)
		return -1;

#ifdef CONFIG_CPU_HAS_PIO4
	pio4_set_periph(pio, mask, config, AT91C_PIO_CFGR_FUNC_PERIPH_C);
#elif defined CONFIG_CPU_HAS_PIO3
	write_pio(pio, PIO_IDR, mask);
	write_pio(pio, ((config && PIO_PULLUP) ? PIO_PPUER : PIO_PPUDR), mask);
	write_pio(pio, ((config & PIO_PULLDOWN) ? PIO_PPDER : PIO_PPDDR), mask);

	pio3_set_drvstr(pin, config);
	pio3_set_slewrate(pin, config);

	write_pio(pio, PIO_SP1, read_pio(pio, PIO_SP1) & ~mask);
	write_pio(pio, PIO_SP2, read_pio(pio, PIO_SP2) | mask);
	write_pio(pio, PIO_PDR, mask);
#endif

	return 0;
}

static int pio_set_d_periph(unsigned pin, int config)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	mask = mask;

	if (pio >= AT91C_NUM_PIO)
		return -1;

	if (config & PIO_PULLUP && config & PIO_PULLDOWN)
		return -1;

#ifdef CONFIG_CPU_HAS_PIO4
	pio4_set_periph(pio, mask, config, AT91C_PIO_CFGR_FUNC_PERIPH_D);
#elif defined CONFIG_CPU_HAS_PIO3
	write_pio(pio, PIO_IDR, mask);
	write_pio(pio, ((config && PIO_PULLUP) ? PIO_PPUER : PIO_PPUDR), mask);
	write_pio(pio, ((config & PIO_PULLDOWN) ? PIO_PPDER : PIO_PPDDR), mask);

	pio3_set_drvstr(pin, config);
	pio3_set_slewrate(pin, config);

	write_pio(pio, PIO_SP1, read_pio(pio, PIO_SP1) | mask);
	write_pio(pio, PIO_SP2, read_pio(pio, PIO_SP2) | mask);
	write_pio(pio, PIO_PDR, mask);
#endif

	return 0;
}

static int pio_set_e_periph(unsigned pin, int config)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	mask = mask;

	if (pio >= AT91C_NUM_PIO)
		return -1;

	if (config & PIO_PULLUP && config & PIO_PULLDOWN)
		return -1;

#ifdef CONFIG_CPU_HAS_PIO4
	pio4_set_periph(pio, mask, config, AT91C_PIO_CFGR_FUNC_PERIPH_E);
#endif

	return 0;
};

static int pio_set_f_periph(unsigned pin, int config)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	mask = mask;

	if (pio >= AT91C_NUM_PIO)
		return -1;

	if (config & PIO_PULLUP && config & PIO_PULLDOWN)
		return -1;

#ifdef CONFIG_CPU_HAS_PIO4
	pio4_set_periph(pio, mask, config, AT91C_PIO_CFGR_FUNC_PERIPH_F);
#endif

	return 0;
};

static int pio_set_g_periph(unsigned pin, int config)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	mask = mask;

	if (pio >= AT91C_NUM_PIO)
		return -1;

	if (config & PIO_PULLUP && config & PIO_PULLDOWN)
		return -1;

#ifdef CONFIG_CPU_HAS_PIO4
	pio4_set_periph(pio, mask, config, AT91C_PIO_CFGR_FUNC_PERIPH_G);
#endif

	return 0;
};

int pio_set_gpio_input(unsigned pin, int config)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	if (pio >= AT91C_NUM_PIO)
		return -1;

	if (config & PIO_PULLUP && config & PIO_PULLDOWN)
		return -1;

#ifdef CONFIG_CPU_HAS_PIO4
	write_pio(pio, PIO_MSKR, mask);

	mask = AT91C_PIO_CFGR_FUNC_GPIO;
	mask |= (config & PIO_DEGLITCH) ? AT91C_PIO_CFGR_IFEN : 0;
	mask |= (config & PIO_PULLUP) ? AT91C_PIO_CFGR_PUEN : 0;
	mask |= (config & PIO_PULLDOWN) ? AT91C_PIO_CFGR_PDEN : 0;

	write_pio(pio, PIO_CFGR, mask);
#else
	write_pio(pio, ((config & PIO_DEGLITCH) ? PIO_IFER : PIO_IFDR), mask);

	write_pio(pio, PIO_IDR, mask);
	write_pio(pio, ((config & PIO_PULLUP) ? PIO_PPUER : PIO_PPUDR), mask);
#ifdef CONFIG_CPU_HAS_PIO3
	write_pio(pio, ((config & PIO_PULLDOWN) ? PIO_PPDER : PIO_PPDDR), mask);
#endif
	write_pio(pio, PIO_ODR, mask);
	write_pio(pio, PIO_PER, mask);
#endif

	return 0;
}

int pio_set_gpio_output(unsigned pin, int value)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	if (pio >= AT91C_NUM_PIO)
		return -1;

#ifdef CONFIG_CPU_HAS_PIO4
	write_pio(pio, PIO_MSKR, mask);

	write_pio(pio, PIO_CFGR,
			(AT91C_PIO_CFGR_FUNC_GPIO | AT91C_PIO_CFGR_DIR));
	write_pio(pio, (value ? PIO_SODR : PIO_CODR), mask);
#else
	write_pio(pio, PIO_IDR, mask);
	write_pio(pio, PIO_PPUDR, mask);
	write_pio(pio, (value ? PIO_SODR : PIO_CODR), mask);
	write_pio(pio, PIO_OER, mask);
	write_pio(pio, PIO_PER, mask);
#endif

	return 0;
}

static int pio_config_gpio_output(unsigned int pin,
				unsigned int config,
				int value)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	if (pio >= AT91C_NUM_PIO)
		return -1;

#ifdef CONFIG_CPU_HAS_PIO4
	unsigned reg_value;

	write_pio(pio, PIO_MSKR, mask);

	reg_value = AT91C_PIO_CFGR_FUNC_GPIO | AT91C_PIO_CFGR_DIR;
	reg_value |= (config & PIO_OPENDRAIN) ? AT91C_PIO_CFGR_OPD : 0;
	write_pio(pio, PIO_CFGR, reg_value);
	write_pio(pio, (value ? PIO_SODR : PIO_CODR), mask);
#else
	write_pio(pio, ((config & PIO_OPENDRAIN) ? PIO_MDER : PIO_MDDR), mask);

#ifdef CONFIG_CPU_HAS_PIO3
	pio3_set_drvstr(pin, config);
	pio3_set_slewrate(pin, config);
#endif

	pio_set_gpio_output(pin, value);
#endif

	return 0;
}

int pio_set_value(unsigned pin, int value)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);

	if (pio >= AT91C_NUM_PIO)
		return -1;

	write_pio(pio, (value ? PIO_SODR : PIO_CODR), mask);

	return 0;
}

int pio_get_value(unsigned pin)
{
	unsigned pio = pin_to_controller(pin);
	unsigned mask = pin_to_mask(pin);
	unsigned int pdsr;

	if (pio >= AT91C_NUM_PIO)
		return -1;

	pdsr = read_pio(pio, PIO_PDSR);

	return ((pdsr & mask) != 0);
}


int pio_configure(const struct pio_desc *pio_desc)
{
	unsigned pio, pin = 0;

	if (pio_desc == 0) return 0;

	/*
	 * Sets all the pio muxing of the corresponding device as defined
	 * in its platform_data struct
	 */
	while (pio_desc->pin_name) {
		pio = pin_to_controller(pio_desc->pin_num);
		if (pio >= AT91C_NUM_PIO) {
			return 0;
		} else if (pio_desc->type == PIO_PERIPH_A) {
			pio_set_a_periph(pio_desc->pin_num,
						pio_desc->attribute);
		} else if (pio_desc->type == PIO_PERIPH_B) {
			pio_set_b_periph(pio_desc->pin_num,
						pio_desc->attribute);
		} else if (pio_desc->type == PIO_PERIPH_C) {
			pio_set_c_periph(pio_desc->pin_num,
						pio_desc->attribute);
		} else if (pio_desc->type == PIO_PERIPH_D) {
			pio_set_d_periph(pio_desc->pin_num,
						pio_desc->attribute);
		} else if (pio_desc->type == PIO_PERIPH_E) {
			pio_set_e_periph(pio_desc->pin_num,
						pio_desc->attribute);
		} else if (pio_desc->type == PIO_PERIPH_F) {
			pio_set_f_periph(pio_desc->pin_num,
						pio_desc->attribute);
		} else if (pio_desc->type == PIO_PERIPH_G) {
			pio_set_g_periph(pio_desc->pin_num,
						pio_desc->attribute);
		} else if (pio_desc->type == PIO_INPUT) {
			pio_set_gpio_input(pio_desc->pin_num,
						pio_desc->attribute);
		} else if (pio_desc->type == PIO_OUTPUT) {
			pio_config_gpio_output(pio_desc->pin_num,
						pio_desc->attribute,
						pio_desc->default_value);
		} else {
			return 0;
		}
		
		++pin;
		++pio_desc;
	}

	return pin;
}
