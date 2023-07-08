/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef HTTP_STATUS_H
#define HTTP_STATUS_H

/* Source: https://en.wikipedia.org/wiki/List_of_HTTP_status_codes */

/*
 * Enum for the HTTP status codes
 */
enum HttpStatusCode {

	/* 2xx success */
	HttpStatus_OK                          = 200,
	HttpStatus_Created                     = 201,
	HttpStatus_Accepted                    = 202,
	HttpStatus_NonAuthoritativeInformation = 203,
	HttpStatus_NoContent                   = 204,
	HttpStatus_ResetContent                = 205,
	HttpStatus_PartialContent              = 206,
	HttpStatus_MultiStatus                 = 207,
	HttpStatus_AlreadyReported             = 208,
	HttpStatus_IMUsed                      = 226,

	/* 3xx re-direction */
	HttpStatus_MultipleChoices             = 300,
	HttpStatus_MovedPermanently            = 301,
	HttpStatus_Found                       = 302,
	HttpStatus_SeeOther                    = 303,
	HttpStatus_NotModified                 = 304,
	HttpStatus_UseProxy                    = 305,
	HttpStatus_TemporaryRedirect           = 307,
	HttpStatus_PermanentRedirect           = 308,

	/* 4xx client error */
	HttpStatus_BadRequest                  = 400,
	HttpStatus_Unauthorized                = 401,
	HttpStatus_PaymentRequired             = 402,
	HttpStatus_Forbidden                   = 403,
	HttpStatus_NotFound                    = 404,
	HttpStatus_MethodNotAllowed            = 405,
	HttpStatus_NotAcceptable               = 406,
	HttpStatus_ProxyAuthenticationRequired = 407,
	HttpStatus_RequestTimeout              = 408,
	HttpStatus_Conflict                    = 409,
	HttpStatus_Gone                        = 410,
	HttpStatus_LengthRequired              = 411,
	HttpStatus_PreconditionFailed          = 412,
	HttpStatus_ContentTooLarge             = 413,
	HttpStatus_PayloadTooLarge             = 413,
	HttpStatus_URITooLong                  = 414,
	HttpStatus_UnsupportedMediaType        = 415,
	HttpStatus_RangeNotSatisfiable         = 416,
	HttpStatus_ExpectationFailed           = 417,
	HttpStatus_ImATeapot                   = 418,
	HttpStatus_MisdirectedRequest          = 421,
	HttpStatus_UnprocessableContent        = 422,
	HttpStatus_UnprocessableEntity         = 422,
	HttpStatus_Locked                      = 423,
	HttpStatus_FailedDependency            = 424,
	HttpStatus_TooEarly                    = 425,
	HttpStatus_UpgradeRequired             = 426,
	HttpStatus_PreconditionRequired        = 428,
	HttpStatus_TooManyRequests             = 429,
	HttpStatus_RequestHeaderFieldsTooLarge = 431,
	HttpStatus_UnavailableForLegalReasons  = 451,

	/* 5xx Server error */
	HttpStatus_InternalServerError         = 500,
	HttpStatus_NotImplemented              = 501,
	HttpStatus_BadGateway                  = 502,
	HttpStatus_ServiceUnavailable          = 503,
	HttpStatus_GatewayTimeout              = 504,
	HttpStatus_HTTPVersionNotSupported     = 505,
	HttpStatus_VariantAlsoNegotiates       = 506,
	HttpStatus_InsufficientStorage         = 507,
	HttpStatus_LoopDetected                = 508,
	HttpStatus_NotExtended                 = 510,
	HttpStatus_NetworkAuthenticationRequired = 511,
};

static const char *HttpStatusStr(int code) {

	switch (code) {

	case 200: return "OK";
	case 201: return "Created";
	case 202: return "Accepted";
	case 203: return "Non-Authoritative Information";
	case 204: return "No Content";
	case 205: return "Reset Content";
	case 206: return "Partial Content";
	case 207: return "Multi-Status";
	case 208: return "Already Reported";
	case 226: return "IM Used";

	case 300: return "Multiple Choices";
	case 301: return "Moved Permanently";
	case 302: return "Found";
	case 303: return "See Other";
	case 304: return "Not Modified";
	case 305: return "Use Proxy";
	case 307: return "Temporary Redirect";
	case 308: return "Permanent Redirect";

	case 400: return "Bad Request";
	case 401: return "Unauthorized";
	case 402: return "Payment Required";
	case 403: return "Forbidden";
	case 404: return "Not Found";
	case 405: return "Method Not Allowed";
	case 406: return "Not Acceptable";
	case 407: return "Proxy Authentication Required";
	case 408: return "Request Timeout";
	case 409: return "Conflict";
	case 410: return "Gone";
	case 411: return "Length Required";
	case 412: return "Precondition Failed";
	case 413: return "Content Too Large";
	case 414: return "URI Too Long";
	case 415: return "Unsupported Media Type";
	case 416: return "Range Not Satisfiable";
	case 417: return "Expectation Failed";
	case 418: return "I'm a teapot";
	case 421: return "Misdirected Request";
	case 422: return "Unprocessable Content";
	case 423: return "Locked";
	case 424: return "Failed Dependency";
	case 425: return "Too Early";
	case 426: return "Upgrade Required";
	case 428: return "Precondition Required";
	case 429: return "Too Many Requests";
	case 431: return "Request Header Fields Too Large";
	case 451: return "Unavailable For Legal Reasons";

	case 500: return "Internal Server Error";
	case 501: return "Not Implemented";
	case 502: return "Bad Gateway";
	case 503: return "Service Unavailable";
	case 504: return "Gateway Timeout";
	case 505: return "HTTP Version Not Supported";
	case 506: return "Variant Also Negotiates";
	case 507: return "Insufficient Storage";
	case 508: return "Loop Detected";
	case 510: return "Not Extended";
	case 511: return "Network Authentication Required";

	default: return "";
	}
}

#endif /* HTTP_STATUS_H */
