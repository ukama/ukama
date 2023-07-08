/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <unistd.h>
#include <pthread.h>
#include <sys/reboot.h>

#include "config.h"
#include "deviced.h"

void* reboot_node(void *args) {

    Config *config = (Config *)args;

    if (strcmp(config->nodeType, UKAMA_TOWER_NODE) == 0) {

        /* send alarm to notify.d */
        if (wc_send_alarm_to_notifyd(config) == USYS_NOK) {
            usys_log_error("Unable to send notification to notify.d");
            usys_log_error("Reboot not processed");
            return;
        }

        usys_log_debug("Reboot alarm sent to notfy.d %ld", time(NULL));

        /* wait for few seconds */
        sleep(WAIT_BEFORE_REBOOT);
       
        /* send command to client-mode device.d to restart */
        /* xxx - to implement */

        /* reboot */
        if (getenv(ENV_DEVICED_DEBUG_MODE) == NULL) {
            sync();
            setuid(0);
            reboot(RB_AUTOBOOT);
        }
    }
}

void process_reboot(Config *config) {

    pthread_t thread;
    
    if (config->clientMode && getenv(ENV_DEVICED_DEBUG_MODE) == NULL) {
        sync();
        setuid(0);
        reboot(RB_AUTOBOOT);
        return;
    }

    pthread_create(&thread, NULL, reboot_node, (void *)config);
}
