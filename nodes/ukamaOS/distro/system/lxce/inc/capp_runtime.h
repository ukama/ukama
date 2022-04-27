/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
