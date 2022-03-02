/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_file.h"

#include "usys_api.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

#define MAX_STR_LENGTH 64

/* Check if file exist */
int usys_file_path_exist(char *fname) {
    int ret = 0;
    if (usys_access(fname, F_OK) != -1) {
        ret = 1;
    } else {
        usys_log_trace("FILE:: File %s is missing.", fname);
    }
    return ret;
}

int usys_file_exist(char *fname) {
    int ret = 0;
    struct stat sb;

    int fd = usys_file_open(fname, O_RDONLY);
    if (fd > 0) {
        ret = 1;
        usys_stat(fname, &sb);

        ret = S_ISREG(sb.st_mode);
        if (!ret) {
            usys_log_error("Err: FILE:: %s is not a file.", fname);
        }
        usys_file_close(fd);

    }

    return ret;
}

int usys_file_open(char *fname, int flags) {
    int fd = 0;

    /* Create input file descriptor */
    fd = usys_open(fname, flags, 0644);
    if (fd == -1) {
        usys_log_error("Opening file %s failed.", fname);
    }

    return fd;
}

int usys_file_remove(void *data) {
    int ret = -1;
    if (data) {
        char *fname = data;

        ret = usys_remove(fname);
        if (!ret) {
            usys_log_debug("FILE:: %s db file deleted successfully.", fname);
        } else {
            usys_log_debug("Err(%d): FILE:: %s db file deleted successfully.", ret,
                      fname);
        }

    }
    return ret;
}

void usys_file_close(int fd) {
    usys_fsync(fd);
    usys_close(fd);
}

/* if file is a symlink */
int usys_file_symlink_exists(const char *path) {
    struct stat sb;
    int ret = 0;
    if (usys_lstat(path, &sb) == 0) {
        usys_log_trace("FILE:: Symbolic link %s exist.", path);
        ret = 1;
    }
    return (ret);
}

char *usys_file_read_sym_link(char *fname) {
    struct stat sb;
    int readbytes = 0;
    if (usys_lstat(fname, &sb) == -1) {
        usys_log_error("lstat for file %s failed.", fname);
        return NULL;
    }

    char *linkname = usys_malloc(sb.st_size + 1);
    if (linkname) {

        readbytes = usys_readlink(fname, linkname, sb.st_size + 1);
        if (readbytes < 0) {
            usys_log_error("read for file %s failed.", fname);
            usys_free(linkname);
            return NULL;
        }
        if (readbytes > sb.st_size) {
            usys_log_error("Err: FILE: symlink increased in size "
                      "between lstat() and readlink()");
            usys_free(linkname);
            return NULL;
        }

        linkname[sb.st_size] = '\0';
        usys_log_trace("FILE:: '%s' points to '%s'\n", fname, linkname);

    } else {
        return NULL;
    }

    return linkname;
}

/*Used for master db info read.*/
int usys_file_raw_read(char *fname, void *buff, off_t offset, uint16_t size) {
    int read_bytes = 0;
    /* Create input file descriptor */
    int fd = usys_open(fname, O_RDONLY, 0644);
    if (fd == -1) {
        usys_log_error("Opening file %s failed.", fname);
        return fd;
    }

    off_t off = usys_lseek(fd, offset, SEEK_SET);
    if (off < offset) {
        read_bytes = -1;
        return read_bytes;
    }

    read_bytes = usys_read(fd, buff, size);
    return read_bytes;
}

int usys_file_read(void *fname, void *buff, off_t offset, uint16_t size) {
    int read_bytes = 0;

    int fd = usys_file_open(fname, O_RDONLY);
    if (fd < 0) {
        read_bytes = -1;
        return read_bytes;
    }

    off_t off = usys_lseek(fd, offset, SEEK_SET);
    if (off < offset) {
        read_bytes = -1;
        return read_bytes;
    }

    read_bytes = usys_read(fd, buff, size);

    usys_file_close(fd);
    usys_log_trace("FILE:: FD(%d) Read %d bytes from offset 0x%x.", fd, read_bytes,
            offset);
    return read_bytes;
}

int usys_file_write(void *fname, void *buff, off_t offset, uint16_t size) {
    int write_bytes = 0;

    int fd = usys_file_open(fname, O_WRONLY);
    if (fd < 0) {
        write_bytes = -1;
        return write_bytes;
    }

    off_t off = usys_lseek(fd, offset, SEEK_SET);
    if (off < offset) {
        write_bytes = -1;
        return write_bytes;
    }

    write_bytes = write(fd, buff, size);

    usys_file_close(fd);
    usys_log_trace("FILE:: FD(%d) Written %d bytes to offset 0x%x.", fd, write_bytes,
              offset);
    return write_bytes;
}

