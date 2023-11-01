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

/* Service configuration */
typedef struct {
	char  *serviceName;
	int   servicePort;
	char* nodedHost;
	int  nodedPort;
	char  *nodedEP;
	char* starterHost;
	int  starterPort;
	char  *starterEP;
	char *nodeId;
	void *updateSession;
	void *runningConfig;
} Config;

#endif /* CONFIG_H_ */
