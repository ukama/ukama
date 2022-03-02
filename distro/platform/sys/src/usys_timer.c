/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
        return false;
    }

    /* Setup a periodic tick of given resolution in Micro seconds with iTimers */
    intervalTimer.it_value.tv_sec = 0;
    intervalTimer.it_value.tv_usec = resolution;
    intervalTimer.it_interval.tv_sec = 0;
    intervalTimer.it_interval.tv_usec = resolution;
    if (setitimer(ITIMER_REAL, &intervalTimer, NULL)) {
        usys_log_error("Error in setting up the timer\n");
        return false;
    }
    return true;
}
