/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "drivers/sysfs.h"
#include "errorcode.h"

#include "usys_api.h"
#include "usys_mem.h"
#include "usys_log.h"
#include "usys_string.h"

/*
 * TODO: Use either of the file file.c or sysfs.c.
 * Decision pending will be decided after having a look
 * and type of operations we will be doing on the sys file.
 */

int sysfs_init(char *name, void *data) {
    return 0;
}

/* Open sysfs file */
int sysfs_open(char *name, int flags) {
    int fd = -1;
    /* Create input file descriptor */
    fd = usys_open(name, flags, 0644);
    if (fd == -1) {
        usys_log_error("Opening file %s failed. Error: %s", usys_error(errno));
    }
    return fd;
}

/* Close a sysfs file.*/
void sysfs_close(int fd) {
    usys_fsync(fd);
    usys_close(fd);
}

/* Check if sysfile exist */
int sysfs_exist(char *name) {
    int ret = 0;
    if (usys_access(name, F_OK) != -1) {
        ret = 1;
    } else {
        usys_log_trace("SYSFS:: File %s is missing.", name);
    }
    return ret;
}

/*Erase file content*/
int sysfs_erase(char *name, uint16_t size) {
    int write_bytes = -1;
    char *buff = usys_malloc(sizeof(char) * size);
    if (buff) {
        usys_memset(buff, 0xff, size);
        int fd = sysfs_open(name, O_WRONLY);
        if (fd < 0) {
            write_bytes = -1;
        } else {
            usys_lseek(fd, SYS_DEF_OFFSET, SEEK_SET);
            write_bytes = usys_write(fd, buff, size);
        }
        sysfs_close(fd);
        if (buff) {
            usys_free(buff);
            buff = NULL;
        }
    }
    usys_log_trace("SYSFS:: Erased bytes: %d from file %s", write_bytes, name);
    return write_bytes;
}

/* Sysfs always store value as char so type casting it to proper type*/
void sysfs_strtotype(void *data, int type, char *val) {
    switch (type) {
    case TYPE_NULL: {
        break;
    }
    case TYPE_CHAR: {
        char *rval = data;
        *rval = usys_strtol(val, NULL, 10);
        usys_log_trace("SYSFS:: TYPE_CHAR Value Read is %c.", *rval);
        break;
    }
    case TYPE_BOOL: {
        bool *rval = data;
        *rval = (usys_atoi(val) ? 1 : 0);
        usys_log_trace("SYSFS:: TYPE_BOOL Value Read is %d.", *rval);
        break;
    }
    case TYPE_UINT8: {
        uint8_t *rval = data;
        *rval = (uint8_t)usys_strtol(val, NULL, 10);
        usys_log_trace("SYSFS:: TYPE_UINT8 Value Read is %d.", *rval);
        break;
    }
    case TYPE_INT8: {
        int8_t *rval = data;
        *rval = (int8_t)usys_strtol(val, NULL, 10);
        usys_log_trace("SYSFS:: TYPE_INT8 Value Read is %d.", *rval);
        break;
    }
    case TYPE_UINT16: {
        uint16_t *rval = data;
        *rval = (uint16_t)usys_strtol(val, NULL, 10);
        usys_log_trace("SYSFS:: TYPE_UINT16 Value Read is %d.", *rval);
        break;
    }
    case TYPE_INT16: {
        int16_t *rval = data;
        *rval = (int16_t)usys_strtol(val, NULL, 10);
        usys_log_trace("SYSFS:: TYPE_INT16 Value Read is %d.", *rval);
        break;
    }
    case TYPE_UINT32: {
        uint32_t *rval = data;
        *rval = (uint32_t)usys_strtol(val, NULL, 10);
        usys_log_trace("SYSFS:: TYPE_UINT32 Value Read is %d.", *rval);
        break;
    }
    case TYPE_INT32: {
        int32_t *rval = data;
        *rval = (int32_t)usys_strtol(val, NULL, 10);
        usys_log_trace("SYSFS:: TYPE_INT32 Value Read is %d.", *rval);
        break;
    }
    case TYPE_INT: {
        int *rval = data;
        *rval = usys_atoi(val);
        usys_log_trace("SYSFS:: TYPE_INT Value Read is %d.", *rval);
        break;
    }
    case TYPE_FLOAT: {
        float *rval = data;
        *rval = (float)usys_strtod(val, NULL);
        usys_log_trace("SYSFS:: TYPE_FLOAT Value Read is %f.", *rval);
        break;
    }
    case TYPE_ENUM: {
        int *rval = data;
        *rval = usys_atoi(val);
        usys_log_trace("SYSFS:: TYPE_ENUM Value Read is %d.", *rval);
        break;
    }
    case TYPE_DOUBLE: {
        double *rval = data;
        *rval = usys_strtod(val, NULL);
        usys_log_trace("SYSFS:: TYPE_DOUBLE Value Read is %lf.", *rval);
        break;
    }
    case TYPE_STRING: {
        /* Should not hit here */
        usys_strcpy(data, val);
        usys_log_trace("SYSFS:: TYPE_STRING Value Read is %s.", data);
        break;
    }
    default: {
    }
    }
}

