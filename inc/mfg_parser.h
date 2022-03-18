/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#ifndef INC_MFG_PARSER_H_
#define INC_MFG_PARSER_H_

#include "schema.h"
#include "json_types.h"

int parser_schema_init(JSONInput* json_ip);

void parser_schema_exit();

StoreSchema *parser_get_mfg_data_by_uuid(char *puuid);

#endif /* INC_MFG_PARSER_H_ */
