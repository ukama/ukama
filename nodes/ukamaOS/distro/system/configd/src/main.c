/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "config.h"
#include "config_macros.h"
#include "service.h"
#include "web.h"
#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"

/**
 * @fn      void handle_sigint(int)
 * @brief   Handle terminate signal for Noded
 *
 * @param   signum
 */
void handle_sigint(int signum) {
    usys_log_debug("Caught terminate signal.\n");
    usys_log_debug("Cleanup complete.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "port",          required_argument, 0, 'p' },
    { "logs",          required_argument, 0, 'l' },
    { "noded-host",    required_argument, 0, 'n' },
    { "noded-port",    required_argument, 0, 's' },
    { "noded-lep",     required_argument, 0, 'e' },
    { "help",          no_argument,       0, 'h' },
    { "version",       no_argument,       0, 'v' },
    { 0,               0,                 0,  0 }
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
        "-h, --help                             Help menu.\n");
    usys_puts(
        "-l, --logs <TRACE> <DEBUG> <INFO>      Log level for the process.\n");
    usys_puts(
        "--p, --port <port>                      Port at which service will"
              "listen.\n");
    usys_puts(
        "-n, --noded-host <host>               Host at which noded service"
                  "will listen.\n");
    usys_puts(
        "-s, --noded-port <port>               Port at which noded service"
                   "will listen.\n");
    usys_puts(
        "-e, --noded-ep </node>                API EP at which noded service"
                       "will enquire for node info.\n");
    usys_puts(
        "-f, --map-file <file-name>         Status map file\n");

    usys_puts(
        "-v, --version                          Software Version.\n");
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
    int ret = USYS_OK, port=0;

    char *debug        = DEF_LOG_LEVEL;
    char *cPort        = DEF_SERVICE_PORT;
    char *nodedHost    = DEF_NODED_HOST;
    char *nodedPort    = DEF_NODED_PORT;
    char *nodedEP      = DEF_NODED_EP;

    UInst serviceInst;

    Config serviceConfig = {0};

    /* Parsing command line args. */
    while (true) {
        int opt = 0;
        int opdIdx = 0;

        opt = getopt_long(argc, argv, "f:p:l:n:hv", longOptions, &opdIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            usage();
            usys_exit(0);
            break;

        case 'v':
            usys_puts(CONFIG_VERSION);
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
        case 'n':
            nodedHost = optarg;
            if (!nodedHost) {
                usage();
                usys_exit(0);
            }
            break;
        case 's':
            nodedPort = optarg;
            if (!nodedPort) {
                usage();
                usys_exit(0);
            }
            break;
        case 'e':
            nodedEP = optarg;
            if (!nodedEP) {
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
    serviceConfig.serviceName  = usys_strdup(SERVICE_NAME);
    serviceConfig.servicePort  = usys_atoi(cPort);
    serviceConfig.nodedEP      = usys_strdup(nodedEP);
    serviceConfig.nodedHost      = usys_strdup(nodedHost);
    serviceConfig.nodedPort     = usys_atoi(nodedPort);

    usys_log_debug("Starting config.d ...");

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    /* Read Node Info from noded */
    if (getenv(ENV_CONFIG_DEBUG_MODE)) {
       serviceConfig.nodeId = usys_strdup(DEF_NODE_ID);
       usys_log_debug("config.d: Using default Node ID: %s", DEF_NODE_ID);
    } else {
        if (get_nodeid_from_noded(&serviceConfig) == STATUS_NOK) {
            usys_log_error("config.d: Unable to connect with node.d");
            goto done;
        }
    }

    if (configd_read_running_config(&serviceConfig.runningConfig)) {
    	usys_log_error("Failed to read last running config.");
    	exit(1);
    }

    if (start_web_services(&serviceConfig, &serviceInst) != USYS_TRUE) {
        usys_log_error("Webservice failed to setup for clients. Exiting.");
        exit(1);
    }

    pause();

done:
    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);

    free(serviceConfig.serviceName);
    free(serviceConfig.nodeId);
    free(serviceConfig.nodedEP);
    free(serviceConfig.nodedHost);
    free_config_data(serviceConfig.runningConfig);
    usys_log_debug("Exiting config.d ...");
    return 1;
}
