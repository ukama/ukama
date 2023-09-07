/*
 * parser.h
 *
 *  Created on: Sep 9, 2020
 *      Author: root
 */

#ifndef UKDB_IDB_SPARSER_H_
#define UKDB_IDB_SPARSER_H_

#include "headers/ubsp/ukdblayout.h"

#define MAX_JSON_SCHEMA     5

int parser_schema_init(JSONInput* json_ip);
void parser_schema_exit();
UKDB* parser_get_mfg_data_by_uuid(char* puuid);

#endif /* UKDB_IDB_SPARSER_H_ */
