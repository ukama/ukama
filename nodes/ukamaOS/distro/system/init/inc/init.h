/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * init.h
 */

#ifndef UKAMA_INIT_H
#define UKAMA_INIT_H

#include <signal.h>

#define UKAMA_BANNER "UkamaOS V00.00.00"

#define TRUE  1
#define FALSE 0

static void setup_console(void);
static void setup_term(void);
static pid_t run_task(char *taskExe, int wait);
static void process_signals(sigset_t *sigset, struct timespec *tspec);

#endif /* UKAMA_INIT_H */
