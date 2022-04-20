/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_sync.h"
#include "usys_log.h"

USysError usys_mutex_init(USysMutex *mutex) {
  pthread_mutexattr_t mutexAttr;

    if (mutex == NULL) {
        usys_log_warn("Mutex Object is NULL");
        return ERR_PLTF_MUTEX_OBJ_NULL;
    }

    if (pthread_mutexattr_init(&mutexAttr) != 0) {
        usys_log_warn("Mutex attribute init failed");
        return ERR_PLTF_MUTEX_ATTR_INIT_FAIL;
    }

    if (pthread_mutexattr_setprotocol(&mutexAttr, PTHREAD_PRIO_INHERIT) != 0) {
        usys_log_warn("Mutex attribute set protocol failed");
        return ERR_PLTF_MUTEX_ATTR_SET_PROTO_FAIL;
    }

    if (pthread_mutexattr_settype(&mutexAttr, PTHREAD_MUTEX_RECURSIVE) != 0) {
        usys_log_warn("Mutex attribute set type RECURSIVE failed");
        return ERR_PLTF_MUTEX_ATTR_SET_TYPE_FAIL;
    }

    if (pthread_mutex_init(mutex, &mutexAttr) != 0) {
        usys_log_warn("Mutex init with attribute failed");
        return ERR_PLTF_MUTEX_INIT_FAILED;
    }

    if (pthread_mutexattr_destroy(&mutexAttr) != 0) {
        usys_log_warn("Mutex attr destroy failed");
    }

    return ERR_PLTF_NONE;
}

USysError usys_mutex_lock(USysMutex *mutex) {
  if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_PLTF_MUTEX_OBJ_NULL;
    }

  if (pthread_mutex_lock(mutex) != 0) {
        usys_log_warn("Mutex lock failed");
        return ERR_PLTF_MUTEX_LOCK_FAILED;
  }

    return ERR_PLTF_NONE;
}

USysError usys_mutex_trylock(USysMutex *mutex) {
    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_PLTF_MUTEX_OBJ_NULL;
    }

    if (pthread_mutex_trylock(mutex) != 0) {
        usys_log_warn("Mutex trylock failed");
        return ERR_PLTF_MUTEX_TRYLOCK_FAILED;
    }

    return ERR_PLTF_NONE;
}

USysError usys_mutex_timedlock_sec(USysMutex *mutex, uint32_t wait_time) {
    struct timespec absTime;

    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_PLTF_MUTEX_OBJ_NULL;
    }

    clock_gettime(CLOCK_REALTIME, &absTime);
    absTime.tv_sec += wait_time;

    if (pthread_mutex_timedlock(mutex, &absTime) != 0) {
        usys_log_warn("Mutex timedlock failed");
        return ERR_PLTF_MUTEX_TIMEDLOCK_FAILED;
    }

    return ERR_PLTF_NONE;
}

USysError usys_mutex_timedlock_nsec(USysMutex *mutex, uint32_t wait_time) {
    struct timespec absTime;

    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_PLTF_MUTEX_OBJ_NULL;
    }

    clock_gettime(CLOCK_REALTIME, &absTime);
    absTime.tv_nsec += wait_time;

    if (pthread_mutex_timedlock(mutex, &absTime) != 0) {
        usys_log_warn("Mutex timedlock failed");
        return ERR_PLTF_MUTEX_TIMEDLOCK_FAILED;
    }

    return ERR_PLTF_NONE;
}

USysError usys_mutex_unlock(USysMutex *mutex) {
    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_PLTF_MUTEX_OBJ_NULL;
    }

    if (pthread_mutex_unlock(mutex) != 0) {
        usys_log_warn("Mutex unlock failed");
        return ERR_PLTF_MUTEX_UNLOCK_FAILED;
    }

    return ERR_PLTF_NONE;
}

USysError usys_mutex_destroy(USysMutex *mutex) {
    if (mutex == NULL) {
        usys_log_warn("Mutex Object NULL");
        return ERR_PLTF_MUTEX_OBJ_NULL;
    }

    if (pthread_mutex_destroy(mutex) != 0) {
        usys_log_warn("Mutex destroy failed");
        return ERR_PLTF_MUTEX_DESTROY_FAILED;
    }

    return ERR_PLTF_NONE;
}

