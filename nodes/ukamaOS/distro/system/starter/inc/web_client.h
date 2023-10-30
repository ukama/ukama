/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
