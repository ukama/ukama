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

#define SVISOR_HEADER                   "; Ukama Node Supervisor config file \n\
[unix_http_server]\n\
 file=/var/run/supervisor.sock\n\
 chmod=0700\n\n"

#define SVISOR_SVISORD 					"[supervisord] \n\
 logfile=/var/log/supervisor/supervisord.log\n\
 pidfile=/var/run/supervisord.pid\n\
 childlogdir=/var/log/supervisor\n\
 nodaemon=true \n\n"

#define SVISOR_RPCINTERFACE 			"[rpcinterface:supervisor]\n\
 supervisor.rpcinterface_factory=\
supervisor.rpcinterface:make_main_rpcinterface\n\n"

#define SVISOR_SVISOR_CTL				"[supervisorctl]\n\
 serverurl=unix:///var/run/supervisor.sock\n\n"

#define SVISOR_INCLUDE					"[include]\n\
 files = /etc/supervisor/conf.d/*.conf\n\n"

#define SVISOR_KICKSTART				"[program:kickstart]\n\
 command=/bin/kickstart.sh\n\
 autostart=true\n\
 autorestart=false\n\
 startretries=1\n\
 stderr_logfile=/var/log/supervisor/%(program_name)s_stderr.log\n\
 stdout_logfile=/var/log/supervisor/%(program_name)s_stdout.log\n\
 stderr_logfile_maxbytes=10MB\n\
 stdout_logfile_maxbytes=10MB\n\
 exitcodes=0\n\n"

#define SVISOR_GROUP_ONBOOT				"[group:on-boot]\n programs=%s\n"
#define SVISOR_GROUP_SYSSVC				"[group:sys-service]\n programs=%s\n"
#define SVISOR_PROGRAM                  "[program:%s_%s]\n"
#define SVISOR_COMMAND                  "%s command=%s/%s \n"
#define SVISOR_COMMAND_WITH_ARGS        "%s command=%s/%s %s \n"
#define SVISOR_AUTOSTART                "%s autostart=false\n"
#define SVISOR_AUTORESTART              "%s autorestart=false\n"
#define SVISOR_STARTRETRIES             "%s startretries=%d\n"
#define SVISOR_STDERR_LOGFILE           "%s stderr_logfile=/var/log/supervisor/%(program_name)s_stderr.log\n"
#define SVISOR_STDERR_LOGFILE_MAX_BYTES "%s stderr_logfile_maxbytes=%dMB\n"
#define SVISOR_STDOUT_LOGFILE           "%s stdout_logfile=/var/log/supervisor/%(program_name)s_stdout.log\n"
#define SVISOR_STDOUT_LOGFILE_MAX_BYTES "%s stdout_logfile_maxbytes=%dMB\n"

#define SVISOR_MAX_SIZE                 4096
#define SVISOR_DEFAULT_STDOUT           "/dev/stdout"
#define SVISOR_DEFAULT_STDERR           "/dev/stdout"
#define SVISOR_DEFAULT_STDERR_MAXBYTES  10
#define SVISOR_DEFAULT_STDOUT_MAXBYTES  10

#define SVISOR_GROUP_LIST_MAX_SIZE      256

#define SVISOR_GROUP_ON_BOOT            "on-boot"
#define SVISOR_GROUP_SYS_SVC            "sys-service"
#define SVISOR_GROUP_NIL                NULL

#define SVISOR_RUNME_SH                 "/bin/runme.sh"
#define SVISOR_WAITFOR_SH               "/bin/waitfor.sh"

int create_supervisor_config(Configs *configs);
void purge_supervisor_config(char *fileName);

#endif /* SUPERVISOR_H */
