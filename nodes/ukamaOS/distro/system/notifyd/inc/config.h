/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
