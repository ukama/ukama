/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <pthread.h>

#include "config.h"
#include "starter.h"
#include "manifest.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_services.h"

SpaceList *gSpaceList = NULL;

void handle_sigint(int signum) {
    usys_log_debug("Terminate signal.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "logs",          required_argument, 0, 'l' },
    { "manifest-file", required_argument, 0, 'm' },
    { "help",          no_argument, 0, 'h' },
    { "version",       no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

void set_log_level(char *slevel) {

    int ilevel = USYS_LOG_TRACE;

    if (!strcmp(slevel, "TRACE")) {
        ilevel = USYS_LOG_TRACE;
    } else if (!strcmp(slevel, "DEBUG")) {
        ilevel = USYS_LOG_DEBUG;
    } else if (!strcmp(slevel, "INFO")) {
        ilevel = USYS_LOG_INFO;
    }
    usys_log_set_level(ilevel);
}

void usage() {

    usys_puts("Usage: starter.d [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-m, --manifest-file <file>    Manifest file");
    usys_puts("-v, --version                 Software version");
}

void fetch_and_update(void *config) {

    SpaceList *spacePtr = NULL;

    while (USYS_TRUE) {
        /* for each capp, with missing pkg, run a thred which fetch via wimc,
         * unpack into its space rootfs and run.
         */
        for (spacePtr = gSpaceList;
             spacePtr;
             spacePtr = spacePtr->next) {

            /* Always skip BOOT space */
            if (strcmp(spacePtr->space->name, SPACE_BOOT) == 0) {
                continue;
            }

            fetch_unpack_run(spacePtr->space, (Config *)config);
        }

        sleep(FETCH_AND_UPDATE_RETRY);
    }

    return NULL;
}

int main(int argc, char **argv) {

    int opt, optIdx;
    char *debug        = DEF_LOG_LEVEL;
    char *manifestFile = DEF_MANIFEST_FILE;
    UInst  serviceInst; 
    Config serviceConfig = {0};

    Manifest  *manifest  = NULL;
    SpaceList *spacePtr  = NULL;
    Space     *bootSpace = NULL;

    pthread_t thread;

    log_set_service(SERVICE_NAME);

    /* Parsing command line args. */
    while (true) {

        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "vh:m:p:l:n:d:w", longOptions,
                               &optIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            usage();
            usys_exit(0);
            break;

        case 'v':
            usys_puts(STARTER_VERSION);
            usys_exit(0);
            break;

        case 'l':
            debug = optarg;
            set_log_level(debug);
            break;

        case 'm':
            manifestFile = optarg;
            if (!manifestFile) {
                usage();
                usys_exit(0);
            }
            break;

        default:
            usage();
            usys_exit(0);
        }
    }

    /* Service config update */
    serviceConfig.servicePort  = usys_find_service_port(SERVICE_NAME);
    serviceConfig.nodedPort    = usys_find_service_port(SERVICE_NODE);
    serviceConfig.notifydPort  = usys_find_service_port(SERVICE_NOTIFY);
    serviceConfig.wimcPort     = usys_find_service_port(SERVICE_WIMC);
    serviceConfig.manifestFile = strdup(manifestFile);
    serviceConfig.nodeID       = NULL;

    if (!serviceConfig.servicePort ||
        !serviceConfig.nodedPort   ||
        !serviceConfig.notifydPort ||
        !serviceConfig.wimcPort) {
        usys_log_error("Unable to determine the port for services");
        usys_exit(1);
    }

    usys_log_debug("Starting %s ... ", SERVICE_NAME);

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /* Read and handle spaces/capps from the manifest file */
    if (!read_manifest_file(&manifest, serviceConfig.manifestFile)) {
        usys_log_error("Error with manifest file: %s",
                       serviceConfig.manifestFile);

        if (start_web_service(&serviceConfig,
                              &serviceInst) != USYS_TRUE) {
            usys_log_error("Webservice failed to setup for clients. Exiting.");
            exit(1);
        }

        pause();
        goto done;
    }

    if (validate_capp_dependency(&manifest) == USYS_FALSE) {

        usys_log_error("Invalid manifest file: %s",
                       serviceConfig.manifestFile);

        if (start_web_service(&serviceConfig,
                              &serviceInst) != USYS_TRUE) {
            usys_log_error("Webservice failed to setup for clients. Exiting.");
            exit(1);
        }

        pause();
        goto done;
    }

    process_manifest_file(&gSpaceList, manifest);
    print_spaces_list(gSpaceList);

    /* for each space: copy their capps into their rootfs at
     * /capps/rootfs/[space_name]/capps/pkg
     * paths are: DEF_CAPP_PATH and DEF_SPACE_ROOTFS_PATH
     */
    copy_capps_to_rootfs(gSpaceList);
    if (unpack_all_capps(gSpaceList) == USYS_FALSE) {
        usys_log_error("Unable to unpack the capps for cspace rootfs.");
        exit(1);
    }

    /* start all the apps - boot is reserved space and is
     * started first. Reboot is also reserved and only executed
     * when the system is rebooting
     */
    if (find_matching_space(&gSpaceList, SPACE_BOOT, &bootSpace)) {
        run_space_all_capps(bootSpace);
    }

    /* and everything else except 'boot' and 'reboot'*/
    for (spacePtr = gSpaceList;
         spacePtr;
         spacePtr = spacePtr->next) {

        /* BOOT space is already running */
        if (strcmp(spacePtr->space->name, SPACE_BOOT) == 0) {
            continue;
        }

        /* REBOOT is only run when system is restarting */
        if (strcmp(spacePtr->space->name, SPACE_REBOOT) == 0) {
            continue;
        }

        run_space_all_capps(spacePtr->space);
    }

    /* for each capp, with missing pkg, run a thread which fetch via wimc,
     *  unpack into its space rootfs and run.
     */
    for (spacePtr = gSpaceList;
         spacePtr;
         spacePtr = spacePtr->next) {

        /* BOOT space is already running */
        if (strcmp(spacePtr->space->name, SPACE_BOOT) == 0) {
            continue;
        }

        fetch_unpack_run(spacePtr->space, &serviceConfig);
    }

    pthread_create(&thread,
                   NULL,
                   fetch_and_update,
                   &serviceConfig);

    /* and finally, start the web service */
    if (start_web_service(&serviceConfig,
                          &serviceInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup for clients. Exiting.");
        exit(1);
    }

    pthread_join(thread, NULL);
    pause();

done:
    free_manifest(manifest);
    usys_log_debug("Exiting %s ...", SERVICE_NAME);

    return USYS_TRUE;
}

