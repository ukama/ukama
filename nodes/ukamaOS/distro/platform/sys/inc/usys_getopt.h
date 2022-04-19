/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_GETOPT_H_
#define USYS_GETOPT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "getopt.h"

/**
 * @typedef UsysOption
 *
 * @brief    The LONG_OPTIONS argument to getopt_long or getopt_long_only
 *           is a vector of 'struct option' terminated by an element containing
 *            a name which is zero.
 *
 */
typedef const struct option UsysOption;

/**
 * @fn     int usys_getopt_long(int, char* const*, const char*,
 *         const struct option*, int*)
 * @brief  Decode options from the vector argv (whose length is argc).
 *         The argument shortopts describes the short options to accept,
 *         just as it does in getopt. The argument longopts describes
 *         the long options to accept
 *
 * @param  argc
 * @param  argv
 * @param  shortopts
 * @param  longopts
 * @param  indexptr
 * @return If an option was successfully found,
 *         then it returns the option character.
 *         If all command-line options have been parsed, then it returns -1.
 */
static inline int usys_getopt_long(int argc, char *const *argv, const char *shortopts,
    const struct option *longopts, int *indexptr) {
  return getopt_long(argc, argv, shortopts, longopts, indexptr);
}

#ifdef __cplusplus
}
#endif

#endif /* USYS_GETOPT_H_ */
