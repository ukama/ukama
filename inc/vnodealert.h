/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef UTILS_VNODEALERT_H_
#define UTILS_VNODEALERT_H_

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

int poll_file(IRQCfg* cfg);
void clean_poll();

#endif /* UTILS_VNODEALERT_H_ */
