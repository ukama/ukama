/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DEVDB_ATT_DRVRDAT31R5A_H_
#define DEVDB_ATT_DRVRDAT31R5A_H_

#include "devdb/att/dat31r5a.h"

int drvr_dat31r5a_init ();
int drvr_dat31r5a_registration(Device* p_dev);
int drvr_dat31r5a_read_properties(DevObj* obj, void* prop, uint16_t* count);
int drvr_dat31r5a_configure(void* p_dev, void* prop, void* data );
int drvr_dat31r5a_read(void* p_dev, void* prop, void* data);
int drvr_dat31r5a_write(void* p_dev, void* prop, void* data);
int drvr_dat31r5a_enable(void* p_dev, void* prop, void* data);
int drvr_dat31r5a_disable(void* p_dev, void* prop, void* data);

#endif /* DEVDB_ATT_DRVRDAT31R5A_H_ */
