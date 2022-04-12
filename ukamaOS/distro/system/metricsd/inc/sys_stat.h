/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>
#include <dirent.h>
#include <ctype.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/statvfs.h>
#include <unistd.h>

#ifndef INC_SYS_STAT_H_
#define INC_SYS_STAT_H_

/* Maximum length of block device name */
#define MAX_DEV_LEN          128
/* Maximum length of network interface name */
#define MAX_IFACE_LEN        16
#define MAX_PF_NAME          1024

/* FILES */
#define PROC_CPU_INFO        "/proc/cpuinfo"
#define PROC_NET_DEV         "/proc/net/dev"
#define PROC_CPU_STAT        "/proc/stat"
#define PROC_MEM_STAT        "/proc/meminfo"
#define SYS_DEV_CPU          "/sys/devices/system/cpu"
#define PROC_UPTIME          "/proc/uptime"

/* CPU */
typedef struct {
    unsigned long long cpuUser;
    unsigned long long cpuNice;
    unsigned long long cpuSys;
    unsigned long long cpuIdle;
    unsigned long long cpuIowait;
    unsigned long long cpuSteal;
    unsigned long long cpuHardirq;
    unsigned long long cpuSoftirq;
    unsigned long long cpuGuest;
    unsigned long long cpuGuestNice;
    double freq;
} SysCPUMetrics;

/* Network */
typedef struct {
    unsigned long long collisions;
    unsigned long long rxErrors;
    unsigned long long txErrors;
    unsigned long long rxDropped;
    unsigned long long txDropped;
    unsigned long long rxFifoErrors;
    unsigned long long txFifoErrors;
    unsigned long long rxOverruns;
    unsigned long long txCarrierErrors;
    unsigned long long rxPackets;
    unsigned long long txPackets;
    unsigned long long rxBytes;
    unsigned long long txBytes;
    unsigned long long rxCompressed;
    unsigned long long txCompressed;
    unsigned long long multicast;
    unsigned int linkspeed;
    unsigned int linkstatus;
    unsigned long long int latency;
    char interface[MAX_IFACE_LEN];
    char duplex;
} SysNetDevMetrics;

/* CPU Frequency */
typedef struct {
    unsigned long cpufreq __attribute__ ((aligned (8)));
} SysCPUFreq;

/* DDR stats */
typedef struct {
    unsigned long long memTotal;
    unsigned long long memUsed;
    unsigned long long memFree;
    unsigned long long memAvail;
    unsigned long long memBuffer;
    unsigned long long memCached;
} SysMemDDRMetrics;

/* Swap memory stats */
typedef struct {
    unsigned long long total;
    unsigned long long used;
    unsigned long long free;
} SysMemSwapMetrics;

/* Memory stats */
typedef struct {
    SysMemDDRMetrics ddr;
    SysMemSwapMetrics swap;
} SysMemMetrics;

/* Generic stats for system */
typedef struct {
    double uptime;
} SysGenMetrics;

/* Storage stats for EMMC */
typedef struct {
    unsigned long long blksize;
    unsigned long long total;
    unsigned long long used;
    unsigned long long free; /* Free for available for unprivileged user */
    unsigned long long pfree; /* Free blocks */
} SysStorageMetrics;

#endif /* INC_SYS_STAT_H_ */
