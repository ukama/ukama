/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* to build capp */

#include <stdlib.h>
#include <stdio.h>

#include "config.h"

#define SCRIPT "./scripts/mk_capp_rootfs.sh"
#define MAX_BUFFER   1024

/*
 * build_capp --
 *
 */
int build_capp(Config *config) {

  char runMe[MAX_BUFFER] = {0};
  BuildConfig *build;

  /* Flow is:
   * 1. create rootfs and busybox
   * 2. build the app 
   * 3. copy the binary to rootfs
   * 4. Make any additional directories within rootfs
   * 5. Copy any config file to rootfs
   * 6. Any misc related command
   * 7. Copy libs for non-static
   * 8. patch ELF to set rpath
   */

  if (config == NULL) return FALSE;
  if (config->build == NULL) return FALSE;

  build = config->build;

  sprintf(runMe, "%s clean", SCRIPT);
  if (system(runMe) < 0) return FALSE;

  sprintf(runMe, "%s clean %s_%s", SCRIPT, config->capp->name,
	  config->capp->version);
  if (system(runMe) < 0) return FALSE;
  
  sprintf(runMe, "%s build busybox", SCRIPT);
  if (system(runMe) < 0) return FALSE;

  sprintf(runMe, "%s build app %s \"%s\"", SCRIPT, build->source, build->cmd);
  if (system(runMe) < 0) return FALSE;

  if (!build->staticFlag) {
    /* set rpath for the executable */
    sprintf(runMe, "%s patchelf %s", SCRIPT, build->binFrom);
    if (system(runMe) < 0 ) return FALSE;
  }

  sprintf(runMe, "%s cp %s %s", SCRIPT, build->binFrom, build->binTo);
  if (system(runMe) < 0) return FALSE;

  sprintf(runMe, "%s mkdir %s", SCRIPT, build->mkdir);
  if (system(runMe) < 0) return FALSE;

  sprintf(runMe, "%s cp %s %s", SCRIPT, build->from, build->to);
  if (system(runMe) < 0) return FALSE;

  if (build->miscFrom && build->miscTo) {
    sprintf(runMe, "%s cp %s %s", SCRIPT, build->miscFrom, build->miscTo);
    if (system(runMe) < 0) return FALSE;
  }

  if (!build->staticFlag) {
    sprintf(runMe, "%s libs %s", SCRIPT, build->binFrom);
    if (system(runMe) < 0) return FALSE;
  }
  return TRUE;
}
