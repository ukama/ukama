/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef AGENT_CALLBACK_H
#define AGENT_CALLBACK_H

int agent_callback_get(const struct _u_request *request,
		       struct _u_response *response, void *user_data);
int agent_callback_post(const struct _u_request *request,
			struct _u_response *response, void *user_data);
int agent_callback_put(const struct _u_request *request,
		       struct _u_response *response, void *user_data);
int agent_callback_delete(const struct _u_request *request,
			  struct _u_response *response, void *user_data);
int agent_callback_stats(const struct _u_request *request,
			 struct _u_response *response, void *user_data);
int agent_callback_default(const struct _u_request *request,
			   struct _u_response *response, void *user_data);

#endif /* AGENT_CALLBACK_H */
