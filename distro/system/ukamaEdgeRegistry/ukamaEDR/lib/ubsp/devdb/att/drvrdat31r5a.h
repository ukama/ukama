/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
