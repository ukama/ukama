/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef WEB_CLIENT_H_
#define WEB_CLIENT_H_

#include "config.h"
#include "lookout.h"

void add_capp_to_list(CappList **list,
                      const char *space,
                      const char *name,
                      const char *tag,
                      const char *status,
                      int pid);
int get_nodeid_from_noded(Config *config);
int send_health_report(Config *config);

#endif /* INC_WEB_CLIENT_H_ */
