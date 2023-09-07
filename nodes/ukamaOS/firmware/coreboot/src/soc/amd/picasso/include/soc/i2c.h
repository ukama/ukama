/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018 Advanced Micro Devices, Inc.
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

#ifndef __PICASSO_I2C_H__
#define __PICASSO_I2C_H__

#include <soc/gpio.h>

struct soc_amd_i2c_save {
	uint32_t control_value;
	uint8_t mux_value;
};

#define GPIO_I2C2_SCL		BIT(2)
#define GPIO_I2C3_SCL		BIT(3)
#define GPIO_I2C_MASK		(BIT(2) | BIT(3))

#define I2C2_SCL_PIN			GPIO_113
#define I2C3_SCL_PIN			GPIO_19

#define GPIO_I2C2_ADDRESS		GPIO_BANK1_CONTROL(I2C2_SCL_PIN)
#define GPIO_I2C3_ADDRESS		GPIO_BANK0_CONTROL(I2C3_SCL_PIN)
#define GPIO_SCL_HIGH			0
#define GPIO_SCL_LOW			GPIO_OUTPUT_ENABLE

#define I2C2_SCL_PIN_IOMUX_GPIOxx	GPIO_113_IOMUX_GPIOxx
#define I2C3_SCL_PIN_IOMUX_GPIOxx	GPIO_19_IOMUX_GPIOxx

void sb_reset_i2c_slaves(void);

#endif /* __PICASSO_I2C_H__ */
