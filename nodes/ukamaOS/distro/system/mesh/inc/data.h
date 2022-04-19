/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef MESH_DATA_H
#define MESH_DATA_H

#include "mesh.h"
#include "config.h"

void clear_request(MRequest **data);
void handle_recevied_data(MRequest *data, Config *config);

#endif /* MESH_DATA_H */
