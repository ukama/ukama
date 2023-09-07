/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2008-2009 coresystems GmbH
 * Copyright (C) 2014 Google Inc.
 * Copyright (C) 2016-2017 Intel Corporation.
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

#ifndef _SKYLAKE_ME_H_
#define _SKYLAKE_ME_H_

/*
 * Management Engine PCI registers
 */
#define  ME_HFS_CWS_RESET	0
#define  ME_HFS_CWS_INIT	1
#define  ME_HFS_CWS_REC		2
#define  ME_HFS_CWS_NORMAL	5
#define  ME_HFS_CWS_WAIT	6
#define  ME_HFS_CWS_TRANS	7
#define  ME_HFS_CWS_INVALID	8
#define  ME_HFS_STATE_PREBOOT	0
#define  ME_HFS_STATE_M0_UMA	1
#define  ME_HFS_STATE_M3	4
#define  ME_HFS_STATE_M0	5
#define  ME_HFS_STATE_BRINGUP	6
#define  ME_HFS_STATE_ERROR	7
#define  ME_HFS_ERROR_NONE	0
#define  ME_HFS_ERROR_UNCAT	1
#define  ME_HFS_ERROR_IMAGE	3
#define  ME_HFS_ERROR_DEBUG	4
#define  ME_HFS_MODE_NORMAL	0
#define  ME_HFS_MODE_DEBUG	2
#define  ME_HFS_MODE_DIS	3
#define  ME_HFS_MODE_OVER_JMPR	4
#define  ME_HFS_MODE_OVER_MEI	5
#define  ME_HFS_BIOS_DRAM_ACK	1
#define  ME_HFS_POWER_SOURCE_AC 1
#define  ME_HFS_POWER_SOURCE_DC 2

