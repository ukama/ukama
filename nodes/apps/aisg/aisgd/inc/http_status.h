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
    HttpStatus_OK                  = 200,
    HttpStatus_BadRequest          = 400,
    HttpStatus_NotFound            = 404,
    HttpStatus_MethodNotAllowed    = 405,
    HttpStatus_InternalServerError = 500,
    HttpStatus_ServiceUnavailable  = 503
};

static const char *HttpStatusStr(int code) {

    switch (code) {
    case 200: return "OK";
    case 400: return "Bad Request";
    case 404: return "Not Found";
    case 405: return "Method Not Allowed";
    case 500: return "Internal Server Error";
    case 503: return "Service Unavailable";
    default:  return "Unknown";
    }
}

#endif /* HTTP_STATUS_H_ */
