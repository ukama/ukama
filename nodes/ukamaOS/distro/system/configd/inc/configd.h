/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_CONFIGD_H_
#define INC_CONFIGD_H_

#include "jserdes.h"
#include "session.h"

int configd_process_incoming_config(const char *service,
		JsonObj *json, Config *config);

int configd_process_complete(const char *service,
		JsonObj *json, Config *config);

int configd_trigger_update(ConfigSession *s);

int configd_read_running_config(ConfigData **c);

void free_config_data(ConfigData *c);
#endif /* INC_NOTIFICATION_H_ */
