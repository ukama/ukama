/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/*
 * capp config.json
 */

#ifndef CAPP_RUNTIME_H
#define CAPP_RUNTIME_H

typedef struct capp_runtime_t {

  int   sockets[2];
  pid_t pid;
  
} CAppRuntime;

int create_and_run_capps(void *apps);

#endif /* CAPP_RUNTIME_H */
