/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */


#ifndef WEB_CLIENT_H
#define WEB_CLIENT_H

#define MAX_URL_LEN  256

#define WIMC_EP "/content/containers"

#define WIMC_RESP_TYPE_RESULT     "result"
#define WIMC_RESP_TYPE_ERROR      "error"
#define WIMC_RESP_TYPE_PROCESSING "processing"

/* For JSON de-serialization */
#define JSON_TYPE            "type"
#define JSON_TYPE_RESULT     "type_result"
#define JSON_VOID_STR        "void"

#define JSON_WIMC_RESPONSE   "wimc_response"

int get_capp_path(Config *config, char *name, char *tag,
                  char **path, int *retCode);
#endif /* WIMC_H */
