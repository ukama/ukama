/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef EDRDB_GPIO_H_
#define EDRDB_GPIO_H_

#include "inc/registry.h"
#include "ifmsg.h"

#define GPIO_TYPE_INPUT 	"in\n"
#define GPIO_TYPE_OUTPUT 	"out\n"
typedef struct {
	PData direction;
	PData state;
	PData counter;
	PData polarity;
	PData debounce;
	PData edge;
	PData ontime;
	PData offtime;
	PData applicationtype;
	PData sensortype;
} DigitalData;

typedef struct {
	PData state;
	PData counter;
	PData ontime;
	PData offtime;
	PData applicationtype;
} DigitalIpData;

typedef struct {
	PData state;
	PData polarity;
	PData applicationtype;
} DigitalOpData;

typedef struct {
	uint8_t instance;
	DevObj obj;
	DigitalData data;
} DRDBDigitalIOSchema;

int drdb_read_gpio_inst_data_from_dev(DRDBSchema* reg, MsgFrame* digt);
int drdb_write_gpio_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);

void drdb_add_gpio_dev_to_reg(void* pdev);
void drdb_add_gpio_inst_to_reg(Device* dev, uint8_t inst, uint8_t subdev);
void free_gpio_data (void* data);
void* copy_gpio_data(void *pdata);

#endif /* EDRDB_GPIO_H_ */