/* Sysfs always store value as char so type casting it to sensor type to char* */
void sysfs_typetostr(void *data, int type, char *val) {
    switch (type) {
    case TYPE_NULL: {
        break;
    }
    case TYPE_CHAR: {
        usys_sprintf(val, "%c", *(int *)data);
        break;
    }
    case TYPE_BOOL: {
        usys_sprintf(val, "%d", *(uint8_t *)data);
        break;
    }
    case TYPE_UINT8: {
        usys_sprintf(val, "%u", *(uint8_t *)data);
        break;
    }
    case TYPE_INT8: {
        usys_sprintf(val, "%d", *(int8_t *)data);
        break;
    }
    case TYPE_UINT16:
        usys_sprintf(val, "%u", *(uint16_t *)data);
        break;
    case TYPE_INT16: {
        usys_sprintf(val, "%d", *(int16_t *)data);
        break;
    }
    case TYPE_UINT32: {
        usys_sprintf(val, "%d", *(uint32_t *)data);
        break;
    }
    case TYPE_INT32: {
        usys_sprintf(val, "%d", *(int32_t *)data);
        break;
    }
    case TYPE_INT: {
        usys_sprintf(val, "%d", *(int *)data);
        break;
    }
    case TYPE_FLOAT: {
        /* TODO: Sysfs won't understand float and double so rounding */
        int64_t val64 = usys_round(*(float *)data);
        usys_sprintf(val, "%ld", val64);
        break;
    }
    case TYPE_ENUM: {
        usys_sprintf(val, "%d", *(int *)data);
        break;
    }
    case TYPE_DOUBLE: {
        /*TODO: Sysfs won't understand float and double so rounding */
        int64_t val64 = usys_round(*(double *)data);
        usys_sprintf(val, "%ld", val64);
        break;
    }
    case TYPE_STRING: {
        /* Should not hit here */
        char *str = data;
        usys_strcpy(val, str);
        break;
    }
    default: {
    }
    }
}

/* Raw read from sysfs file.*/
int sysfs_read_block(char *name, void *buff, uint16_t size) {
    int read_bytes = 0;
    int fd = sysfs_open(name, O_RDONLY);
    if (fd < 0) {
        read_bytes = -1;
    } else {
        usys_lseek(fd, SYS_DEF_OFFSET, SEEK_SET);
        read_bytes = usys_read(fd, buff, size);
        sysfs_close(fd);
    }
    usys_log_trace("SYSFS:: Read %d bytes from %s file from offset 0x%x.",
                   read_bytes, name, SYS_DEF_OFFSET);
    return read_bytes;
}

/* Raw write to sysfs file.*/
int sysfs_write_block(char *name, void *buff, uint16_t size) {
    int write_bytes = 0;
    int fd = sysfs_open(name, O_WRONLY);
    if (fd < 0) {
        write_bytes = -1;
    } else {
        usys_lseek(fd, SYS_DEF_OFFSET, SEEK_SET);
        write_bytes = usys_write(fd, buff, size);
        sysfs_close(fd);
    }
    usys_log_trace("SYSFS:: Written %d bytes to %s file at offset 0x%x.",
                   write_bytes, name, SYS_DEF_OFFSET);
    return write_bytes;
}

/* Formatted read for numbers to sysfs file.*/
int sysfs_read(char *name, void *data, DataType type) {
    int ret = 0;

    /* Considering max number to be of 32 character long .*/
    int size = SYS_FILE_MAX_LENGTH;
    char val[32] = { '\0' };
    uint16_t idx = 0;
    if (sysfs_read_block(name, val, size) < 0) {
        ret = -1;
    } else {
        sysfs_strtotype(data, type, val);
        usys_log_trace("SYSFS:: Read file %s with String:: %s (Number:: 0x%x).",
                       name, val, *(int64_t *)data);
    }

    return ret;
}

/* Formatted write for numbers to sysfs file.*/
int sysfs_write(char *name, void *data, DataType type) {
    int ret = 0;
    uint16_t idx = 0;
    char val[SYS_FILE_MAX_LENGTH] = { '\0' };
    uint8_t dgt = 0;
    sysfs_typetostr(data, type, val);
    dgt = usys_strlen(val) + 1; //+1 is for '\0'
    if (sysfs_write_block(name, val, dgt) < dgt) {
        ret = -1;
    } else {
        usys_log_trace(
            "SYSFS:: Wrote file %s with String:: %s (Number:: 0x%x).", name,
            val, *(int64_t *)data);
    }
    return ret;
}
