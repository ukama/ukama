/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2014, Atmel Corporation
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
#include "debug.h"
#include "macb.h"
#include "pmc.h"

/*
 * EMAC & GMAC Register
 */
#define	MACB_NCR	0x00	/* Network Control Register */
#define	MACB_NCFGR	0x04	/* Network Configuration Register */
#define	MACB_NSR	0x08	/* Network Status Register */
#define	MACB_MAN	0x34	/* Phy Maintenance Register */
#define	MACB_ID		0xfc

/*
 * MACB_NCR: Network Control Register
 */
#define	MACB_NCR_LB		(0x01 << 0)
#define	MACB_NCR_LLB		(0x01 << 1)
#define	MACB_NCR_RE		(0x01 << 2)
#define	MACB_NCR_TE		(0x01 << 3)
#define	MACB_NCR_MPE		(0x01 << 4)
#define	MACB_NCR_CLRSTAT	(0x01 << 5)
#define	MACB_NCR_INCSTAT	(0x01 << 6)
#define	MACB_NCR_WESTAT		(0x01 << 7)
#define	MACB_NCR_BP		(0x01 << 8)
#define	MACB_NCR_TSTART		(0x01 << 9)
#define	MACB_NCR_THALT		(0x01 << 10)

/*
 * MACB_NCFGR: Network Configuration Register
 */
#define MACB_NCFGR_SPD		(0x01 << 0)
#define MACB_NCFGR_FD		(0x01 << 1)
#define MACB_NCFGR_JFRAME	(0x01 << 3)
#define MACB_NCFGR_CAF		(0x01 << 4)
#define MACB_NCFGR_NBC		(0x01 << 5)
#define MACB_NCFGR_MTI		(0x01 << 6)
#define MACB_NCFGR_UNI		(0x01 << 7)
#define MACB_NCFGR_BIG		(0x01 << 8)
#define EMAC_NCFGR_CLK_MASK	(0x03 << 10)
#define		EMAC_NCFGR_CLK_MCK_8	(0x00 << 10)
#define		EMAC_NCFGR_CLK_MCK_16	(0x01 << 10)
#define		EMAC_NCFGR_CLK_MCK_32	(0x02 << 10)
#define		EMAC_NCFGR_CLK_MCK_64	(0x03 << 10)
#define	GMAC_NCFGR_CLK_MASK	(0x07 << 18)
#define		GMAC_NCFGR_CLK_MCK_8	(0x00 << 18)
#define		GMAC_NCFGR_CLK_MCK_16	(0x01 << 18)
#define		GMAC_NCFGR_CLK_MCK_32	(0x02 << 18)
#define		GMAC_NCFGR_CLK_MCK_48	(0x03 << 18)
#define		GMAC_NCFGR_CLK_MCK_64	(0x04 << 18)
#define		GMAC_NCFGR_CLK_MCK_96	(0x05 << 18)

#define MACB_NCFGR_RTY		(0x01 << 12)
#define MACB_NCFGR_FAE		(0x01 << 13)

/*
 * MACB_NSR: Network Status Register
 */
#define MACB_NSR_MDIO		(0x01 << 1)
#define MACB_NSR_IDLE		(0x01 << 2)

/*
 * MACB_MAN: PHY Maintenance Register
 */
#define MACB_MAN_DATA_MASK	(0xffff)
#define		MACB_MAN_DADA(value)	((value) << 0)
#define MACB_MAN_CODE_MASK	(0x03 << 16)
#define		MACB_MAN_CODE		(0x02 << 16)
#define MACB_MAN_REGA_MASK	(0x1f << 18)
#define		MACB_MAN_REGA(value)	((value) << 18)
#define MACB_MAN_PHYA_MASK	(0x1f << 23)
#define		MACB_MAN_PHYA(value)	((value) << 23)
#define MACB_MAN_RW_MASK	(0x03 << 28)
#define		MACB_MAN_RW_WRITE	(0x01 << 28)
#define		MACB_MAN_RW_READ	(0x02 << 28)
#define MACB_MAN_SOF_MASK	(0x03 << 30)
#define		MACB_MAN_SOF		(0x01 << 30)

/*
 * PHY Standard Register Map
 */
#define MII_BMCR		0x00	/* Basic Control Register */
#define MII_BMSR		0x01	/* Basic Status Register  */
#define MII_PHYSID1		0x02	/* PHY Identifier 1 */
#define MII_PHYSID2		0x03	/* PHY Identifier 2 */
#define MII_ADVERTISE		0x04	/* Auto-Negotiation Advertisement */
#define MII_LPA			0x05	/* Auto-Negotiation LPA */
#define MII_EXPANSION		0x06	/* Auto-Negotiation Expansion */
#define MII_NEXT_PAGE		0x07	/* Auto-Negotiation Next Page */
#define MII_LPA_NEXT		0x08	/* Auto-Negotiation LPA Next Page */

