/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef ERROR_H
#define ERROR_H

#define WIMC_OK                    0
#define WIMC_ERROR_EXIST           1
#define WIMC_ERROR_BAD_NAME        2
#define WIMC_ERROR_BAD_ACTION      3
#define WIMC_ERROR_BAD_TYPE        4
#define WIMC_ERROR_BAD_METHOD      5
#define WIMC_ERROR_BAD_URL         6
#define WIMC_ERROR_MEMORY          7
#define WIMC_ERROR_BAD_ID          8
#define WIMC_ERROR_BAD_INTERVAL    9
#define WIMC_ERROR_MISSING_CONTENT 10
#define WIMC_ERROR_MAX_AGENTS      11

#define WIMC_OK_STR                  "OK"
#define WIMC_ERROR_EXIST_STR         "Already Registered"
#define WIMC_ERROR_BAD_NAME_STR      "Invalid/bad name/tag"
#define WIMC_ERROR_BAD_ACTION_STR    "Invalid Action"
#define WIMC_ERROR_BAD_TYPE_STR      "Invalid type"
#define WIMC_ERROR_BAD_ID_STR        "Invalid ID"         
#define WIMC_ERROR_BAD_INTERVAL_STR  "Invalid Interval"
#define WIMC_ERROR_BAD_METHOD_STR    "Invalid method"
#define WIMC_ERROR_BAD_URL_STR       "Invalid URL"
#define WIMC_ERROR_MEMORY_STR        "Internal memory error"
#define WIMC_ERROR_MISSING_CONTENT_STR "Missing Content"
#define WIMC_ERROR_MAX_AGENTS_STR    "Max Agents reached"

extern const char *error_to_str(int error);
#endif /* ERROR_H */
