/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <service.h>
#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_types.h"

#define DEF_LOG_LEVEL  "TRACE"
#define DEF_PORT       "8080"
#define NOTIFY_VERSION "0.0.0"


/**
 * @fn      void noded_service()
 * @brief   Start Noded web service (REST server)
 *
 */
void notify_service() {
    service();
}

/**
 * @fn      int noded_startup(char*, char*, char*)
 * @brief   Do service initialization. Parse the required configs and
 *          initialize web frameworks.
 *
 * @param   invDb
 * @param   pCfg
 * @param   notifServer
 * @return  On success 0,
 *          On failure -1
 */
int notify_startup(int port) {
    int ret = 0;
    ret = service_init(port);
    return ret;
}

/**
 * @fn      void noded_exit()
 * @brief   Service exit procedure. Release the data structure used.
 */
void notify_exit() {
    service_at_exit();
}

/**
 * @fn      void handle_sigint(int)
 * @brief   Handle terminate signal for Noded
 *
 * @param   signum
 */
void handle_sigint(int signum) {
    usys_log_debug("Caught terminate signal.\n");

    /* Exiting NodeD */
    notify_exit();

    usys_log_debug("Cleanup complete.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "port", required_argument, 0, 'p' },
    { "logs", required_argument, 0, 'l' },
    { "help", no_argument, 0, 'h' },
    { "version", no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

/**
 * @fn      void set_log_level(char*)
 * @brief   Set the verbosity level for logs.
 *
 * @param   slevel
 */
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


/**
 * @fn      void usage()
 * @brief   Usage options for the ukamaEDR
 *
 */
void usage() {
    usys_puts("Usage: noded [options] \n");
    usys_puts("Options:\n");
    usys_puts(
        "--h, --help                             Help menu.\n");
    usys_puts(
        "--l, --logs <TRACE> <DEBUG> <INFO>      Log level for the process.\n");
    usys_puts(
        "--p, --port <port>                      Port at which service will"
              "listen.\n");
    usys_puts(
        "--v, --version                          Software Version.\n");
}

/**
 * @fn      int main(int, char**)
 * @brief
 *
 * @param   argc
 * @param   argv
 * @return  Should stay in main function entire time.
 */
int main(int argc, char **argv) {
    int ret = USYS_OK;
    char *debug = DEF_LOG_LEVEL;
    char *cPort = DEF_PORT;
    int  port = 0;

    /* Parsing command line args. */
    while (true) {
        int opt = 0;
        int opdIdx = 0;

        opt = usys_getopt_long(argc, argv, "h:p:l:v:", longOptions, &opdIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            usage();
            usys_exit(0);
            break;

        case 'v':
            usys_puts(NOTIFY_VERSION);
            usys_exit(0);
            break;

        case 'p':
            cPort = optarg;
            if (!cPort) {
                usage();
                usys_exit(0);
            }
            break;

        case 'l':
            debug = optarg;
            set_log_level(debug);
            break;

        default:
            usage();
            usys_exit(0);
        }
    }

    port = usys_atoi(cPort);

    usys_log_debug(
        "Starting notify service for monitoring and reporting node events.");

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /*  pre-startup routine. */
    ret = notify_startup(port);
    if (!ret) {
        /* Starting Service.*/
        notify_service();

        while (1) {
            usys_sleep(30);
        }
    }

    /* Should never reach here */
    notify_exit();

    usys_log_debug("Exiting notify service.");
    return ret;
}
