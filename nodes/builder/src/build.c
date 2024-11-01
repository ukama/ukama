/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <stdio.h>

#include "config.h"

#define SCRIPT "./scripts/mk-app.sh"
#define MAX_BUFFER   1024

int build_app(Config *config) {

  char runMe[MAX_BUFFER] = {0};
  BuildConfig *build;

  /* Flow is:
   * 1. build the app
   * 2. copy the binary to rootfs
   * 3. Make any additional directories within rootfs
   * 4. Copy any config file to rootfs
   * 5. Copy libs for non-static
   * 6. patch ELF to set rpath
   */
  if (config == NULL)        return FALSE;
  if (config->build == NULL) return FALSE;

  build = config->build;

  sprintf(runMe, "%s init %s_%s",
          SCRIPT,
          config->capp->name,
          config->capp->version);
  if (system(runMe) < 0) return FALSE;
  
  sprintf(runMe, "%s build app %s \"%s\"", SCRIPT, build->source, build->cmd);
  if (system(runMe) < 0) return FALSE;

  if (!build->staticFlag) {
      /* set rpath for the executable */
      sprintf(runMe, "%s patchelf %s", SCRIPT, build->binFrom);
      if (system(runMe) < 0 ) return FALSE;
  }

  sprintf(runMe, "%s cp %s %s_%s%s", SCRIPT, build->binFrom,
          config->capp->name, config->capp->version, build->binTo);
  if (system(runMe) < 0) return FALSE;

  if (build->mkdir) {
      sprintf(runMe, "%s mkdir %s_%s%s", SCRIPT,
              config->capp->name, config->capp->version, build->mkdir);
      if (system(runMe) < 0) return FALSE;
  }

  if (build->from && build->to) {
      sprintf(runMe, "%s cp %s %s_%s%s", SCRIPT, build->from,
              config->capp->name, config->capp->version, build->to);
      if (system(runMe) < 0) return FALSE;
  }

  if (build->miscFrom && build->miscTo) {
      sprintf(runMe, "%s cp %s %s_%s%s", SCRIPT, build->miscFrom,
              config->capp->name, config->capp->version, build->miscTo);
      if (system(runMe) < 0) return FALSE;
  }

  if (!build->staticFlag) {
      sprintf(runMe, "%s libs %s", SCRIPT, build->binFrom);
      if (system(runMe) < 0) return FALSE;
  }
  return TRUE;
}