int usys_file_append(void *fname, void *buff, off_t offset, uint16_t size) {
    int write_bytes = 0;

    int fd = usys_file_open(fname, O_WRONLY);
    if (fd < 0) {
        write_bytes = -1;
        return write_bytes;
    }

    off_t off = usys_lseek(fd, offset, SEEK_END);
    if (off < offset) {
        write_bytes = -1;
        return write_bytes;
    }

    write_bytes = usys_write(fd, buff, size);

    usys_file_close(fd);
    usys_log_trace("FILE:: FD(%d) Written %d bytes to offset 0x%x.", fd, write_bytes,
              offset);
    return write_bytes;
}

int usys_file_erase(void *fname, off_t offset, uint16_t size) {
    int write_bytes = 0;
    int fd = -1;

    char *buff = usys_malloc(sizeof(char) * size);
    if (buff) {
        usys_memset(buff, 0xff, size);
        fd = usys_file_open(fname, O_WRONLY);
        if (fd < 0) {
            write_bytes = -1;
            return write_bytes;
        }
        usys_lseek(fd, offset, SEEK_SET);

        write_bytes = usys_write(fd, buff, size);

        usys_file_close(fd);

        usys_free(buff);
    }

    usys_log_trace("FILE:: Erased bytes: %d from %d", write_bytes, fd);
    return write_bytes;
}

int usys_file_read_number(void *fname, void *data, off_t offset, uint16_t count,
                     uint8_t size) {
    int ret = 0;
    char val[8];
    uint16_t idx = 0;
    char *value = (char *)data;

    while (idx < count) {
        if (usys_file_read(fname, val, offset, size) < size) {
            return -1;
        }

        usys_memcpy((value + (idx * size)), val, size);

        for (int i = 0; i < size; i++) {
            usys_log_trace("\t \t File[%d] = 0x%x.", offset,
                      (uint8_t) * (value + (idx * size) + i));
        }

        offset = offset + size;
        idx++;
    }

    return ret;
}

int usys_file_write_number(void *fname, void *data, off_t offset, uint16_t count,
                      uint8_t size) {
    int ret = 0;
    uint16_t idx = 0;
    char val[8];
    char *value = (char *)data;

    while (idx < count) {
        usys_memcpy(val, value + (idx * size), size);

        if (usys_file_write(fname, val, offset, size) < size) {
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

int usys_file_protect(void *fname) {
    //dummy
    return 0;
}

int usys_file_init(void *data) {
    char fname[MAX_STR_LENGTH] = { '\0' };
    int size = usys_strlen((char *)data);
    usys_memcpy(fname, data, size);

    int fd = usys_file_open(fname, O_RDONLY);
    if (fd < 0) {
        /* This means db doesn't exist.*/
        usys_log_warn("FILE:: %s doesn't exist.So creating it", fname);

        fd = usys_file_open(fname, (O_WRONLY | O_CREAT));
        if (fd < 0) {
            return -1;
        }

    }

    usys_file_close(fd);
    usys_log_debug("FILE::File %s is ready.", fname);
    return 0;
}

int usys_file_cleanup(void *fname) {
    int ret = 0;

    ret = usys_remove(fname);
    if (!ret) {
        usys_log_debug("FILE:: DB %s deleted successfully.", fname);
    } else {
        usys_log_debug("FILE:: DB %s deletion failed.", fname);
    }

    return ret;
}

int usys_file_rename(char *old_name, char *new_name) {
    int ret = 0;
    if (usys_rename(old_name, new_name) == 0) {
        usys_log_debug("FILE:: DB %s renamed to %s.", old_name, new_name);
    } else {
        ret = -1;
        usys_log_error("Err:: Unable to rename file %s to %s.", old_name, new_name);
    }
    return ret;
}

int usys_file_add_record(char *filename, char *rowdesc, char *data) {
    int ret = 0;
    /* Check if we need to create a new file */
    if (!usys_file_exist(filename)) {

        ret = usys_file_init(filename);
        if (ret) {
            return ret;
        }
        /* Add column description */
        ret = usys_file_append(filename, rowdesc, 0, usys_strlen(rowdesc));
    }
    /* Add data to file */
    ret = usys_file_append(filename, data, 0, usys_strlen(data));
    return ret;
}
