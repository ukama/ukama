/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "config.h"
#include "switchemu.h"

#include "version.h"

#ifdef __has_include
#  if __has_include("usys_getopt.h")
#    include "usys_getopt.h"
#    define EMU_GETOPT_LONG usys_getopt_long
#    define EMU_OPTARG optarg
#  else
#    include <getopt.h>
#    define UsysOption struct option
#    define EMU_GETOPT_LONG getopt_long
#    define EMU_OPTARG optarg
#  endif
#else
#  include <getopt.h>
#  define UsysOption struct option
#  define EMU_GETOPT_LONG getopt_long
#  define EMU_OPTARG optarg
#endif

#ifdef __has_include
#  if __has_include("usys_string.h")
#    include "usys_string.h"
#    define EMU_STREQ(a, b) (strcmp((a), (b)) == 0)
#  else
#    define EMU_STREQ(a, b) (strcmp((a), (b)) == 0)
#  endif
#else
#  define EMU_STREQ(a, b) (strcmp((a), (b)) == 0)
#endif

#ifdef __has_include
#  if __has_include("usys_log.h")
#    include "usys_log.h"
#    define EMU_LOG_TRACE USYS_LOG_TRACE
#    define EMU_LOG_DEBUG USYS_LOG_DEBUG
#    define EMU_LOG_INFO  USYS_LOG_INFO
#    define EMU_LOG_ERROR USYS_LOG_ERROR
#  else
#    define EMU_LOG_TRACE 0
#    define EMU_LOG_DEBUG 1
#    define EMU_LOG_INFO  2
#    define EMU_LOG_ERROR 3
#  endif
#else
#  define EMU_LOG_TRACE 0
#  define EMU_LOG_DEBUG 1
#  define EMU_LOG_INFO  2
#  define EMU_LOG_ERROR 3
#endif

static UsysOption longOptions[] = {
    { "http-port",      required_argument, 0, 'p' },
    { "snmp-port",      required_argument, 0, 's' },
    { "tftp-port",      required_argument, 0, 't' },
    { "bind-address",   required_argument, 0, 'b' },
    { "state-file",     required_argument, 0, 'f' },
    { "scenario",       required_argument, 0, 'S' },
    { "notify-host",    required_argument, 0, 'H' },
    { "notify-port",    required_argument, 0, 'N' },
    { "notify-path",    required_argument, 0, 'P' },
    { "disable-notify", no_argument,       0, 'd' },
    { "logs",           required_argument, 0, 'l' },
    { "help",           no_argument,       0, 'h' },
    { "version",        no_argument,       0, 'v' },
    { 0, 0, 0, 0 }
};

void config_init_default(EmuConfig *cfg) {
    memset(cfg, 0, sizeof(*cfg));

    cfg->httpPort      = DEF_HTTP_PORT;
    cfg->snmpPort      = DEF_SNMP_PORT;
    cfg->tftpPort      = DEF_TFTP_PORT;
    cfg->notifyPort    = DEF_NOTIFY_PORT;
    cfg->logLevel      = EMU_LOG_INFO;
    cfg->notifyEnabled = 1;

    snprintf(cfg->bindAddr, sizeof(cfg->bindAddr), "%s", DEF_BIND_ADDRESS);
    snprintf(cfg->stateFile, sizeof(cfg->stateFile), "%s", DEF_STATE_FILE);
    snprintf(cfg->scenario, sizeof(cfg->scenario), "%s", DEF_SCENARIO);
    snprintf(cfg->notifyHost, sizeof(cfg->notifyHost), "%s", DEF_NOTIFY_HOST);
    snprintf(cfg->notifyPath, sizeof(cfg->notifyPath), "%s", DEF_NOTIFY_PATH);
}

int config_parse_log_level(const char *slevel) {
    if (slevel == NULL) {
        return EMU_LOG_INFO;
    }

    if (EMU_STREQ(slevel, "TRACE")) {
        return EMU_LOG_TRACE;
    }
    if (EMU_STREQ(slevel, "DEBUG")) {
        return EMU_LOG_DEBUG;
    }
    if (EMU_STREQ(slevel, "INFO")) {
        return EMU_LOG_INFO;
    }
    if (EMU_STREQ(slevel, "ERROR")) {
        return EMU_LOG_ERROR;
    }

    return EMU_LOG_INFO;
}

