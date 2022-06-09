/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_CONFIG_H_
#define INC_CONFIG_H_

/* Service configuration */
typedef struct {
  char* name;
  int port;
  char* nodedHost;
  int nodedPort;
  char* nodedEP;
  char* remoteServer;
} Config;


#endif /* INC_CONFIG_H_ */