/* Infrastructure Progress Values */
#define  ME_HFS2_PHASE_ROM		0
#define  ME_HFS2_PHASE_UKERNEL		2
#define  ME_HFS2_PHASE_BUP		3
#define  ME_HFS2_PHASE_HOST_COMM	6
/* Current State - Based on Infra Progress values. */
/* ROM State */
#define  ME_HFS2_STATE_ROM_BEGIN 0
#define  ME_HFS2_STATE_ROM_DISABLE 6
/* BUP State */
#define  ME_HFS2_STATE_BUP_INIT 0
#define  ME_HFS2_STATE_BUP_DIS_HOST_WAKE 1
#define  ME_HFS2_STATE_BUP_CG_ENABLE 2
#define  ME_HFS2_STATE_BUP_PM_HND_EN 3
#define  ME_HFS2_STATE_BUP_FLOW_DET 4
#define  ME_HFS2_STATE_BUP_PMC_PATCHING 5
#define  ME_HFS2_STATE_BUP_GET_FLASH_VSCC 6
#define  ME_HFS2_STATE_BUP_SET_FLASH_VSCC 7
#define  ME_HFS2_STATE_BUP_VSCC_ERR 8
#define  ME_HFS2_STATE_BUP_EFSS_INIT 9
#define  ME_HFS2_STATE_BUP_CHECK_STRAP 0xa
#define  ME_HFS2_STATE_BUP_PWR_OK_TIMEOUT 0xb
#define  ME_HFS2_STATE_BUP_STRAP_DIS 0xc
#define  ME_HFS2_STATE_BUP_MANUF_OVRD_STRAP 0xd
#define  ME_HFS2_STATE_BUP_M3 0x11
#define  ME_HFS2_STATE_BUP_M0 0x12
#define  ME_HFS2_STATE_BUP_FLOW_DET_ERR 0x13
#define  ME_HFS2_STATE_BUP_M3_CLK_ERR 0x15
#define  ME_HFS2_STATE_BUP_CPU_RESET_DID_TIMEOUT_MEM_MISSING 0x17
#define  ME_HFS2_STATE_BUP_M3_KERN_LOAD 0x18
#define  ME_HFS2_STATE_BUP_T32_MISSING 0x1c
#define  ME_HFS2_STATE_BUP_WAIT_DID 0x1f
#define  ME_HFS2_STATE_BUP_WAIT_DID_FAIL 0x20
#define  ME_HFS2_STATE_BUP_DID_NO_FAIL 0x21
#define  ME_HFS2_STATE_BUP_ENABLE_UMA 0x22
#define  ME_HFS2_STATE_BUP_ENABLE_UMA_ERR 0x23
#define  ME_HFS2_STATE_BUP_SEND_DID_ACK 0x24
#define  ME_HFS2_STATE_BUP_SEND_DID_ACK_ERR 0x25
#define  ME_HFS2_STATE_BUP_M0_CLK 0x26
#define  ME_HFS2_STATE_BUP_M0_CLK_ERR 0x27
#define  ME_HFS2_STATE_BUP_TEMP_DIS 0x28
#define  ME_HFS2_STATE_BUP_M0_KERN_LOAD 0x32
/* Policy Module State */
#define  ME_HFS2_STATE_POLICY_ENTRY 0
#define  ME_HFS2_STATE_POLICY_RCVD_S3 3
#define  ME_HFS2_STATE_POLICY_RCVD_S4 4
#define  ME_HFS2_STATE_POLICY_RCVD_S5 5
#define  ME_HFS2_STATE_POLICY_RCVD_UPD 6
#define  ME_HFS2_STATE_POLICY_RCVD_PCR 7
#define  ME_HFS2_STATE_POLICY_RCVD_NPCR 8
#define  ME_HFS2_STATE_POLICY_RCVD_HOST_WAKE 9
#define  ME_HFS2_STATE_POLICY_RCVD_AC_DC 0xa
#define  ME_HFS2_STATE_POLICY_RCVD_DID 0xb
#define  ME_HFS2_STATE_POLICY_VSCC_NOT_FOUND 0xc
#define  ME_HFS2_STATE_POLICY_VSCC_INVALID 0xd
#define  ME_HFS2_STATE_POLICY_FPB_ERR 0xe
#define  ME_HFS2_STATE_POLICY_DESCRIPTOR_ERR 0xf
#define  ME_HFS2_STATE_POLICY_VSCC_NO_MATCH 0x10
/* Current PM Event Values */
#define  ME_HFS2_PMEVENT_CLEAN_MOFF_MX_WAKE 0
#define  ME_HFS2_PMEVENT_MOFF_MX_WAKE_ERROR 1
#define  ME_HFS2_PMEVENT_CLEAN_GLOBAL_RESET 2
#define  ME_HFS2_PMEVENT_CLEAN_GLOBAL_RESET_ERROR 3
#define  ME_HFS2_PMEVENT_CLEAN_ME_RESET 4
#define  ME_HFS2_PMEVENT_ME_RESET_EXCEPTION 5
#define  ME_HFS2_PMEVENT_PSEUDO_ME_RESET 6
#define  ME_HFS2_PMEVENT_CM0_CM3 7
#define  ME_HFS2_PMEVENT_CM3_CM0 8
#define  ME_HFS2_PMEVENT_NON_PWR_CYCLE_RESET 9
#define  ME_HFS2_PMEVENT_PWR_CYCLE_RESET_M3 0xa
#define  ME_HFS2_PMEVENT_PWR_CYCLE_RESET_MOFF 0xb
#define  ME_HFS2_PMEVENT_CMX_CMOFF 0xc
#define  ME_HFS2_PMEVENT_CM0_CM0PG 0xd
#define  ME_HFS2_PMEVENT_CM3_CM3PG 0xe
#define  ME_HFS2_PMEVENT_CM0PG_CM0 0xf

union me_hfs2 {
	u32 data;
	struct {
	u32 reserved1: 3;
	u32 invoke_mebx: 1;
	u32 cpu_replaced_sts: 1;
	u32 reserved2: 1;
	u32 mfs_failure: 1;
	u32 warm_reset_request: 1;
	u32 cpu_replaced_valid: 1;
	u32 low_power_state: 1;
	u32 power_gating_ind: 1;
	u32 reserved3: 1;
	u32 fw_upd_forced_sb: 1;
	u32 reserved4: 3;
	u32 current_state: 8;
	u32 current_pmevent: 4;
	u32 progress_code: 4;
	} __packed fields;
};

#define  ME_HFS3_FW_SKU_CONSUMER	0x2
#define  ME_HFS3_FW_SKU_CORPORATE	0x3

union me_hfs3 {
	u32 data;
	struct {
	u32 reserved1: 4;
	u32 fw_sku: 3;
	u32 encrypt_key_check: 1;
	u32 pch_config_change: 1;
	u32 reserved2: 21;
	u32 encrypt_key_override: 1;
	u32 power_down_mitigation: 1;
	} __packed fields;
};

#define ME_HFS6_FPF_NOT_COMMITTED	0x0
#define ME_HFS6_FPF_ERROR		0x2

union me_hfs6 {
	u32 data;
	struct {
		u32 reserved1: 30;
		u32 fpf_nvars: 2;
	} __packed fields;
};

#define MKHI_GEN_GROUP_ID	0xff

#define MKHI_GET_FW_VERSION	0x02

void intel_me_status(void);
int send_global_reset(void);

#endif
