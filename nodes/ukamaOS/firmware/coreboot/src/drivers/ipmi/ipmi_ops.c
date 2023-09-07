/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2019 Wiwynn Corp.
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

#include <console/console.h>
#include "ipmi_ops.h"
#include <string.h>

enum cb_err ipmi_init_and_start_bmc_wdt(const int port, uint16_t countdown,
				uint8_t action)
{
	int ret;
	struct ipmi_wdt_req req = {0};
	struct ipmi_rsp rsp;
	printk(BIOS_INFO, "Initializing IPMI BMC watchdog timer\n");
	/* BIOS FRB2 */
	req.timer_use = 1;
	req.timer_actions = action;
	/* clear BIOS FRB2 expiration flag */
	req.timer_use_expiration_flags_clr = 2;
	req.initial_countdown_val = countdown;
	ret = ipmi_kcs_message(port, IPMI_NETFN_APPLICATION, 0x0,
			IPMI_BMC_SET_WDG_TIMER,
			(const unsigned char *) &req, sizeof(req),
			(unsigned char *) &rsp, sizeof(rsp));

	if (ret < sizeof(struct ipmi_rsp) || rsp.completion_code) {
		printk(BIOS_ERR, "IPMI: %s set wdt command failed "
			"(ret=%d resp=0x%x), failed to initialize and start "
			"IPMI BMC watchdog timer\n", __func__,
			ret, rsp.completion_code);
		return CB_ERR;
	}

	/* Reset command to start timer */
	ret = ipmi_kcs_message(port, IPMI_NETFN_APPLICATION, 0x0,
			IPMI_BMC_RESET_WDG_TIMER, NULL, 0,
			(unsigned char *) &rsp, sizeof(rsp));

	if (ret < sizeof(struct ipmi_rsp) || rsp.completion_code) {
		printk(BIOS_ERR, "IPMI: %s reset wdt command failed "
			"(ret=%d resp=0x%x), failed to initialize and start "
			"IPMI BMC watchdog timer\n", __func__,
			ret, rsp.completion_code);
		return CB_ERR;
	}

	printk(BIOS_INFO, "IPMI BMC watchdog initialized and started.\n");
	return CB_SUCCESS;
}

enum cb_err ipmi_stop_bmc_wdt(const int port)
{
	int ret;
	struct ipmi_wdt_req req;
	struct ipmi_wdt_rsp rsp = {0};
	struct ipmi_rsp resp;

	/* Get current timer first */
	ret = ipmi_kcs_message(port, IPMI_NETFN_APPLICATION, 0x0,
			IPMI_BMC_GET_WDG_TIMER, NULL, 0,
			(unsigned char *) &rsp, sizeof(rsp));

	if (ret < sizeof(struct ipmi_rsp) || rsp.resp.completion_code) {
		printk(BIOS_ERR, "IPMI: %s get wdt command failed "
			"(ret=%d resp=0x%x), IPMI BMC watchdog timer may still "
			"be running\n", __func__, ret,
			 rsp.resp.completion_code);
		return CB_ERR;
	}
	/* If bit 6 in timer_use is 0 then it's already stopped. */
	if (!(rsp.data.timer_use & (1 << 6))) {
		printk(BIOS_DEBUG, "IPMI BMC watchdog is already stopped\n");
		return CB_SUCCESS;
	}
	/* Set timer stop running by clearing bit 6. */
	rsp.data.timer_use &= ~(1 << 6);
	rsp.data.initial_countdown_val = 0;
	req = rsp.data;
	ret = ipmi_kcs_message(port, IPMI_NETFN_APPLICATION, 0x0,
			IPMI_BMC_SET_WDG_TIMER,
			(const unsigned char *) &req, sizeof(req),
			(unsigned char *) &resp, sizeof(resp));

	if (ret < sizeof(struct ipmi_rsp) || resp.completion_code) {
		printk(BIOS_ERR, "IPMI: %s set wdt command stop timer failed "
			"(ret=%d resp=0x%x), failed to stop IPMI "
			"BMC watchdog timer\n", __func__, ret,
			resp.completion_code);
		return CB_ERR;
	}
	printk(BIOS_DEBUG, "IPMI BMC watchdog is stopped\n");

	return CB_SUCCESS;
}

enum cb_err ipmi_get_system_guid(const int port, uint8_t *uuid)
{
	int ret;
	struct ipmi_get_system_guid_rsp rsp;

	if (uuid == NULL) {
		printk(BIOS_ERR, "%s failed, null pointer parameter\n",
			 __func__);
		return CB_ERR;
	}

	ret = ipmi_kcs_message(port, IPMI_NETFN_APPLICATION, 0x0,
			IPMI_BMC_GET_SYSTEM_GUID, NULL, 0,
			(unsigned char *) &rsp, sizeof(rsp));

	if (ret < sizeof(struct ipmi_rsp) || rsp.resp.completion_code) {
		printk(BIOS_ERR, "IPMI: %s command failed (ret=%d resp=0x%x)\n",
			 __func__, ret, rsp.resp.completion_code);
		return CB_ERR;
	}

	memcpy(uuid, rsp.data, 16);
	return CB_SUCCESS;
}
