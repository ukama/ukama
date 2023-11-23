/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef UTILS_VNODEALERT_H_
#define UTILS_VNODEALERT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "irqdb.h"

#include <errno.h>
#include <poll.h>
#include <fcntl.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <sys/inotify.h>
#include <sys/types.h>
#include <unistd.h>


#define EVENT_SIZE  		(sizeof (struct inotify_event) )
#define EVENT_BUF_LEN     	(1024 * ( EVENT_SIZE + 16))

#define errExit(msg)    	do { perror(msg); exit(EXIT_FAILURE); \
              } while (0)

#define ERR					-1

/**
 * @fn      int poll_file(IRQCfg*)
 * @brief   Wait for the Inotification event for the file.
 *
 * @param   cfg
 * @return  On success, 0
 *          On failure, -1
 */
int poll_file(IRQCfg* cfg);

#ifdef __cplusplus
}
#endif

#endif /* UTILS_VNODEALERT_H_ */
