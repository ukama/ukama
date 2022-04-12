/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef GLOBALHEADER_H_
#define GLOBALHEADER_H_

#define MAX_NAME_LENGTH                 24
#define MAX_PATH_LENGTH                 64

#define MODULE_UUID_LENGTH          NAME_LENGTH
#define MODULE_UUID_LENGTH          NAME_LENGTH
#define SYSFS_PATH_LENGTH           PATH_LENGTH

#define UBSP_FREE(mem) \
    if (mem) {         \
        free(mem);     \
        mem = NULL;    \
    }

#endif /* GLOBALHEADER_H_ */
