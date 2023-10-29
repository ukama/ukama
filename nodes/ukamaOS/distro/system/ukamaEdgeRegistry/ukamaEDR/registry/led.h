/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */


#ifndef EDRDB_LED_H_
#define EDRDB_LED_H_

#include "inc/registry.h"
#include "ifmsg.h"

typedef struct {
	PData state;
	PData dimmer;
	PData ontime;
	PData cumactivepower;
	PData power;
	PData colour;
	PData units;
	PData applicationtype;
} LedData;

typedef struct {
	uint8_t instance;
	DevObj obj;
	LedData data;
} DRDBLedSchema;


int drdb_read_led_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);
int drdb_write_led_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);

void drdb_add_led_dev_to_reg(void* pdev) ;
void drdb_add_led_inst_to_reg(Device* dev, uint8_t inst, uint8_t subdev);
void free_led_data(void* data);
void* copy_led_data(void *pdata);
#endif /* EDRDB_LED_H_ */
