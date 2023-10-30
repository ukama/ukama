/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef DEVDB_ATT_DAT31R5A_H_
#define DEVDB_ATT_DAT31R5A_H_

#include "inc/driverfxn.h"
#include "devdb/att/att.h"

int dat31r5a_init ();
int dat31r5a_registration(Device* p_dev);
int dat31r5a_read_prop_count(Device* p_dev, uint16_t * count);
int dat31r5a_read_properties(Device* p_dev, void* prop);
int dat31r5a_configure(void* p_dev, void* prop, void* data );
int dat31r5a_read(void* p_dev, void* prop, void* data);
int dat31r5a_write(void* p_dev, void* prop, void* data);
int dat31r5a_enable(void* p_dev, void* prop, void* data);
int dat31r5a_disable(void* p_dev, void* prop, void* data);

#endif /* DEVDB_ATT_DAT31R5A_H_ */
