/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2014 Google Inc.
 * Copyright (C) 2017 Intel Corporation.
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

#ifndef _SOC_MSR_H_
#define _SOC_MSR_H_

#include <intelblocks/msr.h>

#define MSR_PIC_MSG_CONTROL		0x2e
#define MSR_VR_CURRENT_CONFIG		0x601
#define MSR_PL3_CONTROL			0x615
#define MSR_VR_MISC_CONFIG2		0x636
#define MSR_PLATFORM_POWER_LIMIT	0x65c

#endif
