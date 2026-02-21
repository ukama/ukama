/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <sys/wait.h>
#include <unistd.h>

#include "actions.h"
#include "deviced.h"

#include "usys_log.h"

static int _run_cmd(const char *cmd) {

    int status = 0;
    pid_t pid = 0;
    char *argv[4];

    if (!cmd || !*cmd) return STATUS_NOK;

    pid = fork();
    if (pid < 0) {
        return STATUS_NOK;
    }

    if (pid == 0) {
        argv[0] = "/bin/sh";
        argv[1] = "-c";
        argv[2] = (char *)cmd;
        argv[3] = NULL;
        execv(argv[0], argv);
        _exit(127);
    }

    if (waitpid(pid, &status, 0) < 0) {
        return STATUS_NOK;
    }

    if (WIFEXITED(status) && WEXITSTATUS(status) == 0) {
        return STATUS_OK;
    }

    return STATUS_NOK;
}

int actions_service_apply(Config *config, ControlState desired) {

    int ret = STATUS_NOK;
    const char *cmd1 = NULL;
    const char *cmd2 = NULL;

    (void)config;

    if (desired == CONTROL_STATE_OFF) {
        cmd1 = "/epc_package/epc/utils/packaging ./int_setup.sh 0 0 0";
        usys_log_info("service: off");
        ret = _run_cmd(cmd1);
        return ret;
    }

    cmd1 = "/epc_package/epc/utils/packaging sudo ./int_setup.sh 0 0 1";
    cmd2 = "/epc_package/newman run config_srs_ue_1.json";

    usys_log_info("service: on");

    if (_run_cmd(cmd1) != STATUS_OK) {
        return STATUS_NOK;
    }

    if (_run_cmd(cmd2) != STATUS_OK) {
        return STATUS_NOK;
    }

    return STATUS_OK;
}
