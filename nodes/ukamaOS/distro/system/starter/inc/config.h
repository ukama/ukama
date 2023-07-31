/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

#include "usys_types.h"

/* Service configuration */
typedef struct {

    int   servicePort;
    int   nodedPort;
    int   notifydPort;
    int   wimcPort;
    char  *nodeID;
    char  *manifestFile;
} Config;

#endif /* CONFIG_H_ */
