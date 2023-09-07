/*
 * This file is part of the coreboot project.
 *
 * Copyright (c) 2018 Qualcomm Technologies
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2 and
 * only version 2 as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */
#include <types.h>

#ifndef _SDM845_USB_H_
#define _SDM845_USB_H_

/* QSCRATCH_GENERAL_CFG register bit offset */
#define PIPE_UTMI_CLK_SEL			BIT(0)
#define PIPE3_PHYSTATUS_SW			BIT(3)
#define PIPE_UTMI_CLK_DIS			BIT(8)

/* Global USB3 Control  Registers */
#define DWC3_GUSB3PIPECTL_DELAYP1TRANS		BIT(18)
#define DWC3_GUSB3PIPECTL_UX_EXIT_IN_PX		BIT(27)
#define DWC3_GCTL_PRTCAPDIR(n)			((n) << 12)
#define DWC3_GCTL_PRTCAP_OTG			3
#define DWC3_GCTL_PRTCAP_HOST			1

/* Global USB2 PHY Configuration Register */
#define DWC3_GUSB2PHYCFG_USBTRDTIM(n)		((n) << 10)
#define DWC3_GUSB2PHYCFG_USB2TRDTIM_MASK	DWC3_GUSB2PHYCFG_USBTRDTIM(0xf)
#define DWC3_GUSB2PHYCFG_PHYIF(n)		((n) << 3)
#define DWC3_GUSB2PHYCFG_PHYIF_MASK		DWC3_GUSB2PHYCFG_PHYIF(1)
#define USBTRDTIM_UTMI_8_BIT			9
#define UTMI_PHYIF_8_BIT			0

#define DWC3_GCTL_SCALEDOWN(n)			((n) << 4)
#define DWC3_GCTL_SCALEDOWN_MASK		DWC3_GCTL_SCALEDOWN(3)
#define DWC3_GCTL_DISSCRAMBLE			(1 << 3)
#define DWC3_GCTL_U2EXIT_LFPS			(1 << 2)
#define DWC3_GCTL_DSBLCLKGTNG			(1 << 0)

#define PORT_TUNE1_MASK				0xf0

/* QUSB2PHY_PWR_CTRL1 register related bits */
#define POWER_DOWN				BIT(0)

/* DEBUG_CTRL2 register value to program VSTATUS MUX for PHY status */
#define DEBUG_CTRL2_MUX_PLL_LOCK_STATUS		0x4

/* STAT5 register bits */
#define VSTATUS_PLL_LOCK_STATUS_MASK		BIT(0)

/* QUSB PHY register values */
#define QUSB2PHY_PLL_ANALOG_CONTROLS_TWO	0x03
#define QUSB2PHY_PLL_CLOCK_INVERTERS		0x7c
#define QUSB2PHY_PLL_CMODE			0x80
#define QUSB2PHY_PLL_LOCK_DELAY			0x0a
#define QUSB2PHY_PLL_DIGITAL_TIMERS_TWO		0x19
#define QUSB2PHY_PLL_BIAS_CONTROL_1		0x40
#define QUSB2PHY_PLL_BIAS_CONTROL_2		0x20
#define QUSB2PHY_PWR_CTRL2			0x21
#define QUSB2PHY_IMP_CTRL1			0x0
#define QUSB2PHY_IMP_CTRL2			0x58
#define QUSB2PHY_PORT_TUNE1			0x30
#define QUSB2PHY_PORT_TUNE2			0x29
#define QUSB2PHY_PORT_TUNE3			0xca
#define QUSB2PHY_PORT_TUNE4			0x04
#define QUSB2PHY_PORT_TUNE5			0x03
#define QUSB2PHY_CHG_CTRL2			0x0

/* USB3PHY_PCIE_USB3_PCS_PCS_STATUS bit */
#define USB3_PCS_PHYSTATUS		BIT(6)

struct usb_board_data {
	/* Register values going to override from the boardfile */
	u32 pll_bias_control_2;
	u32 imp_ctrl1;
	u32 port_tune1;
};

struct qmp_phy_init_tbl {
	u32 *address;
	u32 val;
};

void setup_usb_host0(struct usb_board_data *data);
void setup_usb_host1(struct usb_board_data *data);
/* Call reset_ before setup_ */
void reset_usb0(void);
void reset_usb1(void);

#endif /* _SDM845_USB_H_ */
