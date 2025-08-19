/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */
#include <pthread.h>
#include <signal.h>
#include <getopt.h>
#include <stdbool.h>

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_services.h"

#include "config.h"
#include "femd.h"
#include "web_service.h"
#include "web_client.h"
#include "gpio_controller.h"
#include "i2c_controller.h"

#include "version.h"

void handle_sigint(int signum) {
    usys_log_debug("Terminate signal.\n");
    usys_exit(0);
}

static UsysOption longOptions[] = {
    { "logs",    required_argument, 0, 'l' },
    { "config",  required_argument, 0, 'c' },
    { "help",    no_argument, 0, 'h' },
    { "version", no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

static void set_log_level(char *slevel) {
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

static void usage() {
    usys_puts("Usage: femd [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-c, --config FILE             Configuration file");
    usys_puts("-v, --version                 Software version");
}

int main(int argc, char **argv) {
    int opt, optIdx;
    int exitCode = USYS_FALSE;

    char *debug = DEF_LOG_LEVEL;
    UInst serviceInst;
    Config serviceConfig = {0};

    GpioController gpioController = {0};
    I2CController  i2cController  = {0};

    usys_log_set_service(SERVICE_NAME);
    usys_log_remote_init(SERVICE_NAME);

    while (true) {
        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "vh:l:", longOptions, &optIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            usage();
            usys_exit(0);
            break;

        case 'v':
            usys_puts(VERSION);
            usys_exit(0);
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

    /* Service config update */
    serviceConfig.serviceName = usys_strdup(SERVICE_NAME);
    serviceConfig.servicePort = usys_find_service_port(SERVICE_NAME);
    serviceConfig.nodedPort   = usys_find_service_port(SERVICE_NODE);
    serviceConfig.notifydPort = usys_find_service_port(SERVICE_NOTIFY);
    serviceConfig.nodeID      = NULL;
    serviceConfig.nodeType    = NULL;

    if (!serviceConfig.servicePort ||
        !serviceConfig.nodedPort   ||
        !serviceConfig.notifydPort) {
        usys_log_error("Unable to determine the port for service(s)");
        usys_exit(1);
    }

    usys_log_debug("Starting %s ... ", SERVICE_NAME);

    /* Signal handler */
    signal(SIGINT, handle_sigint);

    if (getenv(ENV_FEMD_DEBUG_MODE)) {
        serviceConfig.nodeID   = strdup(DEF_NODE_ID);
        serviceConfig.nodeType = strdup(DEF_NODE_TYPE);
        usys_log_debug("%s: using default Node ID: %s Type: %s",
                       SERVICE_NAME,
                       DEF_NODE_ID,
                       DEF_NODE_TYPE);
    } else {
        if (get_nodeid_and_type_from_noded(&serviceConfig) == STATUS_NOK) {
            usys_log_error(
                "%s: unable to connect with node.d", SERVICE_NAME);
            goto done;
        }
    }

    if (gpio_controller_init(&gpioController, NULL) != STATUS_OK) {
        usys_log_error("Failed to initialize GPIO controller");
        exitCode = USYS_TRUE;
        goto done;
    }

    if (i2c_controller_init(&i2cController) != STATUS_OK) {
        usys_log_error("Failed to initialize I2C controller");
        exitCode = USYS_TRUE;
        goto done;
    }

    if (start_web_service(&serviceConfig, &serviceInst, NULL) != USYS_TRUE) {
        usys_free(serviceConfig.serviceName);
        usys_log_error("Webservice failed to setup for clients. Exiting.");
        exitCode = USYS_TRUE;
        usys_exit(1);
    }

    pause();

done:
    ulfius_stop_framework(&serviceInst);
    ulfius_clean_instance(&serviceInst);
    usys_free(serviceConfig.serviceName);

    i2c_controller_cleanup(&i2cController);
    gpio_controller_cleanup(&gpioController);
    
    usys_log_debug("Exiting femd ...");

    return exitCode;
}
