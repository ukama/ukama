/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_sync.h"
#include "usys_log.h"

USysError usys_mutex_init(USysMutex* mutex) {
    if (mutex == NULL) {
        usys_log_warn("Mutex Object is NULL");

        return ERR_MUTEX_OBJ_NULL;
    }

    pthread_mutexattr_t mutex_attr;

    if (pthread_mutexattr_init(&mutex_attr) != 0) {
        usys_log_warn("Mutex attribute init failed");

        return ERR_MUTEX_ATTR_INIT_FAIL;
    }

    if (pthread_mutexattr_setprotocol(&mutex_attr, PTHREAD_PRIO_INHERIT) != 0) {
        usys_log_warn("Mutex attribute set protocol failed");

        return ERR_MUTEX_ATTR_SET_PROTO_FAIL;
    }

    if (pthread_mutexattr_settype(&mutex_attr, PTHREAD_MUTEX_RECURSIVE) != 0) {
        usys_log_warn("Mutex attribute set type RECURSIVE failed");

        return ERR_MUTEX_ATTR_SET_TYPE_FAIL;
    }

    if (pthread_mutex_init(mutex, &mutex_attr) != 0) {
        usys_log_warn("Mutex init with attribute failed");

        return ERR_MUTEX_INIT_FAILED;
    }

    if (pthread_mutexattr_destroy(&mutex_attr) != 0) {
        usys_log_warn("Mutex attr destroy failed");
    }

    return ERR_NONE;
}

USysError usys_mutex_lock(USysMutex* mutex) {
    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_MUTEX_OBJ_NULL;
    }

    if (pthread_mutex_lock(mutex) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Mutex lock failed");
        return ERR_MUTEX_LOCK_FAILED;
    }
}

USysError usys_mutex_trylock(USysMutex* mutex) {
    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_MUTEX_OBJ_NULL;
    }

    if (pthread_mutex_trylock(mutex) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Mutex trylock failed");
        return ERR_MUTEX_TRYLOCK_FAILED;
    }

}

USysError usys_mutex_timedlock_sec(USysMutex* mutex, uint32_t wait_time) {
    struct timespec abs_time;

    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_MUTEX_OBJ_NULL;
    }

    clock_gettime(CLOCK_REALTIME, &abs_time);
    abs_time.tv_sec += wait_time;

    if (pthread_mutex_timedlock(mutex, &abs_time) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Mutex timedlock failed");
        return ERR_MUTEX_TIMEDLOCK_FAILED;
    }
}

USysError usys_mutex_timedlock_nsec(USysMutex* mutex, uint32_t wait_time) {
    struct timespec abs_time;

    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_MUTEX_OBJ_NULL;
    }

    clock_gettime(CLOCK_REALTIME, &abs_time);
    abs_time.tv_nsec += wait_time;

    if (pthread_mutex_timedlock(mutex, &abs_time) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Mutex timedlock failed");
        return ERR_MUTEX_TIMEDLOCK_FAILED;
    }
}

USysError usys_mutex_unlock(USysMutex* mutex) {
    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_MUTEX_OBJ_NULL;
    }

    if (pthread_mutex_unlock(mutex) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Mutex unlock failed");
        return ERR_MUTEX_UNLOCK_FAILED;
    }
}

USysError usys_mutex_destroy(USysMutex* mutex) {
    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_MUTEX_OBJ_NULL;
    }

    if (pthread_mutex_destroy(mutex) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Mutex destroy failed");
        return ERR_MUTEX_DESTROY_FAILED;
    }
}

USysError usys_sem_init(USysSem* sem, uint32_t init_value) {
    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_SEM_OBJ_NULL;
    }

    if (sem_init(sem, 0, init_value) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Semaphore sem_init failed");
        return ERR_SEM_INIT_FAILURE;
    }
}

USysError usys_sem_wait(USysSem* sem) {
    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_SEM_OBJ_NULL;
    }

    if (sem_wait(sem) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Semaphore failed sem_wait failed");
        return ERR_SEM_WAIT_FAIL;
    }
}

USysError usys_sem_trywait(USysSem* sem) {
    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_SEM_OBJ_NULL;
    }

    if (sem_trywait(sem) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Semaphore sem_trywait failed");
        return ERR_SEM_TRYWAIT_FAIL;
    }
}

USysError usys_sem_timedwait_sec(USysSem* sem, uint32_t wait_time) {
    struct timespec abs_time;

    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_SEM_OBJ_NULL;
    }

    clock_gettime(CLOCK_REALTIME, &abs_time);
    abs_time.tv_sec += wait_time;

    if (sem_timedwait(sem,  &abs_time) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Semaphore sem_timedwait failed");
        return ERR_SEM_TIMEDWAIT_FAIL;
    }
}

USysError usys_sem_timedwait_nsec(USysSem* sem, uint32_t wait_time) {
    struct timespec abs_time;

    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_SEM_OBJ_NULL;
    }

    clock_gettime(CLOCK_REALTIME, &abs_time);
    abs_time.tv_nsec += wait_time;

    if (sem_timedwait(sem,  &abs_time) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Semaphore sem_timedwait failed");
        return ERR_SEM_TIMEDWAIT_FAIL;
    }
}

USysError usys_sem_post(USysSem* sem) {
    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_SEM_OBJ_NULL;
    }

    if (sem_post(sem) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Semaphore sem_post failed");
        return ERR_SEM_POST_FAIL;
    }
}

USysError usys_sem_destroy(USysSem* sem) {
    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_SEM_OBJ_NULL;
    }

    if (sem_destroy(sem) == 0) {
        return ERR_NONE;
    } else {
        usys_log_warn("Semaphore sem_destroy failed");
        return ERR_SEM_DESTROY_FAIL;
    }
}

USysError usys_spinlock_init(USysSpinlock* spinlock) {
    int ret = pthread_spin_init(spinlock, PTHREAD_PROCESS_PRIVATE);

    if (ret == 0) {
        return ERR_NONE;
    } else {
        return ERR_SPIN_LOCK_INIT_FAILED;
    }
}

USysError usys_spinlock_lock(USysSpinlock* spinlock) {
    int ret = pthread_spin_lock(spinlock);
    if (ret == 0) {
        return ERR_NONE;
    } else {
        return ERR_SPIN_LOCK_LOCK_FAILED;
    }
}

USysError usys_spinlock_unlock(USysSpinlock* spinlock) {
    int ret = pthread_spin_unlock(spinlock);
    if (ret == 0) {
        return ERR_NONE;
    } else {
        return ERR_SPIN_LOCK_UNLOCK_FAILED;
    }
}

USysError usys_spinlock_destroy(USysSpinlock* spinlock) {
    int ret = pthread_spin_destroy(spinlock);
    if (ret == 0) {
        return ERR_NONE;
    } else {
        return ERR_SPIN_LOCK_DESTROY_FAILED;
    }
}

