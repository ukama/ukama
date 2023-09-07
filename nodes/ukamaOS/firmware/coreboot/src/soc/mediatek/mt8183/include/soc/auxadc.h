/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 MediaTek Inc.
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

#ifndef _MTK_ADC_H
#define _MTK_ADC_H

#include <stdint.h>

typedef struct mtk_auxadc_regs {
	uint32_t con0;
	uint32_t con1;
	uint32_t con1_set;
	uint32_t con1_clr;
	uint32_t con2;
	uint32_t data[16];
	uint32_t reserved[16];
	uint32_t misc;
} mtk_auxadc_regs;

/* Return voltage in uVolt */
int auxadc_get_voltage(unsigned int channel);
#endif
