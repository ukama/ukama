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

#ifdef __cplusplus
extern "C" {
#endif

#include "device.h"
#include "json_types.h"
#include "noded_macros.h"
#include "property.h"
#include "schema.h"

/**
 * @struct Property Map
 * @brief  Stores the properties read from the property config files.
 *         These properties defines what all sensor specific data can be
 *         controlled from software.
 */
typedef struct __attribute__((__packed__)) {
    char name[NAME_LENGTH];
    Version ver;
    uint16_t propCount;
    Property* prop;
} PropertyMap;

/**
 * @fn      int prop_parser_init(char*)
 * @brief   Reads the property config file and parses it.
 *          Stores the properties value to global variable gPropMap
 *
 * @param   ip
 * @return  On success, 0
 *          On failure, non zero value
 */
int prop_parser_init(char *ip);

/**
 * @fn      int prop_parser_get_count(char*)
 * @brief   Reads the properties count in property config for sensor name.
 *
 * @param   name
 * @return  On success, 0
 *          On failure, non zero value
 */
int prop_parser_get_count(char* name);

/**
 * @fn      void prop_parser_exit()
 * @brief   destruct all the memory allocated for parser data to store values.
 *
 */
void prop_parser_exit();

/**
 * @fn      Property prop_parser_get_table*(char*)
 * @brief   Read the property table for the sensor name from the global
 *          variable gMapProp.
 *
 * @param   name
 * @return  On success, 0
 *          On failure, non zero value
 */
Property* prop_parser_get_table(char* name);

/**
 * @fn      Version prop_parser_get_table_version*(char*)
 * @brief   Read the property table version for the sensor name from the global
 *          variable gMapProp.
 *
 * @param   name
 * @return  On success, 0
 *          On failure, non zero value
 */
Version* prop_parser_get_table_version(char* name);

#ifdef __cplusplus
}
#endif

#endif /* INC_PROP_PARSER_H_ */
