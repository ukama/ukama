/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <getopt.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>

#include "config.h"
#include "model.h"
#include "ret_mode.h"
#include "server.h"
#include "web_service.h"
#include "usys_log.h"
#include "version.h"

static volatile sig_atomic_t gRunning = 1;

typedef struct {
    EmuConfig *config;
    EmuModel *model;
} ServerThreadArg;

static void handle_signal(int sig)
{
    (void)sig;
    gRunning = 0;
}

static void *server_thread(void *arg)
{
    ServerThreadArg *threadArg;

    threadArg = arg;
    emu_server_run(threadArg->config, threadArg->model, &gRunning);

    return NULL;
}

static bool install_signal_handlers(void)
{
    struct sigaction sa;

    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = handle_signal;
    sigemptyset(&sa.sa_mask);

    if (sigaction(SIGINT, &sa, NULL) < 0) {
        return false;
    }

    if (sigaction(SIGTERM, &sa, NULL) < 0) {
        return false;
    }

    return true;
}

static void set_log_level(const char *level)
{
    int logLevel = USYS_LOG_INFO;

    if (level == NULL) {
        return;
    }

    if (!strcmp(level, "TRACE")) {
        logLevel = USYS_LOG_TRACE;
    } else if (!strcmp(level, "DEBUG")) {
        logLevel = USYS_LOG_DEBUG;
    } else if (!strcmp(level, "INFO")) {
        logLevel = USYS_LOG_INFO;
    }

    usys_log_set_level(logLevel);
}

static void usage(const char *prog)
{
    fprintf(stderr,
            "Usage: %s [options]\n"
            "\n"
            "Common options:\n"
            "  -c, --config FILE              config file\n"
            "  -l, --log-level LEVEL          TRACE|DEBUG|INFO\n"
            "  -v, --version                  print version\n"
            "  -h, --help                     show this help\n"
            "      --mode contract|ret        emulator mode (default contract)\n"
            "\n"
            "RET mode options:\n"
            "      --pty PATH                 PTY symlink path (default /tmp/aisg-ret0)\n"
            "      --vendor XX                two-character AISG vendor code\n"
            "      --serial SERIAL            serial number, max 17 chars\n"
            "      --requires-config BOOL     true/false\n"
            "      --initial-tilt DEG         initial tilt in degrees\n"
            "      --min-tilt DEG             min tilt in degrees\n"
            "      --max-tilt DEG             max tilt in degrees\n",
            prog);
}

