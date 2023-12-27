/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>

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

/* shutdown.c */
extern bool shutdown_all_systems(char *systems, char *ukamaRepo, char *authRepo);
extern bool shutdown_node(char *id);

#define CMD_BUILD  1
#define CMD_DEPLOY 2
#define CMD_STATUS 3
#define CMD_ALL    4
#define CMD_DOWN   5

#define TARGET_ALL     1
#define TARGET_NODES   2
#define TARGET_SYSTEMS 3

#define IS_S3_PATH(source) (strncmp((source), "s3://", 5) == 0 ? USYS_TRUE : USYS_FALSE)

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

    usys_puts("Usage: builder [nodes | systems]  "
              "[build | deploy | status | down] [options]");
    usys_puts("Options:");
    usys_puts("-h, --help                    Help menu");
    usys_puts("-l, --logs <TRACE|DEBUG|INFO> Log level for the process");
    usys_puts("-c, --config-file <file>      Builder config file");
    usys_puts("-v, --version                 Software version");
}

void processArguments(int argc, char *argv[], int *target, int *cmd) {

    bool isTargetSet = USYS_FALSE;
    bool isCmdSet    = USYS_FALSE;

    *target = TARGET_ALL;
    *cmd    = CMD_ALL;

    for (int i = 1; i < argc; i++) {
        if (strcasecmp(argv[i], "nodes") == 0 ||
            strcasecmp(argv[i], "systems") == 0) {

            if (!isTargetSet) {
                if (strcasecmp(argv[1], "nodes") == 0) {
                    *target = TARGET_NODES;
                } else if (strcasecmp(argv[i], "systems") == 0) {
                    *target = TARGET_SYSTEMS;
                }
                isTargetSet = USYS_TRUE;
            }
        } else if (strcasecmp(argv[i], "build") == 0 ||
                   strcasecmp(argv[i], "deploy") == 0 ||
                   strcasecmp(argv[i], "status") == 0 ||
                   strcasecmp(argv[i], "down") == 0) {

            if (!isCmdSet) {
                if (strcasecmp(argv[1], "build") == 0) {
                    *cmd = CMD_BUILD;
                } else if (strcasecmp(argv[1], "deploy") == 0) {
                    *cmd = CMD_DEPLOY;
                } else if (strcasecmp(argv[1], "status") == 0) {
                    *cmd = CMD_STATUS;
                } else if (strcasecmp(argv[1], "down") == 0) {
                    *cmd = CMD_DOWN;
                }
                isCmdSet = USYS_TRUE;
            }
        } else if (strcasecmp(argv[i], "help") == 0) {
            usage();
            usys_exit(0);
        }
    }
}

void extract_filename(const char *path, char *filename, size_t size) {

    const char *lastSlash = strrchr(path, '/');

    if (lastSlash) {
        strncpy(filename, lastSlash + 1, size);
        filename[size - 1] = '\0';
    } else {
        strncpy(filename, path, size);
    }
}

bool is_s3path_and_fetch(char *source, char *dest) {

    char runMe[2*1024+1] = {0};
    char fileName[1024] = {0};

    extract_filename(source, fileName, sizeof(fileName));
    snprintf(runMe, sizeof(runMe), "aws s3 cp %s %s/%s",
             source, dest, fileName);

    if (system(runMe) == 0) {
        usys_log_debug("File successfully copied from S3 to local"
                       "%s -> %s", source, dest);
        return USYS_TRUE;
    } else {
        usys_log_error("Unable to copy from S3 (%s) to local (%s)",
                       source, dest);
        return USYS_FALSE;
    }
}

bool fetch_all_s3_files(char *kernel, char *initRAM, char *disk) {

    if (IS_S3_PATH(kernel)) {
        if (!is_s3path_and_fetch(kernel, "./scripts")) {
            usys_log_error("Unable to fetch: %s", kernel);
            return USYS_FALSE;
        }
    }

    if (IS_S3_PATH(initRAM)) {
        if (!is_s3path_and_fetch(initRAM, "./scripts")) {
            usys_log_error("Unable to fetch: %s", initRAM);
            return USYS_FALSE;
        }
    }

    if (IS_S3_PATH(disk)) {
        if (!is_s3path_and_fetch(disk, "./scripts")) {
            usys_log_error("Unable to fetch: %s", disk);
            return USYS_FALSE;
        }
    }

    return USYS_TRUE;
}

int main(int argc, char **argv) {

    int opt, optIdx;
    int cmd = CMD_ALL, target = TARGET_ALL;
    char *debug      = DEF_LOG_LEVEL;
    char *configFile = DEF_CONFIG_FILE;

    Config *config = NULL;

    processArguments(argc, argv, &cmd, &target);

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

        if (target == TARGET_ALL || target == TARGET_SYSTEMS) {
            if (!build_all_systems(config->build->systemsList,
                                   config->setup->ukamaRepo,
                                   config->setup->authRepo)) {
                usys_log_error("Build (systems) error. Exiting ...");
                goto done;
            }
        }

        if (target == TARGET_ALL || target == TARGET_NODES) {

            if (config->build->kernelImage != NULL &&
                config->build->initRAMImage != NULL &&
                config->build->diskImage != NULL) {

                if (!fetch_all_s3_files(config->build->kernelImage,
                                        config->build->initRAMImage,
                                        config->build->diskImage)) {
                    usys_log_error("Unable to fetch img files");
                    goto done;
                }

            } else {
                if (!build_nodes(config->build->nodeCount,
                                 config->setup->ukamaRepo,
                                 config->build->nodeIDsList)) {
                    usys_log_error("Build (node) error. Exiting ...");
                    goto done;
                }
            }

            if (cmd == CMD_BUILD) {
                free_config(config);
                return USYS_TRUE;
            }
        }
    }

    if (cmd == CMD_ALL || cmd == CMD_DEPLOY) {

        usys_log_debug("Deploying the node(s) and system(s) ...");

        if (target == TARGET_ALL || target == TARGET_NODES) {
            if (!deploy_node(config->build->nodeIDsList)) {
                usys_log_error("Unable to deploy the node. Existing ...");
                goto done;
            }
        }

        if (target == TARGET_ALL || target == TARGET_SYSTEMS) {
            if (!deploy_all_systems(config->fileName,
                                    config->deploy,
                                    config->setup->ukamaRepo,
                                    config->setup->authRepo)) {
                usys_log_error("Unable to deploy the system. Exiting ...");
                goto done;
            }
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

    if (cmd == CMD_ALL || cmd == CMD_DOWN) {

        if (target == TARGET_ALL || target == TARGET_NODES) {
            if (!shutdown_node(config->deploy->nodeIDsList)) {
                usys_log_error("Node Shutdown FAILED: %s Try manually",
                               config->deploy->nodeIDsList);
            }
        }

        if (target == TARGET_ALL || target == TARGET_SYSTEMS) {
            if (!shutdown_all_systems(config->deploy->systemsList,
                                      config->setup->ukamaRepo,
                                      config->setup->authRepo)) {
                usys_log_error("Systems Shutdown FAILED");
                goto done;
            }
        }
    }

done:
    free_config(config);
    return USYS_TRUE;
}
