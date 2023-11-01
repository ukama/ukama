/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
