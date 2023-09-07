/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2014 Intel Corporation
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

#include <soc/ramstage.h>

void board_silicon_USB2_override(SILICON_INIT_UPD *params)
{
	if (SocStepping() >= SocD0) {

		//Follow Intel recommendation to set
		//BSW D-stepping PERPORTRXISET 2 (low strength)
		params->D0Usb2Port0PerPortRXISet = 2;
		params->D0Usb2Port1PerPortRXISet = 2;
		params->D0Usb2Port2PerPortRXISet = 2;
		params->D0Usb2Port3PerPortRXISet = 2;
		params->D0Usb2Port4PerPortRXISet = 2;
	}
}
