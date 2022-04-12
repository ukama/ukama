/*
 * pparser.h
 *
 *  Created on: Sep 12, 2020
 *      Author: root
 */

#ifndef UTILS_PPARSER_H_
#define UTILS_PPARSER_H_

#include "headers/ubsp/devices.h"
#include "headers/ubsp/property.h"
#include "headers/ubsp/ukdblayout.h"

#include <string.h>

#define MAX_JSON_DEVICE         32
#define PROP_NAME_LENGTH        32

typedef struct __attribute__((__packed__)) {
    char name[PROP_NAME_LENGTH];
    Version ver;
    uint16_t prop_count;
    Property* prop;
} PropertyMap;

/* No free require for these function */

int parser_property_init(char *ip);
int parser_get_property_count(char* name);
void parser_property_exit();

Property* parser_get_property_table(char* name);
Version* parser_get_property_table_version(char* name);

#endif /* UTILS_PPARSER_H_ */
