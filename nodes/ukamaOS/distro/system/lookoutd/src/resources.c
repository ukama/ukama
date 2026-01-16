/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>

#include "lookout.h"

int get_memory_usage(int pid) {

    FILE *file = NULL;

    char filename[MAX_BUFFER] = {0};
    char line[MAX_BUFFER]     = {0};
    long int vmRSS            = 0;

    sprintf(filename, "/proc/%d/status", pid);
    file = fopen(filename, "r");
    if (!file) {
        usys_log_error("Failed to open status file: %s", filename);
        return 0;
    }

    while (fgets(line, sizeof(line), file)) {
        if (strncmp(line, "VmRSS:", 6) == 0) {
            sscanf(line, "VmRSS: %ld kB", &vmRSS);
            break;
        }
    }

    fclose(file);
    return vmRSS;
}

int get_disk_usage(int pid) {

    FILE *file = NULL;
    
    char filename[MAX_BUFFER] = {0};
    char line[MAX_BUFFER]     = {0};
    long int bytes            = 0;
    
    sprintf(filename, "/proc/%d/io", pid);

    file = fopen(filename, "r");
    if (!file) {
        usys_log_error("Failed to open status file: %s", filename);
        return 0;
    }

    while (fgets(line, sizeof(line), file)) {
        if (strncmp(line, "write_bytes:", 12) == 0) {
            sscanf(line, "write_bytes: %ld", &bytes);
        }
    }

    fclose(file);
    return bytes;
}

double get_cpu_usage(int pid) {

    FILE *file = NULL;

    char filename[MAX_BUFFER] = {0};
    char comm[256];
    char state;
    int  ppid;
    int  pgrp;
    int  session;
    int  tty_nr;
    int  tpgid;

    unsigned int flags;
    unsigned long minflt;
    unsigned long cminflt;
    unsigned long majflt;
    unsigned long cmajflt;
    unsigned long utime;
    unsigned long stime;

    long cutime;
    long cstime;
    long priority;
    long nice;
    long num_threads;
    long itrealvalue;
    unsigned long long starttime=0;
    
    sprintf(filename, "/proc/%d/stat", pid);

    file = fopen(filename, "r");
    if (!file) {
        usys_log_error("Failed to open status file: %s", filename);
        return 0;
    }

fscanf(file, "%d %255s %c %d %d %d %d %d %u %lu %lu %lu %lu %lu %lu %ld %ld %ld %ld %ld %ld %llu",
       &pid, comm, &state, &ppid, &pgrp, &session, &tty_nr, &tpgid, &flags,
       &minflt, &cminflt, &majflt, &cmajflt, &utime, &stime, &cutime, &cstime,
       &priority, &nice, &num_threads, &itrealvalue, &starttime);

    fclose(file);

    return (double)100.0 * ( utime + stime + cutime + cstime) / sysconf(_SC_CLK_TCK);
}

char *get_radio_status(void) {

    return "on";
}