static void config_load_from_env(EmuConfig *cfg) {
    const char *env = NULL;

    env = getenv("SWITCHEMU_HTTP_PORT");
    if (env != NULL) cfg->httpPort = atoi(env);

    env = getenv("SWITCHEMU_SNMP_PORT");
    if (env != NULL) cfg->snmpPort = atoi(env);

    env = getenv("SWITCHEMU_TFTP_PORT");
    if (env != NULL) cfg->tftpPort = atoi(env);

    env = getenv("SWITCHEMU_NOTIFY_PORT");
    if (env != NULL) cfg->notifyPort = atoi(env);

    env = getenv("SWITCHEMU_BIND_ADDRESS");
    if (env != NULL) snprintf(cfg->bindAddr, sizeof(cfg->bindAddr), "%s", env);

    env = getenv("SWITCHEMU_STATE_FILE");
    if (env != NULL) snprintf(cfg->stateFile, sizeof(cfg->stateFile), "%s", env);

    env = getenv("SWITCHEMU_SCENARIO");
    if (env != NULL) snprintf(cfg->scenario, sizeof(cfg->scenario), "%s", env);

    env = getenv("SWITCHEMU_NOTIFY_HOST");
    if (env != NULL) snprintf(cfg->notifyHost, sizeof(cfg->notifyHost), "%s", env);

    env = getenv("SWITCHEMU_NOTIFY_PATH");
    if (env != NULL) snprintf(cfg->notifyPath, sizeof(cfg->notifyPath), "%s", env);

    env = getenv("SWITCHEMU_NOTIFY_ENABLED");
    if (env != NULL) cfg->notifyEnabled = atoi(env) ? 1 : 0;

    env = getenv("SWITCHEMU_LOG_LEVEL");
    if (env != NULL) cfg->logLevel = config_parse_log_level(env);
}

void config_usage(void) {
    puts("Usage: switchemu.d [options]");
    puts("Options:");
    puts("-h, --help                    Help menu");
    puts("-v, --version                 Software version");
    puts("-l, --logs <TRACE|DEBUG|INFO|ERROR>");
    puts("                              Log level");
    puts("-p, --http-port <port>        HTTP debug service port");
    puts("-s, --snmp-port <port>        SNMP emulator port");
    puts("-t, --tftp-port <port>        TFTP emulator port");
    puts("-b, --bind-address <addr>     Bind address for listeners");
    puts("-f, --state-file <path>       State persistence file");
    puts("-S, --scenario <name>         Startup scenario");
    puts("-H, --notify-host <host>      notify.d host");
    puts("-N, --notify-port <port>      notify.d port");
    puts("-P, --notify-path <path>      notify.d alarm path");
    puts("-d, --disable-notify          Disable alarm delivery");
    puts("");
    puts("Environment:");
    puts("  SWITCHEMU_HTTP_PORT");
    puts("  SWITCHEMU_SNMP_PORT");
    puts("  SWITCHEMU_TFTP_PORT");
    puts("  SWITCHEMU_BIND_ADDRESS");
    puts("  SWITCHEMU_STATE_FILE");
    puts("  SWITCHEMU_SCENARIO");
    puts("  SWITCHEMU_NOTIFY_HOST");
    puts("  SWITCHEMU_NOTIFY_PORT");
    puts("  SWITCHEMU_NOTIFY_PATH");
    puts("  SWITCHEMU_NOTIFY_ENABLED");
    puts("  SWITCHEMU_LOG_LEVEL");
}

void config_version(void) {
#ifdef VERSION
    puts(VERSION);
#else
    puts("unknown");
#endif
}

int config_load(EmuConfig *cfg, int argc, char **argv) {
    int opt    = 0;
    int optIdx = 0;

    config_init_default(cfg);
    config_load_from_env(cfg);

    while (1) {
        opt = EMU_GETOPT_LONG(argc, argv, "hvl:p:s:t:b:f:S:H:N:P:d",
                              longOptions, &optIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
        case 'h':
            config_usage();
            exit(0);
        case 'v':
            config_version();
            exit(0);
        case 'l':
            cfg->logLevel = config_parse_log_level(EMU_OPTARG);
            break;
        case 'p':
            cfg->httpPort = atoi(EMU_OPTARG);
            break;
        case 's':
            cfg->snmpPort = atoi(EMU_OPTARG);
            break;
        case 't':
            cfg->tftpPort = atoi(EMU_OPTARG);
            break;
        case 'b':
            snprintf(cfg->bindAddr, sizeof(cfg->bindAddr), "%s", EMU_OPTARG);
            break;
        case 'f':
            snprintf(cfg->stateFile, sizeof(cfg->stateFile), "%s", EMU_OPTARG);
            break;
        case 'S':
            snprintf(cfg->scenario, sizeof(cfg->scenario), "%s", EMU_OPTARG);
            break;
        case 'H':
            snprintf(cfg->notifyHost, sizeof(cfg->notifyHost), "%s", EMU_OPTARG);
            break;
        case 'N':
            cfg->notifyPort = atoi(EMU_OPTARG);
            break;
        case 'P':
            snprintf(cfg->notifyPath, sizeof(cfg->notifyPath), "%s", EMU_OPTARG);
            break;
        case 'd':
            cfg->notifyEnabled = 0;
            break;
        default:
            config_usage();
            return STATUS_NOK;
        }
    }

    return STATUS_OK;
}
