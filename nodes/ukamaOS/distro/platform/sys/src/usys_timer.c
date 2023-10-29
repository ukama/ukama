/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include "usys_timer.h"
#include "usys_log.h"

bool usys_timer(uint32_t resolution, void (*tick_handler)()) {
    sigset_t thMask;
    struct sigaction sa;
    struct itimerval intervalTimer;

    /* Create a tick handler */
    sigemptyset(&thMask);
    sigaddset(&thMask, SIGALRM);
    pthread_sigmask(SIG_UNBLOCK, &thMask, NULL);

    sigfillset(&sa.sa_mask);
    sa.sa_handler = tick_handler;
    sa.sa_flags = SA_SIGINFO;
    if (sigaction(SIGALRM, &sa, NULL) == -1) {
        usys_log_error("Error in Sigaction\n");
        return USYS_FALSE;
    }

    /* Setup a periodic tick of given resolution in Micro seconds with iTimers */
    intervalTimer.it_value.tv_sec = 0;
    intervalTimer.it_value.tv_usec = resolution;
    intervalTimer.it_interval.tv_sec = 0;
    intervalTimer.it_interval.tv_usec = resolution;
    if (setitimer(ITIMER_REAL, &intervalTimer, NULL)) {
        usys_log_error("Error in setting up the timer\n");
        return USYS_FALSE;
    }
    return USYS_TRUE;
}
