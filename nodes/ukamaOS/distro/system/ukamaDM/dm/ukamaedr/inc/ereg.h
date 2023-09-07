/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_EREG_H_
#define INC_EREG_H_

#include <stdbool.h>
#include <stdint.h>
#include <string.h>

int ereg_exec_sensor(uint16_t inst, uint16_t stype, uint16_t rid, void* data, size_t* size);
int ereg_read(uint16_t instance, uint16_t misc, uint16_t resourceid, void* data, size_t size);
int ereg_read_inst_count(uint16_t stype, void* data, size_t* size);
int ereg_read_inst(uint16_t inst, uint16_t stype, uint16_t rid, void* data, size_t* size);
int ereg_write(uint16_t instance, uint16_t misc, uint16_t resourceid, void* data,size_t size);
int ereg_write_inst(uint16_t inst, uint16_t stype, uint16_t rid, void* data, size_t* size);
int ereg_handle_alarm(void* ctx, void *data, size_t *size);


#endif /* INC_EREG_H_ */
