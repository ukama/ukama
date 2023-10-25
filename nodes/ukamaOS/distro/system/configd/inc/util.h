/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef UTIL_H_
#define UTIL_H_

#include <stdio.h>
#include <stdlib.h>

#include "jansson.h"
#include "configd.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"

int is_valid_json(const char *json_string);

int make_path(const char* path);

int move_dir(const char *source, const char *destination);

int remove_dir(const char *path);


int create_config(ConfigData* c);

int create_backup_config();

int restore_config() ;

int store_config(char* version);

int prepare_for_new_config(ConfigData* c);

#endif /* UTIL_H_ */
