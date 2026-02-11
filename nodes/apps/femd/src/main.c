



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
#include <unistd.h>

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
#include "safety_monitor.h"

#include "version.h"

/* network.c */
int start_web_service(ServerConfig *serverConfig, UInst *serviceInst);

/* Graceful shutdown flag and handlers */
volatile sig_atomic_t g_running = 1;

static void handle_terminate(int signum) {
    usys_log_debug("Terminate signal");
    (void)signum;
    g_running = 0;
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

static void install_signal_handlers(void) {
    struct sigaction sa;

    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = handle_terminate;
    sigemptyset(&sa.sa_mask);

    sigaction(SIGINT,  &sa, NULL);
    sigaction(SIGTERM, &sa, NULL);

#ifdef SIGPIPE
    signal(SIGPIPE, SIG_IGN);
#endif
}

static int validate_fem_band_env(void) {

    const char *env = getenv(ENV_FEM_BAND);
    char band[16];
    size_t i = 0;
    size_t j = 0;

    if (!env || !*env) {
        usys_log_error("Band env not set: %s Supported values: B1, B41, B48",
                       ENV_FEM_BAND);
        return STATUS_NOK;
    }

    while (env[i] && j < sizeof(band) - 1) {
        if (env[i] != ' ' && env[i] != '\t' &&
            env[i] != '\n' && env[i] != '\r') {
            band[j++] = env[i];
        }
        i++;
    }
    band[j] = '\0';

    if (strcasecmp(band, "B1") != 0 &&
        strcasecmp(band, "B41") != 0 &&
        strcasecmp(band, "B48") != 0) {

        usys_log_error("Invalid %s='%s'. Supported values: B1, B41, B48",
                       ENV_FEM_BAND, env);
        return STATUS_NOK;
    }

    return STATUS_OK;
}

int main(int argc, char **argv) {

    int opt;
    int optIdx;
    int exitCode = 0;

    char *debugLevel = DEF_LOG_LEVEL;

    /* zero-init for safety */
    UInst serviceInst = {0};
    Config serviceConfig = {0};
    GpioController gpioController = {0};
    I2CController i2cController = {0};
    ServerConfig serverConfig = {0};
    SafetyMonitor safetyMonitor = {0};

    /* state tracking flags */
    bool serviceNameAllocated = false;
    bool nodeIdAllocated = false;
    bool nodeTypeAllocated = false;
    bool gpioInitialized = false;
    bool i2cInitialized = false;
    bool safetyInitialized = false;
    bool webServiceStarted = false;

    usys_log_set_service(SERVICE_NAME);

    if (validate_fem_band_env() != STATUS_OK) {
        exitCode = 1;
        goto done;
    }

    while (true) {
        opt = usys_getopt_long(argc, argv, "vh:l:", longOptions, &optIdx);
        if (opt == -1) break;

        switch (opt) {
        case 'h':
            usage();
            return 0;

        case 'v':
            usys_puts(VERSION);
            return 0;

        case 'l':
            debugLevel = optarg;
            set_log_level(debugLevel);
            break;

        default:
            usage();
            return 0;
        }
    }

    /* Service configuration */
    serviceConfig.serviceName = usys_strdup(SERVICE_NAME);
    serviceNameAllocated = (serviceConfig.serviceName != NULL);

    serviceConfig.servicePort = usys_find_service_port(SERVICE_NAME);
    serviceConfig.nodedPort   = usys_find_service_port(SERVICE_NODE);
    serviceConfig.notifydPort = usys_find_service_port(SERVICE_NOTIFY);

    if (!serviceNameAllocated ||
        !serviceConfig.servicePort ||
        !serviceConfig.nodedPort ||
        !serviceConfig.notifydPort) {

        usys_log_error("Unable to determine service configuration");
        exitCode = 1;
        goto done;
    }

    install_signal_handlers();

    if (getenv(ENV_FEMD_DEBUG_MODE)) {
        serviceConfig.nodeID   = usys_strdup(DEF_NODE_ID);
        serviceConfig.nodeType = usys_strdup(DEF_NODE_TYPE);

        nodeIdAllocated = (serviceConfig.nodeID != NULL);
        nodeTypeAllocated = (serviceConfig.nodeType != NULL);

        if (!nodeIdAllocated || !nodeTypeAllocated) {
            exitCode = 1;
            goto done;
        }
    } else {
        if (get_nodeid_and_type_from_noded(&serviceConfig) != STATUS_OK) {
            exitCode = 1;
            goto done;
        }

        nodeIdAllocated = (serviceConfig.nodeID != NULL);
        nodeTypeAllocated = (serviceConfig.nodeType != NULL);

        if (!serviceConfig.nodeType ||
            strcmp(serviceConfig.nodeType, "Amplifier") != 0) {

            usys_log_error("Fem.d only runs on amplifier node");
            exitCode = 1;
            goto done;
        }
    }

    if (gpio_controller_init(&gpioController, NULL) != STATUS_OK) {
        exitCode = 1;
        goto done;
    }
    gpioInitialized = true;

    if (i2c_controller_init(&i2cController) != STATUS_OK) {
        exitCode = 1;
        goto done;
    }
    i2cInitialized = true;

    if (safety_monitor_init(&safetyMonitor,
                            &gpioController,
                            &i2cController,
                            &serviceConfig) != STATUS_OK) {
        exitCode = 1;
        goto done;
    }
    safetyInitialized = true;

    if (safety_monitor_start(&safetyMonitor) != STATUS_OK) {
        exitCode = 1;
        goto done;
    }

    serverConfig.config = &serviceConfig;
    serverConfig.gpioController = &gpioController;
    serverConfig.i2cController = &i2cController;

    if (start_web_service(&serverConfig, &serviceInst) != USYS_TRUE) {
        exitCode = 1;
        goto done;
    }
    webServiceStarted = true;

    usys_log_info("FEM.d started successfully");

    while (g_running) {
        usleep(200000);
    }

done:
    /* orderly shutdown */
    if (webServiceStarted) {
        ulfius_stop_framework(&serviceInst);
        ulfius_clean_instance(&serviceInst);
    }

    if (safetyInitialized)     safety_monitor_cleanup(&safetyMonitor);
    if (i2cInitialized)        i2c_controller_cleanup(&i2cController);
    if (gpioInitialized)       gpio_controller_cleanup(&gpioController);
    if (serviceNameAllocated)  usys_free(serviceConfig.serviceName);
    if (nodeIdAllocated)       usys_free(serviceConfig.nodeID);
    if (nodeTypeAllocated)     usys_free(serviceConfig.nodeType);

    usys_log_debug("Exiting %s ...", SERVICE_NAME);

    return exitCode;
}