USysError usys_sem_init(USysSem *sem, uint32_t init_value) {
    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_PLTF_SEM_OBJ_NULL;
    }

    if (sem_init(sem, 0, init_value) != 0) {
        usys_log_warn("Semaphore sem_init failed");
        return ERR_PLTF_SEM_INIT_FAILURE;
    }

    return ERR_PLTF_NONE;
}

USysError usys_sem_wait(USysSem *sem) {
    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_PLTF_SEM_OBJ_NULL;
    }

    if (sem_wait(sem) != 0) {
        usys_log_warn("Semaphore failed sem_wait failed");
        return ERR_PLTF_SEM_WAIT_FAIL;
    }

    return ERR_PLTF_NONE;
}

USysError usys_sem_trywait(USysSem *sem) {
    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_PLTF_SEM_OBJ_NULL;
    }

    if (sem_trywait(sem) != 0) {
        usys_log_warn("Semaphore sem_trywait failed");
        return ERR_PLTF_SEM_TRYWAIT_FAIL;
    }

    return ERR_PLTF_NONE;
}

USysError usys_sem_timedwait_sec(USysSem *sem, uint32_t wait_time) {
    struct timespec absTime;

    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_PLTF_SEM_OBJ_NULL;
    }

    clock_gettime(CLOCK_REALTIME, &absTime);
    absTime.tv_sec += wait_time;

    if (sem_timedwait(sem, &absTime) != 0) {
        usys_log_warn("Semaphore sem_timedwait failed");
        return ERR_PLTF_SEM_TIMEDWAIT_FAIL;
    }

    return ERR_PLTF_NONE;
}

USysError usys_sem_timedwait_nsec(USysSem *sem, uint32_t wait_time) {
    struct timespec absTime;

    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_PLTF_SEM_OBJ_NULL;
    }

    clock_gettime(CLOCK_REALTIME, &absTime);
    absTime.tv_nsec += wait_time;

    if (sem_timedwait(sem, &absTime) != 0) {
        usys_log_warn("Semaphore sem_timedwait failed");
        return ERR_PLTF_SEM_TIMEDWAIT_FAIL;
    }

    return ERR_PLTF_NONE;
}

USysError usys_sem_post(USysSem *sem) {
    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_PLTF_SEM_OBJ_NULL;
    }

    if (sem_post(sem) != 0) {
        usys_log_warn("Semaphore sem_post failed");
        return ERR_PLTF_SEM_POST_FAIL;
    }

    return ERR_PLTF_NONE;
}

USysError usys_sem_destroy(USysSem *sem) {
    if (sem == NULL) {
        usys_log_warn("Semaphore Object NULL");
        return ERR_PLTF_SEM_OBJ_NULL;
    }

    if (sem_destroy(sem) != 0) {
        usys_log_warn("Semaphore sem_destroy failed");
        return ERR_PLTF_SEM_DESTROY_FAIL;
    }

    return ERR_PLTF_NONE;
}

USysError usys_spinlock_init(USysSpinlock *spinlock) {
  if(pthread_spin_init(spinlock, PTHREAD_PROCESS_PRIVATE) != 0) {
    usys_log_warn("spinlock  initialization failed");
    return ERR_PLTF_SPIN_LOCK_INIT_FAILED;
  }

  return ERR_PLTF_NONE;
}

USysError usys_spinlock_lock(USysSpinlock *spinlock) {
  if(pthread_spin_lock(spinlock) != 0) {
    usys_log_warn("spinlock lock failed");
    return ERR_PLTF_SPIN_LOCK_LOCK_FAILED;
  }

  return ERR_PLTF_NONE;
}

USysError usys_spinlock_unlock(USysSpinlock *spinlock) {
  if(pthread_spin_unlock(spinlock) != 0) {
    usys_log_warn("spinlock unlock failed");
    return ERR_PLTF_SPIN_LOCK_UNLOCK_FAILED;
  }

  return ERR_PLTF_NONE;
}

USysError usys_spinlock_destroy(USysSpinlock *spinlock) {
  if(pthread_spin_destroy(spinlock) != 0) {
    usys_log_warn("spinlock destroy failed");
    return ERR_PLTF_SPIN_LOCK_DESTROY_FAILED;
  }

  return ERR_PLTF_NONE;
}
