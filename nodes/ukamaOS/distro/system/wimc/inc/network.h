/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WIMC_NETWORK_H
#define WIMC_NETWORK_H

int init_frameworks(struct _u_instance *adminInst,
		    struct _u_instance *clientInst, int adminPort,
		    int clientPort);
void setup_admin_endpoints(WimcCfg *cfg, struct _u_instance *instance);
void setup_client_endpoints(WimcCfg *cfg, struct _u_instance *instance);
int start_framework(struct _u_instance *instance);
int start_web_services(WimcCfg *cfg, struct _u_instance *adminInst,
		       struct _u_instance *clientInst);

#endif /* WIMC_NETWORK_H */
