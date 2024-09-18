/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include "file.h"

#include "usys_log.h"

#define MAX_STR_LENGTH 64

int file_path_exist(char *fname) {
  int ret = 0;
  if (access(fname, F_OK) != -1) {
    ret = 1;
  } else {
    usys_log_trace("File %s is missing.", fname);
  }
  return ret;
}

int file_exist(char *fname) {
  int ret = 0;
  struct stat sb;
  int fd = file_open(fname, O_RDONLY);
  if (fd > 0) {
    stat(fname, &sb);
    ret = S_ISREG(sb.st_mode);
    if (!ret) {
      usys_log_error("Error %s is not a file.", fname);
      return ret;
    }
    file_close(fd);
    ret = 1;
  }
  return ret;
}

int file_open(char *fname, int flags) {
  int fd = 0;
  /* Create input file descriptor */
  fd = open(fname, flags, 0644);
  if (fd == -1) {
    perror("open");
  }
  return fd;
}

int file_remove(void *data) {
  int ret = -1;
  if (data) {
    char *fname = data;
    ret = remove(fname);
    if (!ret) {
      usys_log_debug("%s db file deleted successfully.", fname);
    } else {
      usys_log_debug("%s db file deleted successfully.", ret, fname);
    }
  }
  return ret;
}

void file_close(int fd) {
  fsync(fd);
  close(fd);
}

int file_cleanup(void *fname) {
  int ret = 0;
  ret = remove(fname);
  if (!ret) {
    usys_log_debug("DB %s deleted successfully.", fname);
  } else {
    usys_log_debug("DB %s deletion failed.", fname);
  }
  return ret;
}

int file_rename(char *old_name, char *new_name) {
  int ret = 0;
  if (rename(old_name, new_name) == 0) {
    usys_log_debug("DB %s renamed to %s.", old_name, new_name);
  } else {
    ret = -1;
    usys_log_error("Unable to rename file %s to %s.", old_name, new_name);
  }
  return ret;
}
