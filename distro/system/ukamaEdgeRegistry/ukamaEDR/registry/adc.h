/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */


#ifndef EDRDB_ADC_H_
#define EDRDB_ADC_H_

#include "inc/registry.h"
#include "headers/edr/ifmsg.h"

typedef struct {
	PData outputcurr; /*  TODO: Check if we just report other value like in dB */
	PData minrange;
	PData maxrange;
	//PData units;
	PData applicationtype;
} AdcData;

typedef struct  {
	uint8_t instance;
	DevObj obj;
	AdcData data;
} DRDBAdcSchema;


int drdb_read_adc_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);
int drdb_write_adc_inst_data_from_dev(DRDBSchema* reg, MsgFrame* rqmsg);
int drdb_add_ads1015_inst_to_reg(Device *dev, Property *prop, uint8_t inst,
		uint8_t subdev);

void free_adc_data(void* data);
void drdb_add_adc_dev_to_reg(void* pdev);

void* copy_adc_data(void *pdata);
#endif /* EDRDB_ADC_H_ */
