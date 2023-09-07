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

#ifndef __SOC_MEDIATEK_COMMON_MTCMOS_H__
#define __SOC_MEDIATEK_COMMON_MTCMOS_H__

void mtcmos_audio_power_on(void);
void mtcmos_display_power_on(void);

void mtcmos_protect_display_bus(void);

#endif /* __SOC_MEDIATEK_COMMON_MTCMOS_H__ */
