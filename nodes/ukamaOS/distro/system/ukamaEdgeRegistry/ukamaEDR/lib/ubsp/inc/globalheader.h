/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
