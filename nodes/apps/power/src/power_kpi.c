/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <stdio.h>

#include "power_kpi.h"
#include "usys_log.h"

static double clamp0(double v) {
	return (v < 0) ? 0 : v;
}

static PowerSeverity worst(PowerSeverity a, PowerSeverity b) {
	return (a > b) ? a : b;
}

/*
 * For now, keep thresholds simple and stable.
 * If you later want config-driven thresholds, we can add it without changing JSON schema.
 */
static const double TEMP_WARN_C = 80.0;
static const double TEMP_CRIT_C = 90.0;

void power_eval(PowerSnapshot *s) {

	PowerSeverity sev = POWER_SEV_OK;
	char why[128];

	if (!s) return;

	/* Rail power */
	s->rail_in.w  = clamp0(s->rail_in.v  * s->rail_in.i);
	s->rail_aux.w = clamp0(s->rail_aux.v * s->rail_aux.i);

	/* Total power */
	s->total_w = s->rail_in.w + s->rail_aux.w;

	/* Start with rail severities */
	sev = worst(sev, s->rail_in.severity);
	sev = worst(sev, s->rail_aux.severity);

	/* Temperature severity */
	why[0] = '\0';
	if (s->temp_board_c > 0) {
		if (s->temp_board_c >= TEMP_CRIT_C) {
			sev = worst(sev, POWER_SEV_CRIT);
			snprintf(why, sizeof(why), "board temp critical (%.1fC)", s->temp_board_c);
		} else if (s->temp_board_c >= TEMP_WARN_C) {
			sev = worst(sev, POWER_SEV_WARN);
			snprintf(why, sizeof(why), "board temp high (%.1fC)", s->temp_board_c);
		}
	}

	/* Final overall */
	s->overall_severity = sev;

	if (why[0]) {
		snprintf(s->overall_reason, sizeof(s->overall_reason), "%s", why);
	} else {
		switch (sev) {
		case POWER_SEV_OK:
			snprintf(s->overall_reason, sizeof(s->overall_reason), "ok");
			break;
		case POWER_SEV_WARN:
			snprintf(s->overall_reason, sizeof(s->overall_reason), "warn");
			break;
		case POWER_SEV_CRIT:
			snprintf(s->overall_reason, sizeof(s->overall_reason), "crit");
			break;
		default:
			snprintf(s->overall_reason, sizeof(s->overall_reason), "ok");
			break;
		}
	}
}
