/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "builder.h"
#include "config.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_services.h"

/* build.c */
extern bool build_all_systems(char *systemsList, char *ukamaRepo, char *authRepo);
extern bool build_nodes(int count, char *list, char *repo);

/* deploy.c */
extern bool deploy_all_systems(char *file, DeployConfig *deployConfig, char *ukamaRepo, char *authRepo);
extern bool display_all_systems_status(char *systems, int interval);
extern bool deploy_node(char *id);

#define CMD_BUILD  1
#define CMD_DEPLOY 2
#define CMD_STATUS 3
#define CMD_ALL    4

static UsysOption longOptions[] = {
    { "logs",        required_argument, 0, 'l' },
    { "config-file", required_argument, 0, 'c' },
    { "help",        no_argument, 0, 'h' },
    { "version",     no_argument, 0, 'v' },
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

    usys_puts("Usage: builder [build | deploy | status] [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-c, --config-file <file>      Builder config file");
    usys_puts("-v, --version                 Software version");
}

int main(int argc, char **argv) {

    int opt, optIdx;
    int cmd = CMD_ALL;
    char *debug      = DEF_LOG_LEVEL;
    char *configFile = DEF_CONFIG_FILE;

    Config *config = NULL;

    if (strcasecmp(argv[1], "build") == 0) {
        cmd = CMD_BUILD;
    } else if (strcasecmp(argv[1], "deploy") == 0) {
        cmd = CMD_DEPLOY;
    } else if (strcasecmp(argv[1], "status") == 0) {
        cmd = CMD_STATUS;
    } else {
        cmd = CMD_ALL;
    }

    /* Parsing command line args. */
    while (true) {

        opt = 0;
        optIdx = 0;

        opt = usys_getopt_long(argc, argv, "vh:c:p:l", longOptions,
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
            usys_puts(BUILDER_VERSION);
            usys_exit(0);
            break;

        case 'l':
            debug = optarg;
            set_log_level(debug);
            break;

        case 'c':
            configFile = optarg;
            if (!configFile) {
                usage();
                usys_exit(0);
            }
            break;

        default:
            usage();
            usys_exit(0);
        }
    }

    if (read_config_file(&config, configFile) != USYS_TRUE) {
        usys_log_error("Unable to read builder's config file: %s",
                       configFile);
        goto done;
    }

    if (cmd == CMD_ALL || cmd == CMD_BUILD) {

        /* build all systems */
        if (!build_all_systems(config->build->systemsList,
                               config->setup->ukamaRepo,
                               config->setup->authRepo)) {
            usys_log_error("Build (systems) error. Exiting ...");
            goto done;
        }

        /* build node(s) */
        if (!build_nodes(config->build->nodeCount,
                         config->setup->ukamaRepo,
                         config->build->nodeIDsList)) {
            usys_log_error("Build (node) error. Exiting ...");
            goto done;
        }

        if (cmd == CMD_BUILD) {
            free_config(config);
            return USYS_TRUE;
        }
    }

    if (cmd == CMD_ALL || cmd == CMD_DEPLOY) {

        usys_log_debug("Deploying the node(s) and system(s) ...");

        if (!deploy_node(config->build->nodeIDsList)) {
            usys_log_error("Unable to deploy the node. Existing ...");
            goto done;
        }

        if (!deploy_all_systems(config->fileName,
                                config->deploy,
                                config->setup->ukamaRepo,
                                config->setup->authRepo)) {
            usys_log_error("Unable to deploy the system. Exiting ...");
            goto done;
        }

        if (cmd == CMD_DEPLOY) {
            free_config(config);
            return USYS_TRUE;
        }
    }

    if (cmd == CMD_ALL || cmd == CMD_STATUS) {
        display_all_systems_status(config->deploy->systemsList,
                                   config->setup->statusInterval);
    }

done:
    free_config(config);
    return USYS_TRUE;
}

