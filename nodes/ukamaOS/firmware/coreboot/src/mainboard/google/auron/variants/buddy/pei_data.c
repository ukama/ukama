/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2014 Google Inc.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation; version 2 of
 * the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <stdint.h>
#include <soc/gpio.h>
#include <soc/pei_data.h>
#include <soc/pei_wrapper.h>

void mainboard_fill_pei_data(struct pei_data *pei_data)
{
	pei_data->ec_present = 1;

	/* P0: Side USB3.0 port, USB3S1 */
	pei_data_usb2_port(pei_data, 0, 0x0150, 1, 0,
			   USB_PORT_INTERNAL);
	/* P1: Rear USB3.0 port, USB3R1 */
	pei_data_usb2_port(pei_data, 1, 0x0040, 1, 0,
			   USB_PORT_INTERNAL);
	/* P2: Rear USB3.0 port, USB3R2 */
	pei_data_usb2_port(pei_data, 2, 0x0080, 1, 1,
			   USB_PORT_INTERNAL);
	/* P3: Card Rearder, CRS1 */
	pei_data_usb2_port(pei_data, 3, 0x0040, 1, USB_OC_PIN_SKIP,
			   USB_PORT_INTERNAL);
	/* P4: Rear USB2.0 port, USB2R1 */
	pei_data_usb2_port(pei_data, 4, 0x0040, 1, 2,
			   USB_PORT_INTERNAL);
	/* P5: 2D Camera */
	pei_data_usb2_port(pei_data, 5, 0x0000, 1, USB_OC_PIN_SKIP,
			   USB_PORT_INTERNAL);
	/* P6: VP8 */
	pei_data_usb2_port(pei_data, 6, 0x0150, 1, USB_OC_PIN_SKIP,
			   USB_PORT_MINI_PCIE);
	/* P7: WLAN & BT */
	pei_data_usb2_port(pei_data, 7, 0x0000, 1, USB_OC_PIN_SKIP,
			   USB_PORT_MINI_PCIE);

	/* P1: Side USB3.0 port, USB3S1 */
	pei_data_usb3_port(pei_data, 0, 1, 0, 0);
	/* P2: Rear USB3.0 port, USB3R1 */
	pei_data_usb3_port(pei_data, 1, 1, 0, 0);
	/* P3: Rear USB3.0 port, USB3R2 */
	pei_data_usb3_port(pei_data, 2, 1, 1, 0);
	/* P4: Card Rearder, CRS1 */
	pei_data_usb3_port(pei_data, 3, 1, USB_OC_PIN_SKIP, 0);
}
