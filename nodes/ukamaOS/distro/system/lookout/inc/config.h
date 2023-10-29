/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

/* Service configuration */
typedef struct {

    int   servicePort;
    int   nodedPort;
    int   starterdPort;
    int   nodeSystemPort;
    char  *nodeID;
} Config;

#endif /* CONFIG_H_ */
