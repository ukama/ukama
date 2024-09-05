/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
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

bool is_valid_json(const char *json_string);
bool make_path(const char* path);
int move_dir(const char *source, const char *destination);
int remove_dir(const char *path);
int clean_empty_dir(char* path);
bool create_config(ConfigData* c);
bool remove_config(ConfigData *c);
int create_backup_config();
int restore_config() ;
int store_config(char* version);
int prepare_for_new_config(ConfigData* c);

#endif /* UTIL_H_ */
