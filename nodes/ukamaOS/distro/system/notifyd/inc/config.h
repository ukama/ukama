/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

#define MAX_LINE_LENGTH 128
#define MAX_ENTRIES     128

typedef struct {

    char serviceName[MAX_LINE_LENGTH];
    char moduleName[MAX_LINE_LENGTH];
    char propertyName[MAX_LINE_LENGTH];
    char type[MAX_LINE_LENGTH];
    char severity[MAX_LINE_LENGTH];
    int  code;
} Entry;

/* Service configuration */
typedef struct {

    char  *serviceName;
    int   servicePort;
    char  *nodedHost;
    int   nodedPort;
    char  *nodedEP;
    char  *remoteServer;
    char  *nodeID;
    int   numEntries;
    Entry entries[MAX_ENTRIES];
} Config;

#endif /* CONFIG_H_ */
