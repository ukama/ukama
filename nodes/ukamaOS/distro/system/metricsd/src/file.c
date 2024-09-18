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

bool file_validate_offset(off_t offset) {
  bool ret = false;
  if ((offset <= MAX_UKDB_OFFSET) && (offset >= MIN_UKDB_OFFSET)) {
    ret = true;
  }
  return ret;
}

/* Check if file exist */
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

/* if file is a symlink */
int file_symlink_exists(const char *path) {
  struct stat sb;
  int ret = 0;
  if (lstat(path, &sb) == 0) {
    usys_log_trace("FILE:: Symbolic link %s exist.", path);
    ret = 1;
  }
  return (ret);
}

char *file_read_sym_link(char *fname) {
  struct stat sb;
  int readbytes = 0;
  if (lstat(fname, &sb) == -1) {
    perror("lstat");
    return NULL;
  }
  char *linkname = malloc(sb.st_size + 1);
  if (linkname) {
    readbytes = readlink(fname, linkname, sb.st_size + 1);
    if (readbytes < 0) {
      perror("lstat");
      free(linkname);
      return NULL;
    }
    if (readbytes > sb.st_size) {
      usys_log_error("Symlink increased in size between lstat() and readlink()");
      free(linkname);
      return NULL;
    }
    linkname[sb.st_size] = '\0';

    usys_log_trace("'%s' points to '%s'\n", fname, linkname);
  } else {
    return NULL;
  }
  return linkname;
}

/*Used for master db info read.*/
int file_raw_read(char *fname, void *buff, off_t offset, uint16_t size) {
  int read_bytes = 0;
  /* Create input file descriptor */
  int fd = open(fname, O_RDONLY, 0644);
  if (fd == -1) {
    perror("open");
    return fd;
  }
  off_t off = lseek(fd, offset, SEEK_SET);
  if (off < offset) {
    read_bytes = -1;
    return read_bytes;
  }
  read_bytes = read(fd, buff, size);
  return read_bytes;
}

int file_read(void *fname, void *buff, off_t offset, uint16_t size) {
  int read_bytes = 0;
  int fd = file_open(fname, O_RDONLY);
  if (fd < 0) {
    read_bytes = -1;
    return read_bytes;
  }

  off_t off = lseek(fd, offset, SEEK_SET);
  if (off < offset) {
    read_bytes = -1;
    return read_bytes;
  }

  if (file_validate_offset(offset)) {
    read_bytes = read(fd, buff, size);
  }
  file_close(fd);
  usys_log_trace("FD(%d) Read %d bytes from offset 0x%x.", fd, read_bytes,
                 offset);
  return read_bytes;
}

int file_write(void *fname, void *buff, off_t offset, uint16_t size) {
  int write_bytes = 0;
  int fd = file_open(fname, O_WRONLY);
  if (fd < 0) {
    write_bytes = -1;
    return write_bytes;
  }
  off_t off = lseek(fd, offset, SEEK_SET);
  if (off < offset) {
    write_bytes = -1;
    return write_bytes;
  }
  if (file_validate_offset(offset)) {
    write_bytes = write(fd, buff, size);
  }
  file_close(fd);
  usys_log_trace("FD(%d) Written %d bytes to offset 0x%x.", fd, write_bytes,
                 offset);
  return write_bytes;
}

int file_append(void *fname, void *buff, off_t offset, uint16_t size) {
  int write_bytes = 0;
  int fd = file_open(fname, O_WRONLY);
  if (fd < 0) {
    write_bytes = -1;
    return write_bytes;
  }
  off_t off = lseek(fd, offset, SEEK_END);
  if (off < offset) {
    write_bytes = -1;
    return write_bytes;
  }
  if (file_validate_offset(offset)) {
    write_bytes = write(fd, buff, size);
  }
  file_close(fd);
  usys_log_trace("FD(%d) Written %d bytes to offset 0x%x.", fd, write_bytes,
                 offset);
  return write_bytes;
}

int file_erase(void *fname, off_t offset, uint16_t size) {
  int write_bytes = 0;
  int fd = -1;
  char *buff = malloc(sizeof(char) * size);
  if (buff) {
    memset(buff, 0xff, size);
    fd = file_open(fname, O_WRONLY);
    if (fd < 0) {
      write_bytes = -1;
      return write_bytes;
    }
    lseek(fd, offset, SEEK_SET);
    if (file_validate_offset(offset)) {
      write_bytes = write(fd, buff, size);
    }
    file_close(fd);
  }
  if (buff) {
    free(buff);
  }
  usys_log_trace("Erased bytes: %d from %d", write_bytes, fd);
  return write_bytes;
}

int file_read_number(void *fname, void *data, off_t offset, uint16_t count,
                     uint8_t size) {
  int ret = 0;
  char val[8];
  uint16_t idx = 0;
  char *value = (char *)data;
  while (idx < count) {
    if (file_read(fname, val, offset, size) < size) {
      return -1;
    }
    memcpy((value + (idx * size)), val, size);
    for (int i = 0; i < size; i++) {
      usys_log_trace("\t \t File[%d] = 0x%x.", offset,
                     (uint8_t) * (value + (idx * size) + i));
    }
    offset = offset + size;
    idx++;
  }
  return ret;
}

int file_write_number(void *fname, void *data, off_t offset, uint16_t count,
                      uint8_t size) {
  int ret = 0;
  uint16_t idx = 0;
  char val[8];
  char *value = (char *)data;
  while (idx < count) {
    memcpy(val, value + (idx * size), size);
    if (file_write(fname, val, offset, size) < size) {
      return -1;
    }
    for (int i = 0; i < size; i++) {
      usys_log_trace("\t \t File[%d] = 0x%x.", offset,
                     (uint8_t) * (value + (idx * size) + i));
    }
    offset = offset + size;
    idx++;
  }
  return ret;
}

int file_protect(void *fname) {
  // dummy
  return 0;
}

int file_init(void *data) {
  char fname[MAX_STR_LENGTH] = {'\0'};
  int size = strlen((char *)data);
  memcpy(fname, data, size);
  int fd = file_open(fname, O_RDONLY);
  if (fd < 0) {
    /* This means db doesn't exist.*/
    usys_log_debug("%s doesn't exist.So creating it", fname);
    fd = file_open(fname, (O_WRONLY | O_CREAT));
    if (fd < 0) {
      return -1;
    }
  }
  file_close(fd);
  return 0;
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

int file_add_record(char *filename, char *rowdesc, char *data) {
  int ret = 0;
  /* Check if we need to create a new file */
  if (!file_exist(filename)) {
    ret = file_init(filename);
    if (ret) {
      return ret;
    }
    /* Add column description */
    file_append(filename, rowdesc, 0, strlen(rowdesc));
  }
  /* Add data to file */
  file_append(filename, data, 0, strlen(data));
  return ret;
}
