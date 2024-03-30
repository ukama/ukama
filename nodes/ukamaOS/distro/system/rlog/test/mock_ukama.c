/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <jansson.h>
#include <ulfius.h>

#include "http_status.h"
#include "usys_services.h"
#include "usys_file.h"

#include "rlogd.h"

#define EP_SERVICE       "v1/nodes/logger/:nodeID"

typedef struct _u_request  req_t;
typedef struct _u_response resp_t;
typedef struct _u_map      map_t;

void print_log(const struct _u_request *request) {

    fprintf(stdout, "Recevied log message: \n");
    fprintf(stdout, " Length: %ld \n Data: %s \n",
            request->binary_body_length, (char *)request->binary_body);
}

int callback_default(const URequest *request, UResponse *response,
                     void *data) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_Unauthorized,
                                    HttpStatusStr(HttpStatus_Unauthorized));

    return U_CALLBACK_CONTINUE;
}

int callback_post_log(const URequest *request, UResponse *response,
                      void *data) {

    print_log(request);
    
    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

int main(int argc, char **argv) {

	struct _u_instance inst;
	int port = 0;

    port = usys_find_service_port(SERVICE_REMOTE);
    if (port == 0) {
        fprintf(stderr, "Unable to find remote service port");
        exit(1);
    }
	
	if (ulfius_init_instance(&inst, port, NULL, NULL) != U_OK) {
		fprintf(stderr, "Error ulfius_init_instance, abort\n");
		return 1;
	}

	/* Endpoint list declaration */
	ulfius_add_endpoint_by_val(&inst, "POST", EP_SERVICE, NULL, 0,
							   &callback_post_log, NULL);
    ulfius_set_default_endpoint(&inst, &callback_default, NULL);
    
	/* Start the framework */
	if (ulfius_start_framework(&inst) == U_OK) {
		fprintf(stdout, "Famework start on port %d\n", inst.port);
	} else {
		fprintf(stderr, "Error starting framework\n");
        exit(1);
	}

    pause();

	ulfius_stop_framework(&inst);
	ulfius_clean_instance(&inst);
	
	return 0;
}
