/* ----------------------------------------------------------------------------
 *         Microchip Technology AT91Bootstrap project
 * ----------------------------------------------------------------------------
 * Copyright (c) 2018, Microchip Technology Inc. and its subsidiaries
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Microchip's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY MICROCHIP "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT ARE
 * DISCLAIMED. IN NO EVENT SHALL MICROCHIP BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

#ifndef __UMCTL2_H__
#define __UMCTL2_H__

struct umctl2_config_state
{
	/* PHY callbacks to PHY driver */
		/* Pointer to PHY init function */
	void (*phy_init)(void *);
		/* Pointer to phy init done function */
	int (*phy_idone)(void);
		/* Pointer to PHY start function */
	int (*phy_start)(void);
		/* Pointer to PHY training function */
	int (*phy_train)(void);

	/* Policies configuration */
		/* pageclose mechanism: precharge banks only after pageclose_timer
		expires, after there are no more page hit transactions in CAM.
		Disabling this means it's open-page policy (dont precharge).
		Can change to closed page policy from HIF (close immediately).
		This mechanism is mid-way between open and closed */
	unsigned int pageclose;
		/* timer for pageclose policy: wait this timer until precharge
		after there are no more page hit transactions waiting in CAM */
	unsigned int pageclose_timer;
		/* DEBUG ONLY: precharge or not after colission cases:
		r after w, w after r, w after w with disabled write combine */
	unsigned int dis_collision_page_opt;
		/* DEBUG ONLY: disable write combine */
	unsigned int dis_wc;

	/* bitmap of AXI ports opened in UCMTL2_MP.
	Some are needed for DMA, others for CPU or other busses... depends on architecture
	Port 0 : CPU : 128 bit, single quue for read, single queue write.
	Port 1 : HSS : 64 bit, single read, single write. this is a bridge to AHB/APB
	Port 2 : PSS: 64 bit, dual queue for read (red/blue), single queue write.
	Port 3 : ISS: 64 bit: single read/single write : camera et al.
	Port 4 : MSS: SDMMC and GMAC ... single read/single write
	 */
	unsigned int axi_port_bitmap;

	/* Qos regions definitions */
	/* this array represents the last of region0 QOS value for indexed port for read requests */
	unsigned int port_x_read_region0_last[5];
	/* this array represents the last of region1 QOS value for indexed port for read requests
	   only for dual queue ports - two regions. region2 will end at value 0xF */
	unsigned int port_x_read_region1_last[5];
	/* this array represents the last of region0 QOS value for indexed port for write requests */
	unsigned int port_x_write_region0_last[5];
	/* this array represents the last of region1 QOS value for indexed port for write requests */
	unsigned int port_x_write_region1_last[5];

	/* Enable ordered read/writes. If set to 1, preserves the
	 * ordering between read transaction and write transaction
	 * issued to the same address, on a given port.
	 */
	unsigned int port_x_rdwr_ordered_en;

	/* Transaction stores configuration*/
	/* Read transaction store */
	/* number of entries in low prio store (LPR): */
	unsigned int lpr_num_entries;
	/* number of entries in high prio store (HPR): */
		/* this value is MEMC_NO_OF_ENTRY - lpr_num_entries */
	/* Write transaction store */
		/* this store doesn't have a split, all in a single queue */

	/* Store starvation configuration */
		/* read queues starvation configuration */
		/* clocks until queue goes critical */
		/* low priority read queue */
	unsigned int lpr_max_starve;
		/* high priority read queue */
	unsigned int hpr_max_starve;
		/* write queue */
	unsigned int w_max_starve;
		/* once queues are critical, how long should they run, in number of transactions */
	unsigned int lpr_xact_run_length;
	unsigned int hpr_xact_run_length;
	unsigned int w_xact_run_length;
	/* Store configuration */
		/* clocks to delay switching transaction store from preffered
			to alternate (depending on prefer_write or not ) */
	unsigned int rdwr_idle_gap;
		/* select which transaction stores is preffered.
			This store takes priority if it has transactions pending */
	unsigned int prefer_write;
};

int umctl2_init(struct umctl2_config_state *state);

#define MP_AXI_PORT_ENABLE(x) (1 << (x))

#endif
