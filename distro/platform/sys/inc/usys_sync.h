/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef SYS_SYNC_H
#define SYS_SYNC_H

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"
#include "usys_error.h"

/**
 * @brief
 *
 * @param mutex
 * @return USysError
 */
USysError usys_mutex_init(USysMutex *mutex);

/**
 * @brief
 *
 * @param mutex
 * @return USysError
 */
USysError usys_mutex_lock(USysMutex *mutex);

/**
 * @brief
 *
 * @param mutex
 * @return USysError
 */
USysError usys_mutex_trylock(USysMutex *mutex);

/**
 * @brief
 *
 * @param mutex
 * @param wait_time
 * @return USysError
 */
USysError usys_mutex_timedlock_sec(USysMutex *mutex, uint32_t wait_time);

/**
 * @brief
 *
 * @param mutex
 * @param wait_time
 * @return USysError
 */
USysError usys_mutex_timedlock_nsec(USysMutex *mutex, uint32_t wait_time);

/**
 * @brief
 *
 * @param mutex
 * @return USysError
 */
USysError usys_mutex_unlock(USysMutex *mutex);

/**
 * @brief
 *
 * @param mutex
 * @return USysError
 */
USysError usys_mutex_destroy(USysMutex *mutex);

/**
 * @brief
 *
 * @param sem
 * @param init_value
 * @return USysError
 */
USysError usys_sem_init(USysSem *sem, uint32_t init_value);

/**
 * @fn     USysMutex usys_sem_open*(const char*, int, mode_t, unsigned int)
 * @brief  Creates a new POSIX semaphore or opens an existing semaphore.
 *         The semaphore is identified by name.
 *
 * @param  name
 * @param  oflag
 * @param  mode
 * @param  value
 * @return On success returns the address of the new semaphore
 *         On error return SEM_FAILED
 */
static inline USysSem *usys_sem_open(const char *name, int oflag, mode_t mode,
                                     unsigned int value) {
    return sem_open(name, oflag, mode, value);
}

/**
 * @fn     int usys_sem_close(USysSem*)
 * @brief  closes the named semaphore referred to by sem.
 *
 * @param  sem
 * @return On success 0
 *         On error -1
 */
static inline int usys_sem_close(USysSem *sem) {
    return sem_close(sem);
}
/**
 * @brief
 *
 * @param sem
 * @return USysError
 */
USysError usys_sem_wait(USysSem *sem);

/**
 * @brief
 *
 * @param sem
 * @return USysError
 */
USysError usys_sem_trywait(USysSem *sem);

/**
 * @brief
 *
 * @param sem
 * @param wait_time
 * @return USysError
 */
USysError usys_sem_timedwait_sec(USysSem *sem, uint32_t wait_time);

/**
 * @brief
 *
 * @param sem
 * @param wait_time
 * @return USysError
 */
USysError usys_sem_timedwait_nsec(USysSem *sem, uint32_t wait_time);

/**
 * @brief
 *
 * @param sem
 * @return USysError
 */
USysError usys_sem_post(USysSem *sem);

/**
 * @brief
 *
 * @param sem
 * @return USysError
 */
USysError usys_sem_destroy(USysSem *sem);

/**
 * @brief
 *
 * @param spinlock
 * @return USysError
 */
USysError usys_spinlock_init(USysSpinlock *spinlock);

/**
 * @brief
 *
 * @param spinlock
 * @return USysError
 */
USysError usys_spinlock_lock(USysSpinlock *spinlock);

/**
 * @brief
 *
 * @param spinlock
 * @return USysError
 */
USysError usys_spinlock_unlock(USysSpinlock *spinlock);

/**
 * @brief
 *
 * @param spinlock
 * @return USysError
 */
USysError usys_spinlock_destroy(USysSpinlock *spinlock);

#ifdef __cplusplus
}
#endif
#endif /*! SYS_SYNC_H */
