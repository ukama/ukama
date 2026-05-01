/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef HTTP_STATUS_H_
#define HTTP_STATUS_H_

enum HttpStatusCode {

    HttpStatus_OK               = 200,
    HttpStatus_NotFound         = 404,
    HttpStatus_MethodNotAllowed = 405,
};

static inline const char *HttpStatusStr(int code) {

    switch (code) {
    case HttpStatus_OK:
        return "OK";

    case HttpStatus_NotFound:
        return "Not Found";

    case HttpStatus_MethodNotAllowed:
        return "Method Not Allowed";

    default:
        return "";
    }
}

#endif /* HTTP_STATUS_H_ */
