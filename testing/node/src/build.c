/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <stdlib.h>
#include <stdio.h>

#include "config.h"

#define SCRIPT     "./scripts/mk_vnode_capps.sh"
#define MAX_BUFFER 1024

int build_capp(Config *config) {

	char runMe[MAX_BUFFER] = {0};
	BuildConfig *build;
	CappConfig *capp;

	if (config == NULL)        return FALSE;
	if (config->build == NULL) return FALSE;
	if (config->capp == NULL)  return FALSE;

	build = config->build;
	capp  = config->capp;

	sprintf(runMe, "%s clean %s_%s", SCRIPT, capp->name, capp->version);
	if (system(runMe) < 0) return FALSE;

	sprintf(runMe, "%s mkdir %s_%s", SCRIPT, capp->name, capp->version);
	if (system(runMe) < 0) return FALSE;

	sprintf(runMe, "%s build app %s \"%s\"", SCRIPT, build->source, build->cmd);
	if (system(runMe) != 0) return FALSE;

	sprintf(runMe, "%s mkdir %s_%s/%s", SCRIPT, capp->name,	capp->version,
			build->binTo);
	if (system(runMe) < 0) return FALSE;

	sprintf(runMe, "%s cp %s %s_%s/%s", SCRIPT, build->binFrom,
			capp->name, capp->version, build->binTo);
	if (system(runMe) < 0) return FALSE;

	if (build->mkdir) {
		sprintf(runMe, "%s mkdir %s_%s/%s", SCRIPT, capp->name, capp->version,
				build->mkdir);
		if (system(runMe) < 0) return FALSE;
	}

	if (build->from) {
		sprintf(runMe, "%s cp-config %s %s", SCRIPT, build->from, build->to);
		if (system(runMe) < 0) return FALSE;
	}

	if (!build->staticFlag) {
	    sprintf(runMe, "%s libs %s %s_%s", SCRIPT, build->binFrom,
				capp->name, capp->version);
		if (system(runMe) < 0) return FALSE;
	}

    if (!build->staticFlag) {
		/* set rpath for the executable */
		sprintf(runMe, "%s patchelf %s", SCRIPT, build->binFrom);
		if (system(runMe) < 0 ) return FALSE;
	}

    return TRUE;
}