/*
 * Basic Control Register Description
 */
#define BMCR_CTST		(0x01 << 7)	/* Collision Test */
#define BMCR_FULLDPLX		(0x01 << 8)	/* Full Duplex Mode */
#define BMCR_ANRESTART		(0x01 << 9)	/* Auto negotiation Restart */
#define BMCR_ISOLATE		(0x01 << 10)	/* Isolation of PHY from MII */
#define BMCR_PDOWN		(0x01 << 11)	/* Power Down */
#define BMCR_ANENABLE           (0x01 << 12)	/* Enable Auto Negotiation */
#define BMCR_SPEED100		(0x01 << 13)	/* Speed Select 100Mbps */
#define BMCR_LOOPBACK		(0x01 << 14)	/* TXD loopback bits */
#define BMCR_RESET		(0x01 << 15)	/* Software Reset */

#define PHY_ID_NUMBER		(0x0022)

struct mii_bus {
	const char *name;
	void *reg_base;
	unsigned int phy_addr;
};

static inline unsigned int macb_read(void *base, unsigned int offset)
{
	return readl(base + offset);
}

static inline void macb_write(void *base,
				unsigned int offset,
				unsigned int value)
{
	writel(value, base + offset);
}

static int macb_mdio_read(struct mii_bus *bus,
				unsigned int regnum,
				unsigned int *value)
{
	unsigned int timeout = 10000;
	unsigned int reg;

	reg = (MACB_MAN_SOF
			| MACB_MAN_RW_READ
			| MACB_MAN_PHYA(bus->phy_addr)
			| MACB_MAN_REGA(regnum)
			| MACB_MAN_CODE);
	macb_write(bus->reg_base, MACB_MAN, reg);

	do {
		reg = macb_read(bus->reg_base, MACB_NSR);
	} while ((--timeout) && (!(reg & MACB_NSR_IDLE)));

	if (!timeout) {
		dbg_loud("MDIO: Timeout to read\n");
		return -1;
	}

	reg = macb_read(bus->reg_base, MACB_MAN);
	*value = MACB_MAN_DATA_MASK & reg;

	return 0;
}

static int macb_mdio_write(struct mii_bus *bus,
				unsigned int regnum,
				unsigned int value)
{
	unsigned int timeout = 10000;
	unsigned int reg;

	reg = (MACB_MAN_SOF
			| MACB_MAN_RW_WRITE
			| MACB_MAN_PHYA(bus->phy_addr)
			| MACB_MAN_REGA(regnum)
			| MACB_MAN_DADA(value)
			| MACB_MAN_CODE);
	macb_write(bus->reg_base, MACB_MAN, reg);

	do {
		reg = macb_read(bus->reg_base, MACB_NSR);
	} while ((--timeout) && (!(reg & MACB_NSR_IDLE)));

	if (!timeout) {
		dbg_loud("MDIO: Timeout to write\n");
		return -1;
	}

	return 0;
}

static void macb_enable_managementport(struct mii_bus *bus, unsigned char on)
{
	unsigned int reg = macb_read(bus->reg_base, MACB_NCR);

	if (on)
		macb_write(bus->reg_base, MACB_NCR, reg | MACB_NCR_MPE);
	else
		macb_write(bus->reg_base, MACB_NCR, (reg & ~MACB_NCR_MPE));
}

static unsigned int macb_findphy(struct mii_bus *bus)
{
	unsigned int phy_address = bus->phy_addr;
	unsigned int value;
	unsigned int i;
	unsigned int rc = 0xff;

	if (macb_mdio_read(bus, MII_PHYSID1, &value)) {
		dbg_loud("MACB: Failed to read MII_PHYID1\n");
		return rc;
	}

	rc = phy_address;
	if (value != PHY_ID_NUMBER) {
		rc = 0xff;
		for (i = 0; i < 32; i++) {
			phy_address = (phy_address + 1) & 0x1f;
			bus->phy_addr = phy_address;
			if (macb_mdio_read(bus, MII_PHYSID1, &value))
				dbg_loud("MACB: Failed to read MII_PHYID1\n");
			if (value == PHY_ID_NUMBER) {
				rc = phy_address;
				break;
			}
		}
	}

	return rc;
}

static int macb_is_gem(struct mii_bus *bus)
{
	unsigned int reg = macb_read(bus->reg_base, MACB_ID);

	return (((reg >> 16) & 0xffff) == 0x02) ? 1 : 0;
}

