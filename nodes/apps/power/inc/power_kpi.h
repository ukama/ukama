/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
#ifndef __POWER_KPI_H__
#define __POWER_KPI_H__

#include <stdint.h>

#include "metrics_store.h"

/*
 * Simple calibration model.
 * Keep env minimal: only POWER_CFG (JSON) overrides these defaults.
 */
typedef struct {
	/* ADS1015 channels for rails (if applicable) */
	int		ch_12v;
	int		ch_5v;
	int		ch_3v3;
	int		ch_28v;

	/* Voltage dividers: Vrail = Vadcin * v_gain + v_offset */
	double		v_gain_12v;
	double		v_gain_5v;
	double		v_gain_3v3;
	double		v_gain_28v;

	double		v_off_12v;
	double		v_off_5v;
	double		v_off_3v3;
	double		v_off_28v;

	/*
	 * Current sense model (LTC6102-style):
	 * I = Vout * (RIN/ROUT) / RSENSE
	 */
	double		rsense_12v;
	double		rin_12v;
	double		rout_12v;

	double		rsense_28v;
	double		rin_28v;
	double		rout_28v;

	/* Limits */
	double		lim_12v_min, lim_12v_max;
	double		lim_5v_min, lim_5v_max;
	double		lim_3v3_min, lim_3v3_max;
	double		lim_28v_min, lim_28v_max;

	double		lim_i_12v_max;
	double		lim_i_28v_max;

	double		lim_temp_board_max;
	double		lim_temp_supply_max;
} PowerCal;

void power_cal_defaults(PowerCal *c);
int power_cal_from_env(PowerCal *c);	/* reads POWER_CFG if present */
void power_eval(PowerSnapshot *s);	    /* sets rail w/severity/overall */

#endif /* __POWER_KPI_H__ */