int main(int argc, char **argv)
{
    EmuConfig config;
    EmuModel model;
    UInst web;
    ServerThreadArg threadArg;
    pthread_t tid;
    const char *configFile;
    int opt;
    int optIndex = 0;
    bool loadedConfig = false;

    static struct option longOpts[] = {
        { "config",          required_argument, 0, 'c' },
        { "log-level",       required_argument, 0, 'l' },
        { "version",         no_argument,       0, 'v' },
        { "help",            no_argument,       0, 'h' },
        { "mode",            required_argument, 0, 1000 },
        { "pty",             required_argument, 0, 1001 },
        { "vendor",          required_argument, 0, 1002 },
        { "serial",          required_argument, 0, 1003 },
        { "requires-config", required_argument, 0, 1004 },
        { "initial-tilt",    required_argument, 0, 1005 },
        { "min-tilt",        required_argument, 0, 1006 },
        { "max-tilt",        required_argument, 0, 1007 },
        { 0, 0, 0, 0 }
    };

    configFile = AISG_EMU_CONFIG_FILE;

    if (!emu_config_init(&config)) {
        return 1;
    }

    emu_model_init(&model);
    usys_log_set_level(LOG_INFO);
    usys_log_set_service(AISG_EMU_SERVICE_NAME);

    while ((opt = getopt_long(argc,
                              argv,
                              "c:l:hv",
                              longOpts,
                              &optIndex)) != -1) {
        switch (opt) {
        case 'c':
            configFile = optarg;
            break;
        case 'l':
            set_log_level(optarg);
            break;
        case 'v':
            printf("%s\n", VERSION);
            emu_model_free(&model);
            emu_config_free(&config);
            return 0;
        case 'h':
            usage(argv[0]);
            emu_model_free(&model);
            emu_config_free(&config);
            return 0;
        case 1000:
            if (!emu_config_set_mode(&config, optarg)) {
                fprintf(stderr, "invalid mode: %s\n", optarg);
                usage(argv[0]);
                emu_model_free(&model);
                emu_config_free(&config);
                return 1;
            }
            break;
        case 1001:
            snprintf(config.retPtyPath, sizeof(config.retPtyPath), "%s", optarg);
            break;
        case 1002:
            snprintf(config.retVendorCode, sizeof(config.retVendorCode), "%.2s", optarg);
            break;
        case 1003:
            snprintf(config.retSerial, sizeof(config.retSerial), "%s", optarg);
            break;
        case 1004:
            if (!emu_config_set_bool(&config.retRequiresConfig, optarg)) {
                fprintf(stderr, "invalid --requires-config value: %s\n", optarg);
                emu_model_free(&model);
                emu_config_free(&config);
                return 1;
            }
            break;
        case 1005:
            config.retInitialTiltTenths = emu_config_tilt_arg_to_tenths(
                optarg,
                config.retInitialTiltTenths);
            break;
        case 1006:
            config.retMinTiltTenths = emu_config_tilt_arg_to_tenths(
                optarg,
                config.retMinTiltTenths);
            break;
        case 1007:
            config.retMaxTiltTenths = emu_config_tilt_arg_to_tenths(
                optarg,
                config.retMaxTiltTenths);
            break;
        default:
            usage(argv[0]);
            emu_model_free(&model);
            emu_config_free(&config);
            return 1;
        }
    }

    if (configFile != NULL) {
        loadedConfig = emu_config_load(&config, configFile);
        if (!loadedConfig) {
            usys_log_warn("using default aisg-emu configuration");
        }
    }

    /* Re-parse long options that may intentionally override config file values. */
    optind = 1;
    while ((opt = getopt_long(argc,
                              argv,
                              "c:l:hv",
                              longOpts,
                              &optIndex)) != -1) {
        switch (opt) {
        case 1000:
            emu_config_set_mode(&config, optarg);
            break;
        case 1001:
            snprintf(config.retPtyPath, sizeof(config.retPtyPath), "%s", optarg);
            break;
        case 1002:
            snprintf(config.retVendorCode, sizeof(config.retVendorCode), "%.2s", optarg);
            break;
        case 1003:
            snprintf(config.retSerial, sizeof(config.retSerial), "%s", optarg);
            break;
        case 1004:
            emu_config_set_bool(&config.retRequiresConfig, optarg);
            break;
        case 1005:
            config.retInitialTiltTenths = emu_config_tilt_arg_to_tenths(
                optarg,
                config.retInitialTiltTenths);
            break;
        case 1006:
            config.retMinTiltTenths = emu_config_tilt_arg_to_tenths(
                optarg,
                config.retMinTiltTenths);
            break;
        case 1007:
            config.retMaxTiltTenths = emu_config_tilt_arg_to_tenths(
                optarg,
                config.retMaxTiltTenths);
            break;
        default:
            break;
        }
    }

    (void)loadedConfig;

    if (!install_signal_handlers()) {
        usys_log_error("failed to install signal handlers");
        emu_model_free(&model);
        emu_config_free(&config);
        return 1;
    }

    if (config.mode == EmuModeRet) {
        bool ok;

        ok = ret_mode_run(&config, &gRunning);
        emu_model_free(&model);
        emu_config_free(&config);
        return ok ? 0 : 1;
    }

    emu_model_load_scenario(&model, config.scenario);

    threadArg.config = &config;
    threadArg.model  = &model;

    if (pthread_create(&tid, NULL, server_thread, &threadArg) != 0) {
        usys_log_error("failed to create emulator socket thread");
        emu_model_free(&model);
        emu_config_free(&config);
        return 1;
    }

    start_web_service(&web, &config, &model);

    while (gRunning) {
        pause();
    }

    usys_log_info("stopping aisg-emu");

    stop_web_service(&web);
    pthread_join(tid, NULL);

    emu_model_free(&model);
    emu_config_free(&config);

    return 0;
}
