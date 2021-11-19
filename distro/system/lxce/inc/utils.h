/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * utilies functions
 */

#ifndef LXCE_UTILS_H
#define LXCE_UTILS_H

#define PARENT_SOCKET 0
#define CHILD_SOCKET  1

int set_integer_object_value(json_t *json, int *param, char *objName,
			     int mandatory, int defValue);
int set_str_object_value(json_t *json, char **param, char *objName,
			 int mandatory, char *defValue);
int namespaces_flag(char *ns);

#endif /* LXCE_UTILS_H */
