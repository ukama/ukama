/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2013 Google Inc.
 * Copyright (C) 2013 Sage Electronic Engineering, LLC.
 * Copyright (C) 2014 Intel Corporation
 * Copyright (C) 2017-present Facebook, Inc.
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

#include <console/console.h>
#include <soc/gpio.h>
#include <soc/romstage.h>
#include <pc80/mc146818rtc.h>
#include <drivers/intel/fsp1_0/fsp_util.h>
#include <soc/intel/fsp_baytrail/chip.h>
#include <cpu/x86/tsc.h>

/**
 * /brief mainboard call for setup that needs to be done before fsp init
 *
 */
void early_mainboard_romstage_entry(void)
{

}

/**
 * Get function disables - most of these will be done automatically
 * @param fd_mask
 * @param fd2_mask
 */
void get_func_disables(uint32_t *fd_mask, uint32_t *fd2_mask)
{

}

/**
 * /brief mainboard call for setup that needs to be done after fsp init
 *
 */
void late_mainboard_romstage_entry(void)
{

}

static const MEMORY_DOWN_DATA minnowmax_memory_config[] = {
        /* 1066 */
        {
        .EnableMemoryDown = 1,
        .DRAMSpeed = 1,      /* DRAM Speed: 0=800, 1=1066, 2=1333, 3=1600*/
        .DRAMType = 1,       /* DRAM Type: 0=DDR3, 1=DDR3L, 2=DDR3U, 4=LPDDR2, 5=LPDDR3, 6=DDR4*/
        .DIMM0Enable = 1,    /* DIMM 0 Enable */
        .DIMM1Enable = 0,    /* DIMM 1 Enable */
        .DIMMDWidth = 1,     /* DRAM device data width: 0=x8, 1=x16, 2=x32*/
        .DIMMDensity = 1,    /* DRAM device data density: 0=1Gb, 1=2Gb, 2=4Gb, 3=8Gb */
        .DIMMBusWidth = 3,   /* DIMM Bus Width: 0=8bit, 1=16bit, 2=32bit, 3=64bit */
        .DIMMSides = 0,      /* Ranks Per DIMM: 0=1rank, 1=2rank */
        .DIMMtCL = 11,        /* tCL */
        .DIMMtRPtRCD = 11,    /* tRP and tRCD in DRAM clk - 5:12.5ns, 6:15ns, etc. */
        .DIMMtWR = 12,        /* tWR in DRAM clk  */
        .DIMMtWTR = 6,       /* tWTR in DRAM clk */
        .DIMMtRRD = 6,       /* tRRD in DRAM clk */
        .DIMMtRTP = 6,       /* tRTP in DRAM clk */
        .DIMMtFAW = 20,      /* tFAW in DRAM clk */
        },
        /* 1333 */
        {
        .EnableMemoryDown = 1,
        .DRAMSpeed = 2,      /* DRAM Speed: 0=800, 1=1066, 2=1333, 3=1600*/
        .DRAMType = 1,       /* DRAM Type: 0=DDR3, 1=DDR3L, 2=DDR3U, 4=LPDDR2, 5=LPDDR3, 6=DDR4*/
        .DIMM0Enable = 1,    /* DIMM 0 Enable */
        .DIMM1Enable = 0,    /* DIMM 1 Enable */
        .DIMMDWidth = 1,     /* DRAM device data width: 0=x8, 1=x16, 2=x32*/
        .DIMMDensity = 1,    /* DRAM device data density: 0=1Gb, 1=2Gb, 2=4Gb, 3=8Gb */
        .DIMMBusWidth = 3,   /* DIMM Bus Width: 0=8bit, 1=16bit, 2=32bit, 3=64bit */
        .DIMMSides = 0,      /* Ranks Per DIMM: 0=1rank, 1=2rank */
        .DIMMtCL = 9,        /* tCL */
        .DIMMtRPtRCD = 9,    /* tRP and tRCD in DRAM clk - 5:12.5ns, 6:15ns, etc. */
        .DIMMtWR = 10,       /* tWR in DRAM clk  */
        .DIMMtWTR = 5,       /* tWTR in DRAM clk */
        .DIMMtRRD = 4,       /* tRRD in DRAM clk */
        .DIMMtRTP = 5,       /* tRTP in DRAM clk */
        .DIMMtFAW = 30,      /* tFAW in DRAM clk */
        }
};

void romstage_fsp_rt_buffer_callback(FSP_INIT_RT_BUFFER *FspRtBuffer)
{
        UPD_DATA_REGION *UpdData = FspRtBuffer->Common.UpdDataRgnPtr;
        u8 use_xhci = UpdData->PcdEnableXhci;
        int is_1333_sku;
        is_1333_sku = !!(tsc_freq_mhz() >= 1700);

        printk(BIOS_INFO, "DDR Using %d MHz DDR3 settings.\n",
                is_1333_sku ? 1333 : 1066);

	printk(BIOS_NOTICE, "This is the Ukama ComV1 compute board.\n");
        UpdData->PcdMemoryParameters = minnowmax_memory_config[is_1333_sku];
        UpdData->PcdMemoryParameters.DIMMDensity
                += (DIMM_DENSITY_4G_BIT - DIMM_DENSITY_2G_BIT);
        printk(BIOS_NOTICE, "%s GB Ukama ComV1 detected.\n",
                        "2 / 4");
        /* Disable 2nd DIMM */
        //UpdData->PcdMrcInitSPDAddr2 = 0x00; /* cannot use SPD_ADDR_DISABLED at this point */

        /* Update XHCI UPD value if required */
        get_option(&use_xhci, "use_xhci_over_ehci");
        if ((use_xhci < 2) && (use_xhci != UpdData->PcdEnableXhci)) {
                UpdData->PcdEnableXhci = use_xhci;
                printk(FSP_INFO_LEVEL, "Xhci updated from CMOS:\t\t\t%s\n",
                        UpdData->PcdEnableXhci?"Enabled":"Disabled");
        }

        return;
}

