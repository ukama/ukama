/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_SYS_TYPES_H
#define USYS_SYS_TYPES_H

#ifdef __cplusplus
extern "C" {
#endif

#include <assert.h>
#include <ctype.h>
#include <endian.h>
#include <errno.h>
#include <fcntl.h>
#include <math.h>
#include <netdb.h>
#include <pthread.h>
#include <sched.h>
#include <semaphore.h>
#include <signal.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>
#include <asm/types.h>
#include <arpa/inet.h>
#include <linux/if_ether.h>
#include <linux/if_packet.h>
#include <linux/if_tun.h>
#include <linux/sctp.h>
#include <linux/types.h>
#include <net/if.h>
#include <netinet/ip.h>
#include <netinet/udp.h>
#include <sys/ioctl.h>
#include <sys/mman.h>
#include <sys/poll.h>
#include <sys/queue.h>
#include <sys/socket.h>
#include <sys/stat.h>
#include <sys/time.h>
#include <sys/types.h>


/**
 * @typedef USysMachineWord
 *
 * @brief  A machine word is the amount of data that a machine can process at one time
 *
 */
#if (ARCH == 64)
typedef uint64_t USysMachineWord;
#else
typedef uint32_t USysMachineWord;
#endif

/**
 * @typedef USysMutex
 *
 * @brief mutex object
 *
 */
typedef pthread_mutex_t USysMutex;

/**
 * @typedef USysSem
 *
 * @brief semaphore object
 *
 */
typedef sem_t USysSem;


/**
 * @typedef USysSpinLock
 *
 * @brief spinlock object
 *
 */
typedef pthread_spinlock_t USysSpinlock;

/**
 * @typedef USysThreadId
 *
 * @brief Thread identifier
 *
 */
typedef pthread_t USysThreadId;

/**
 * @typedef  USysThreadKey
 *
 * @brief  thread keys
 *
 */
typedef pthread_key_t USysThreadKey;

/**
 * @typedef USysSharedMemMgrHandle
 *
 * @brief  System Shared Memory Management Handle Type Definition
 *
 */
typedef void* USysSharedMemMgrHandle;


/**
 * @typedef USysPhysAddr
 *
 * @brief  System Shared Memory Management Handle Type Definition
 *
 */

#if (ARCH == 64)
typedef uint64_t USysPhysAddr;
#else
/* ARCH == 32 */
typedef uint32_t USysPhysAddr;
#endif


#ifdef __cplusplus
}
#endif
#endif
