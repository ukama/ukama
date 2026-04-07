/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include "file.h"

#include "usys_log.h"

int file_exist(char *fname) {

    int ret        = 0;
    int fd         = 0;
    struct stat sb;

    fd = file_open(fname, O_RDONLY);
    if (fd <= 0) {
        return 0;
    }

    stat(fname, &sb);
    ret = S_ISREG(sb.st_mode);
    if (ret == 0) {
        usys_log_error("%s is not a regular file", fname);
        file_close(fd);
        return 0;
    }

    file_close(fd);

    return 1;
}

int file_open(char *fname, int flags) {

    int fd = 0;

    fd = open(fname, flags, 0644);
    if (fd == -1) {
        perror("open");
    }

    return fd;
}

int file_remove(void *data) {

    int ret      = -1;
    char *fname  = NULL;

    if (data == NULL) {
        return ret;
    }

    fname = data;
    ret = remove(fname);
    if (ret == 0) {
        usys_log_debug("%s deleted successfully", fname);
    } else {
        usys_log_debug("failed to delete %s", fname);
    }

    return ret;
}

void file_close(int fd){

    fsync(fd);
    close(fd);
}
