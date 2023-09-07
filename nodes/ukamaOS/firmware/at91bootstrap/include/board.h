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
#ifndef __BOARD_H__
#define __BOARD_H__

#ifdef CONFIG_AT91SAM9260EK
#include "at91sam9260ek.h"
#endif

#ifdef CONFIG_AT91SAM9XEEK
#include "at91sam9xeek.h"
#endif

#ifdef CONFIG_AT91SAM9G10EK
#include "at91sam9g10ek.h"
#endif

#ifdef CONFIG_AT91SAM9G20EK
#include "at91sam9g20ek.h"
#endif

#ifdef CONFIG_AT91SAM9M10G45EK
#include "at91sam9m10g45ek.h"
#endif

#ifdef CONFIG_AT91SAM9261EK
#include "at91sam9261ek.h"
#endif

#ifdef CONFIG_AT91SAM9263EK
#include "at91sam9263ek.h"
#endif

#ifdef CONFIG_AT91SAM9RLEK
#include "at91sam9rlek.h"
#endif

#ifdef CONFIG_AT91SAM9X5EK
#include "at91sam9x5ek.h"
#endif

#ifdef CONFIG_AT91SAM9N12EK
#include "at91sam9n12ek.h"
#endif

#ifdef CONFIG_SAM9X60EK
#include "sam9x60ek.h"
#endif

#ifdef CONFIG_SAM9X60_DDR2_SIP_EB
#include "sam9x60_ddr2_sip_eb.h"
#endif

#ifdef CONFIG_SAM9X60_SDR_SIP_EB
#include "sam9x60_sdr_sip_eb.h"
#endif

#ifdef CONFIG_SAMA5D3XEK
#include "sama5d3xek.h"
#endif

#ifdef CONFIG_SAMA5D3_XPLAINED
#include "sama5d3_xplained.h"
#endif

#ifdef CONFIG_SAMA5D3X_CMP
#include "sama5d3x_cmp.h"
#endif

#ifdef CONFIG_SAMA5D4EK
#include "sama5d4ek.h"
#endif

#ifdef CONFIG_SAMA5D4_XPLAINED
#include "sama5d4_xplained.h"
#endif

#ifdef CONFIG_SAMA5D2_PTC_EK
#include "sama5d2_ptc_ek.h"
#endif

#ifdef CONFIG_SAMA5D2_XPLAINED
#include "sama5d2_xplained.h"
#endif

#ifdef CONFIG_SAMA5D27_SOM1_EK
#include "sama5d27_som1_ek.h"
#endif

#ifdef CONFIG_SAMA5D27_WLSOM1_EK
#include "sama5d27_wlsom1_ek.h"
#endif

#ifdef CONFIG_SAMA5D27_UKAMA_ANODE
#include "sama5d27_ukama_anode.h"
#endif

#ifdef CONFIG_SAMA5D2_LPDDR2SIP_VB
#include "sama5d2_lpddr2sip_vb.h"
#endif

#ifdef CONFIG_SAMA5D2_ICP
#include "sama5d2_icp.h"
#endif

#ifdef CONFIG_SAMA7G5EK
#include "sama7g5ek.h"
#endif

#include "contrib_board.h"

/*
 * Functions Prototype
 */
extern void hw_init(void);

extern void nandflash_hw_init(void);

extern void at91_spi0_hw_init(void);

extern void at91_qspi_hw_init(void);

extern void at91_mci0_hw_init(void);
extern void at91_mci1_hw_init(void);
extern void at91_mci2_hw_init(void);

extern void at91_sdhc_hw_init(void);

extern void at91_board_set_dtb_name(char *of_name);

extern void norflash_hw_init(void);

extern char *board_override_cmd_line(void);

/* PM Support */
extern unsigned int at91_twi0_hw_init(void);
extern unsigned int at91_twi1_hw_init(void);
extern unsigned int at91_twi2_hw_init(void);
extern unsigned int at91_twi3_hw_init(void);
extern void at91_board_config_twi_bus(void);

extern unsigned int at91_eth0_hw_init(void);
extern unsigned int at91_eth1_hw_init(void);
extern void at91_disable_mac_clock(void);

extern void at91_disable_smd_clock(void);

extern int at91_board_act8865_set_reg_voltage(void);

#endif /* #ifndef __BOARD_H__ */
