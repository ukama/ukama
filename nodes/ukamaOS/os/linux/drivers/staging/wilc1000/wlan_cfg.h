/* SPDX-License-Identifier: GPL-2.0 */
/*
 * Copyright (c) 2012 - 2018 Microchip Technology Inc., and its subsidiaries.
 * All rights reserved.
 */

#ifndef WILC_WLAN_CFG_H
#define WILC_WLAN_CFG_H

struct wilc_cfg_byte {
	u16 id;
	u8 val;
};

struct wilc_cfg_hword {
	u16 id;
	u16 val;
};

struct wilc_cfg_word {
	u32 id;
	u32 val;
};

struct wilc_cfg_str {
	u16 id;
	u8 *str;
};

struct wilc_cfg_bin {
	u16 id;
	u8 *bin;
};

struct wilc_cfg_str_vals {
	u8 mac_address[7];
	u8 firmware_version[129];
	u8 assoc_rsp[256];
};

struct wilc_bin_vals {
	u8 antenna_param[5];
};

struct wilc_cfg {
	struct wilc_cfg_byte *b;
	struct wilc_cfg_hword *hw;
	struct wilc_cfg_word *w;
	struct wilc_cfg_str *s;
	struct wilc_cfg_str_vals *str_vals;
	struct wilc_cfg_bin *bin;
	struct wilc_bin_vals *bin_vals;
};

struct wilc;
int cfg_set_wid(struct wilc_vif *vif, u8 *frame, u32 offset, u16 id, u8 *buf,
		int size);
int cfg_get_wid(u8 *frame, u32 offset, u16 id);
int cfg_get_val(struct wilc *wl, u16 wid, u8 *buffer, u32 buffer_size);
void cfg_indicate_rx(struct wilc *wilc, u8 *frame, int size,
			       struct wilc_cfg_rsp *rsp);
int cfg_init(struct wilc *wl);
void cfg_deinit(struct wilc *wl);

#endif
