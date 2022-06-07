/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef SUPERVISOR_H
#define SUPERVISOR_H

#include "config.h"

#define SVISOR_FILENAME                 "supervisor.conf"
#define SVISOR_HEADER                   "[supervisord] \n nodaemon=true \n\n"

#define SVISOR_PROGRAM                  "[program: %s_%s]\n"
#define SVISOR_COMMAND                  "%s command=%s/%s\n"
#define SVISOR_COMMAND_WITH_ARGS        "%s command=%s/%s %s\n"
#define SVISOR_AUTOSTART                "%s autostart=%s\n"
#define SVISOR_AUTORESTART              "%s autorestart=%s\n"
#define SVISOR_STDERR_LOGFILE           "%s stderr_logfile=%s\n"
#define SVISOR_STDERR_LOGFILE_MAX_BYTES "%s stderr_logfile_maxbytes=%d\n"
#define SVISOR_STDOUT_LOGFILE           "%s stdout_logfile=%s\n"
#define SVISOR_STDOUT_LOGFILE_MAX_BYTES "%s stdout_logfile_maxbytes=%d\n"

#define SVISOR_MAX_SIZE                2048
#define SVISOR_DEFAULT_STDOUT          "/dev/stdout"
#define SVISOR_DEFAULT_STDERR          "/dev/stdout"
#define SVISOR_DEFAULT_STDERR_MAXBYTES 0
#define SVISOR_DEFAULT_STDOUT_MAXBYTES 0

int create_supervisor_config(Configs *configs);
void purge_supervisor_config(char *fileName);

#endif /* SUPERVISOR_H */
