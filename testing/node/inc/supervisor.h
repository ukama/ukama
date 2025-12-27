/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef SUPERVISOR_H
#define SUPERVISOR_H

#include "config.h"

#define SVISOR_FILENAME                 "supervisor.conf"

#define SVISOR_HEADER                   "; Ukama Node Supervisor config file \n\
[unix_http_server]\n\
 file=/var/run/supervisor.sock\n\
 chmod=0700\n\n"

#define SVISOR_SVISORD                  "[supervisord] \n\
 pidfile=/var/run/supervisord.pid\n\
 nodaemon=true \n\
 loglevel=warn\n\n"

#define SVISOR_RPCINTERFACE             "[rpcinterface:supervisor]\n\
 supervisor.rpcinterface_factory=\
supervisor.rpcinterface:make_main_rpcinterface\n\n"

#define SVISOR_SVISOR_CTL               "[supervisorctl]\n\
 serverurl=unix:///var/run/supervisor.sock\n\n"

#define SVISOR_INCLUDE                  "[include]\n\
 files = /etc/supervisor/conf.d/*.conf\n\n"

#define SVISOR_KICKSTART                "[program:kickstart]\n\
 command=/bin/kickstart.sh\n\
 autostart=true\n\
 autorestart=false\n\
 startretries=0\n\
 exitcodes=0\n\
 stdout_logfile=/dev/stdout\n\
 stdout_logfile_maxbytes=0\n\
 stderr_logfile=/dev/stderr\n\n\
 stderr_logfile_maxbytes=0\n\n"

#define SVISOR_STDOUT_LOGFILE_MAXBYTES  "%s stdout_logfile_maxbytes=0\n"
#define SVISOR_STDERR_LOGFILE_MAXBYTES  "%s stderr_logfile_maxbytes=0\n"
#define SVISOR_GROUP_ONBOOT             "[group:on-boot]\n programs=%s\n\n"
#define SVISOR_GROUP_SYSSVC             "[group:sys-service]\n programs=%s\n\n"

#define SVISOR_PROGRAM                  "[program:%s_%s]\n"
#define SVISOR_COMMAND                  "%s command=%s/%s \n"
#define SVISOR_COMMAND_WITH_ARGS        "%s command=%s/%s %s \n"
#define SVISOR_AUTOSTART                "%s autostart=%s\n"
#define SVISOR_AUTORESTART              "%s autorestart=%s\n"
#define SVISOR_STARTRETRIES             "%s startretries=%d\n"
#define SVISOR_EXITCODES                "%s exitcodes=0\n"
#define SVISOR_STARTSECS                "%s startsecs=%d\n"
#define SVISOR_STDERR_LOGFILE           "%s stderr_logfile=/dev/stderr\n"
#define SVISOR_STDOUT_LOGFILE           "%s stdout_logfile=/dev/stdout\n"
#define SVISOR_MAX_SIZE                 4096
#define SVISOR_GROUP_LIST_MAX_SIZE      256
#define SVISOR_GROUP_ON_BOOT            "on-boot"
#define SVISOR_GROUP_SYS_SVC            "sys-service"
#define SVISOR_GROUP_NIL                NULL

int create_supervisor_config(Configs *configs);
void purge_supervisor_config(char *fileName);

#endif /* SUPERVISOR_H */
