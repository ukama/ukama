/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_CONFIG_H_
#define INC_CONFIG_H_

#ifdef __cplusplus
extern "C" {
#endif

/* Service configuration */
typedef struct {
  char* name;
  int port;
  char* nodedHost;
  int nodedPort;
  char* nodedEP;
  char* remoteServer;
} Config;

#ifdef __cplusplus
}
#endif

#endif /* INC_CONFIG_H_ */
