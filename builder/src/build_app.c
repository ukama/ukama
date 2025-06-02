/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <stdio.h>

#include "config_app.h"

#define SCRIPT       "builder/scripts/make-app.sh"
#define LIB_USYS     "nodes/ukamaOS/distro/platform/build/libusys.so"
#define MAX_BUFFER   1024

int build_app(Config *config) {

    char *ukamaRoot = NULL;
    char runMe[MAX_BUFFER] = {0};
    BuildConfig *build;

  if (config == NULL)        return FALSE;
  if (config->build == NULL) return FALSE;
  if ((ukamaRoot = getenv("UKAMA_ROOT")) == NULL) return FALSE;

  build = config->build;

  sprintf(runMe, "%s/%s init %s_%s",
          ukamaRoot,
          SCRIPT,
          config->capp->name,
          config->capp->version);
  if (system(runMe) < 0) return FALSE;
  
  sprintf(runMe, "%s/%s build app %s \"%s\"",
          ukamaRoot, SCRIPT, build->source, build->cmd);
  if (system(runMe) < 0) return FALSE;

#if 0
#ifndef ALPINE_BUILD
  if (!build->staticFlag) {
      /* set rpath for the executable */
      sprintf(runMe, "%s/%s patchelf %s", ukamaRoot, SCRIPT, build->binFrom);
      if (system(runMe) < 0 ) return FALSE;
  }
#endif
#endif

  sprintf(runMe, "%s/%s cp %s %s_%s%s", ukamaRoot, SCRIPT,
          build->binFrom, config->capp->name, config->capp->version, build->binTo);
  if (system(runMe) < 0) return FALSE;

  if (build->mkdir) {
      sprintf(runMe, "%s/%s mkdir %s_%s%s", ukamaRoot, SCRIPT,
              config->capp->name, config->capp->version, build->mkdir);
      if (system(runMe) < 0) return FALSE;
  }

  if (build->from && build->to) {
      sprintf(runMe, "%s/%s cp %s %s_%s%s", ukamaRoot, SCRIPT,
              build->from, config->capp->name, config->capp->version, build->to);
      if (system(runMe) < 0) return FALSE;
  }

  if (build->miscFrom && build->miscTo) {
      sprintf(runMe, "%s/%s cp %s %s_%s%s", ukamaRoot, SCRIPT,
              build->miscFrom, config->capp->name,
              config->capp->version, build->miscTo);
      if (system(runMe) < 0) return FALSE;
  }

  if (!build->staticFlag) {
      sprintf(runMe, "%s/%s libs %s %s_%s", ukamaRoot, SCRIPT,
              build->binFrom, config->capp->name, config->capp->version);
      if (system(runMe) < 0) return FALSE;
  }

#if 0
  // Currently, we are focusing on using alpine - commented it out.
  if (!build->staticFlag) {
      sprintf(runMe, "%s/%s cp %s/%s %s_%s/lib",
              ukamaRoot, SCRIPT, ukamaRoot, LIB_USYS,
              config->capp->name, config->capp->version);
      if (system(runMe) < 0) return FALSE;
  }
#endif

  return TRUE;
}