static unsigned int gem_mdc_clk_div(void)
{
	unsigned int clk_div;
	unsigned int mck = at91_get_ahb_clock();

	if (mck < 20000000)
		clk_div = GMAC_NCFGR_CLK_MCK_8;
	else if (mck < 40000000)
		clk_div = GMAC_NCFGR_CLK_MCK_16;
	else if (mck < 80000000)
		clk_div = GMAC_NCFGR_CLK_MCK_32;
	else if (mck < 120000000)
		clk_div = GMAC_NCFGR_CLK_MCK_48;
	else if (mck < 160000000)
		clk_div = GMAC_NCFGR_CLK_MCK_64;
	else
		clk_div = GMAC_NCFGR_CLK_MCK_96;

	return clk_div;
}

static unsigned int macb_mdc_clk_div(void)
{
	unsigned int clk_div;
	unsigned int mck = at91_get_ahb_clock();

	if (mck < 20000000)
		clk_div = EMAC_NCFGR_CLK_MCK_8;
	else if (mck < 40000000)
		clk_div = EMAC_NCFGR_CLK_MCK_16;
	else if (mck < 80000000)
		clk_div = EMAC_NCFGR_CLK_MCK_32;
	else
		clk_div = EMAC_NCFGR_CLK_MCK_64;

	return clk_div;
}

static int macb_set_mdc_clk(struct mii_bus *bus)
{
	unsigned int config;
	unsigned int clk_div;

	if (macb_is_gem(bus))
		clk_div = gem_mdc_clk_div();
	else
		clk_div = macb_mdc_clk_div();

	config = macb_read(bus->reg_base, MACB_NCFGR);
	if (macb_is_gem(bus))
		config &= ~EMAC_NCFGR_CLK_MASK;
	else
		config &= ~GMAC_NCFGR_CLK_MASK;
	config |= clk_div;
	macb_write(bus->reg_base, MACB_NCFGR, config);

	return 0;
}

static int phy_software_reset(struct mii_bus *bus)
{
	unsigned int value;
	int timeout = 10;

	value = BMCR_RESET;
	if (macb_mdio_write(bus, MII_BMCR, value)) {
		dbg_loud("MACB: Failed to write MII_BMCR\n");
		return -1;
	}

	do {
		if (macb_mdio_read(bus, MII_BMCR, &value)) {
			dbg_loud("MACB: Failed to read MII_BMCR\n");
			return -1;
		}
	} while ((--timeout) && (value & BMCR_RESET));

	return 0;
}

static int phy_power_down(struct mii_bus *bus)
{
	unsigned int value;

	if (macb_mdio_read(bus, MII_BMCR, &value)) {
		dbg_loud("MACB: Failed to read MII_BMCR\n");
		return -1;
	}

	if (macb_mdio_write(bus, MII_BMCR, (value | BMCR_PDOWN))) {
		dbg_loud("MACB: Failed to write MII_BMCR\n");
		return -1;
	}

	return 0;
}

static int phy_power_down_mode(struct mii_bus *bus)
{
	unsigned int phy_addr;
	int ret = 0;

	macb_set_mdc_clk(bus);

	macb_enable_managementport(bus, 1);

	phy_addr = macb_findphy(bus);
	if (phy_addr == 0xff) {
		ret = -1;
		dbg_loud("PHY: Failed to find PHY\n");
		goto error;
	}

	bus->phy_addr = phy_addr;

	if (phy_software_reset(bus)) {
		ret = -1;
		dbg_loud("PHY: Failed to reset PHY\n");
		goto error;
	}

	if (phy_power_down(bus)) {
		ret = -1;
		dbg_loud("PHY: Failed to enter power down mode\n");
		goto error;
	}

	dbg_loud("PHY: %s: Enter power down mode\n", bus->name);

error:
	macb_enable_managementport(bus, 0);
	return ret;
}

int phys_enter_power_down(void)
{
	struct mii_bus macb_mii_bus;
	unsigned int base_addr;

#if defined(CONFIG_MAC0_PHY)
	base_addr = at91_eth0_hw_init();

	macb_mii_bus.name = "ETH0 PHY";
	macb_mii_bus.reg_base = (void *)base_addr;
	macb_mii_bus.phy_addr = 1;

	if (phy_power_down_mode(&macb_mii_bus)) {
		dbg_loud("%s: Failed to enter power down mode\n",
						macb_mii_bus.name);
	}
#endif

#if defined(CONFIG_MAC1_PHY)
	base_addr = at91_eth1_hw_init();

	macb_mii_bus.name = "ETH1 PHY";
	macb_mii_bus.reg_base = (void *)base_addr;
	macb_mii_bus.phy_addr = 1;

	if (phy_power_down_mode(&macb_mii_bus)) {
		dbg_loud("%s: Failed to enter power down mode\n",
						macb_mii_bus.name);
	}
#endif

	at91_disable_mac_clock();

	return 0;
}
