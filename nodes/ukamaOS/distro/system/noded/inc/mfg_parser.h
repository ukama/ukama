/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_MFG_PARSER_H_
#define INC_MFG_PARSER_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "schema.h"
#include "json_types.h"

/**
 * @fn      int parser_schema_init(JSONInput*)
 * @brief   Reads the JSON schema and parses it.
 *
 * @param   json_ip
 * @return  On success, 0
 *          On failure, non zero value.
 */
int parser_schema_init(JSONInput* json_ip);

/**
 * @fn      void parser_schema_exit()
 * @brief   Release all the memory allocated by parser to tore the value read
 *          after parsing is completed successfully.
 *
 */
void parser_schema_exit();

/**
 * @fn      StoreSchema parser_get_mfg_data_by_uuid*(char*)
 * @brief   Reads the schema value read by parser for the module with UUID
 *          mentioned in the argument.
 *
 * @param   puuid
 * @return  On success, pointer to schema.
 *          on failure, NULL
 */
StoreSchema *parser_get_mfg_data_by_uuid(char *puuid);

#ifdef __cplusplus
}
#endif

#endif /* INC_MFG_PARSER_H_ */
