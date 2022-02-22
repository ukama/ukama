/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef SYS_SYNC_H
#define SYS_SYNC_H

#ifdef __cplusplus
extern "C" {
#endif

#include "sys_types.h"
#include "sys_error.h"

/**
 * @brief
 *
 * @param mutex
 * @return USysError
 */
USysError usys_mutex_init(USysMutex* mutex);


/**
 * @brief
 *
 * @param mutex
 * @return USysError
 */
USysError usys_mutex_lock(USysMutex* mutex);


/**
 * @brief
 *
 * @param mutex
 * @return USysError
 */
USysError usys_mutex_trylock(USysMutex* mutex);


/**
 * @brief
 *
 * @param mutex
 * @param wait_time
 * @return USysError
 */
USysError usys_mutex_timedlock_sec(USysMutex* mutex, uint32_t wait_time);


/**
 * @brief
 *
 * @param mutex
 * @param wait_time
 * @return USysError
 */
USysError usys_mutex_timedlock_nsec(USysMutex* mutex, uint32_t wait_time);


/**
 * @brief
 *
 * @param mutex
 * @return USysError
 */
USysError usys_mutex_unlock(USysMutex* mutex);

/**
 * @brief
 *
 * @param mutex
 * @return USysError
 */
USysError usys_mutex_destroy(USysMutex* mutex);

/**
 * @brief
 *
 * @param sem
 * @param init_value
 * @return USysError
 */
USysError usys_sem_init(USysSem* sem, uint32_t init_value);

/**
 * @brief
 *
 * @param sem
 * @return USysError
 */
USysError usys_sem_wait(USysSem* sem);

/**
 * @brief
 *
 * @param sem
 * @return USysError
 */
USysError usys_sem_trywait(USysSem* sem);

/**
 * @brief
 *
 * @param sem
 * @param wait_time
 * @return USysError
 */
USysError usys_sem_timedwait_sec(USysSem* sem, uint32_t wait_time);

/**
 * @brief
 *
 * @param sem
 * @param wait_time
 * @return USysError
 */
USysError usys_sem_timedwait_nsec(USysSem* sem, uint32_t wait_time);

/**
 * @brief
 *
 * @param sem
 * @return USysError
 */
USysError usys_sem_post(USysSem* sem);

/**
 * @brief
 *
 * @param sem
 * @return USysError
 */
USysError usys_sem_destroy(USysSem* sem);

/**
 * @brief
 *
 * @param spinlock
 * @return USysError
 */
USysError usys_spinlock_init(USysSpinlock* spinlock);

/**
 * @brief
 *
 * @param spinlock
 * @return USysError
 */
USysError usys_spinlock_lock(USysSpinlock* spinlock);

/**
 * @brief
 *
 * @param spinlock
 * @return USysError
 */
USysError usys_spinlock_unlock(USysSpinlock* spinlock);

/**
 * @brief
 *
 * @param spinlock
 * @return USysError
 */
USysError usys_spinlock_destroy(USysSpinlock* spinlock);

#ifdef __cplusplus
}
#endif
#endif /*! SYS_SYNC_H */
