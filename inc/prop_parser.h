/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_PROP_PARSER_H_
#define INC_PROP_PARSER_H_


#include "device.h"
#include "json_types.h"
#include "noded_macros.h"
#include "property.h"
#include "schema.h"

typedef struct __attribute__((__packed__)) {
    char name[NAME_LENGTH];
    Version ver;
    uint16_t propCount;
    Property* prop;
} PropertyMap;

/* No free require for these function */

int prop_parser_init(char *ip);
int prop_parser_get_count(char* name);
void prop_parser_exit();

Property* prop_parser_get_table(char* name);
Version* prop_parser_get_table_version(char* name);

#endif /* INC_PROP_PARSER_H_ */
