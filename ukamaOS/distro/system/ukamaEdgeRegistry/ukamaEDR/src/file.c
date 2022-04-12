/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "headers/utils/file.h"
#include "dmt.h"

#include "headers/utils/log.h"

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
        log_trace("FILE:: File %s is missing.", fname);
    }
    return ret;
}

int file_exist(char *fname) {
    int ret = 0;
    struct stat sb;
    int fd = file_open(fname, O_RDONLY);
    if (fd > 0) {
        ret = 1;
        stat(fname, &sb);
        ret = S_ISREG(sb.st_mode);
        if (!ret) {
            log_error("Err: FILE:: %s is not a file.", fname);
        }
        file_close(fd);
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
            log_debug("FILE:: %s db file deleted successfully.", fname);
        } else {
            log_debug("Err(%d): FILE:: %s db file deleted successfully.", ret,
                      fname);
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
        log_trace("FILE:: Symbolic link %s exist.", path);
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
    char *linkname = dmt_malloc(sb.st_size + 1);
    if (linkname) {
        readbytes = readlink(fname, linkname, sb.st_size + 1);
        if (readbytes < 0) {
            perror("lstat");
            dmt_free(linkname);
            return NULL;
        }
        if (readbytes > sb.st_size) {
            log_error("Err: FILE: symlink increased in size "
                      "between lstat() and readlink()");
            dmt_free(linkname);
            return NULL;
        }
        linkname[sb.st_size] = '\0';

        log_trace("FILE:: '%s' points to '%s'\n", fname, linkname);
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
    log_trace("FILE:: FD(%d) Read %d bytes from offset 0x%x.", fd, read_bytes,
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
    log_trace("FILE:: FD(%d) Written %d bytes to offset 0x%x.", fd, write_bytes,
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
    log_trace("FILE:: FD(%d) Written %d bytes to offset 0x%x.", fd, write_bytes,
              offset);
    return write_bytes;
}

int file_erase(void *fname, off_t offset, uint16_t size) {
    int write_bytes = 0;
    int fd = -1;
    char *buff = dmt_malloc(sizeof(char) * size);
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
        dmt_free(buff);
    }
    log_trace("FILE:: Erased bytes: %d from %d", write_bytes, fd);
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
            log_trace("\t \t File[%d] = 0x%x.", offset,
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
            log_trace("\t \t File[%d] = 0x%x.", offset,
                      (uint8_t) * (value + (idx * size) + i));
        }
        offset = offset + size;
        idx++;
    }
    return ret;
}

int file_protect(void *fname) {
    //dummy
    return 0;
}

int file_init(void *data) {
    char fname[MAX_STR_LENGTH] = { '\0' };
    int size = strlen((char *)data);
    memcpy(fname, data, size);
    int fd = file_open(fname, O_RDONLY);
    if (fd < 0) {
        /* This means db doesn't exist.*/
        log_debug("FILE:: %s doesn't exist.So creating it", fname);
        fd = file_open(fname, (O_WRONLY | O_CREAT));
        if (fd < 0) {
            return -1;
        }
    }
    file_close(fd);
    log_debug("FILE::File %s is ready.", fname);
    return 0;
}

int file_cleanup(void *fname) {
    int ret = 0;
    ret = remove(fname);
    if (!ret) {
        log_debug("FILE:: DB %s deleted successfully.", fname);
    } else {
        log_debug("FILE:: DB %s deletion failed.", fname);
    }
    return ret;
}

int file_rename(char *old_name, char *new_name) {
    int ret = 0;
    if (rename(old_name, new_name) == 0) {
        log_debug("FILE:: DB %s renamed to %s.", old_name, new_name);
    } else {
        ret = -1;
        log_error("Err:: Unable to rename file %s to %s.", old_name, new_name);
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

#if 0
int test_case1() {
    int ret = -1;
    char buff[32] = "MY NAME IS VISHAL";
    char testbuff[32] = { '\0' };
    int size = strlen(buff);
    file_write(buff, 0, size);
    file_read(testbuff, 0, size);
    printf("Read data is %s.\n", testbuff);
    if (strcmp(buff, testbuff)) {
        ret = 0;
    }
    return ret;
}

int test_case2() {
    int ret = -1;
    int type = 12;
    int test_type;
    char ty1[4];
    char ty[4];
    memcpy(ty, &type, 4);
    printf("Write: %x %x %x %x\n", ty[0], ty[1], ty[2], ty[3]);
    file_write(ty, 18, 4);
    file_read(ty1, 18, 4);
    printf("Read: %x %x %x %x \n", ty1[0], ty[1], ty[2], ty[3]);
    memcpy(&test_type, ty1, 4);
    printf("Read data is %s and %d.", ty1, test_type);
    if (strncmp(ty1, ty, 4)) {
        ret = 0;
    }
    return ret;
}

int test_case3() {
    int ret = -1;
    uint16_t write[3] = { 455, 35, 6335 };
    uint16_t read[3] = { 0 };
    file_write_number(write, 22, 3, sizeof(int));
    file_read_number(read, 22, 3, sizeof(int));
    printf("Write %d %d %d  and read %d %d %d.\n", write[0], write[1], write[2],
           read[0], read[1], read[2]);
    if ((read[0] == write[0]) && (read[1] == write[1]) &&
        (read[2] == write[2])) {
        ret = 0;
    }
    return ret;
}
#endif
#if FILE_TEST
int main(int argc, char *argv[]) {
    test_case1();
    test_case2();
    test_case3();
    test_case4();

    return (EXIT_SUCCESS);
}
#endif
