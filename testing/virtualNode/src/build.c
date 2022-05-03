/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* build capps for virtual node */

#include <stdlib.h>
#include <stdio.h>

#include "config.h"

#define SCRIPT     "./scripts/mk_vnode_capps.sh"
#define MAX_BUFFER 1024

/*
 * build_capp --
 *
 */
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
	if (system(runMe) < 0) return FALSE;

	if (!build->staticFlag) {
		/* set rpath for the executable */
		sprintf(runMe, "%s patchelf %s", SCRIPT, build->binFrom);
		if (system(runMe) < 0 ) return FALSE;
	}

	sprintf(runMe, "%s mkdir %s_%s/%s", SCRIPT, capp->name,	capp->version,
			build->binTo);
	if (system(runMe) < 0) return FALSE;
	sprintf(runMe, "%s cp %s %s_%s/%s", SCRIPT, build->binFrom,
			capp->name, capp->version, build->binTo);
	if (system(runMe) < 0) return FALSE;

	sprintf(runMe, "%s mkdir %s_%s/%s", SCRIPT, capp->name, capp->version,
			build->mkdir);
	if (system(runMe) < 0) return FALSE;

	sprintf(runMe, "%s cp %s %s_%s/%s", SCRIPT, build->from, capp->name,
			capp->version, build->to);
	if (system(runMe) < 0) return FALSE;

	if (!build->staticFlag) {
	    sprintf(runMe, "%s libs %s %s_%s", SCRIPT, build->binFrom,
				capp->name, capp->version);
		if (system(runMe) < 0) return FALSE;
	}

	return TRUE;
}
