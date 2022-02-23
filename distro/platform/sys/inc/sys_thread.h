/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#ifndef USYS_SYS_THREAD_H
#define USYS_SYS_THREAD_H

#ifdef __cplusplus
extern "C" {
#endif

#include "sys_types.h"

/**
 * @fn     int usys_thread_attr_init()
 * @brief  initializes the thread
 *         attributes object pointed to by attr with default attribute
 *         values.
 *
 * @param  attr
 * @return On success, these functions return 0; on error, they return a
 *         nonzero error number.
 */
static inline int usys_thread_attr_init(USysThreadAttr *attr) {
    return pthread_attr_init(attr);
}

/**
 * @fn     int usys_thread_attr_destroy()
 * @brief  thread attributes object is destroyed using this function
 *
 * @return On success, these functions return 0; on error, they return a
 *         nonzero error number.
 */
static inline int usys_thread_attr_destroy(USysThreadAttr *attr) {
    return pthread_attr_destroy(attr);
}

/**
 * @fn    int usys_thread_exit(void*)
 * @brief Terminates the calling thread
 *
 * @param status
 */
static inline void usys_thread_exit(void* status) {
   return pthread_exit(status);
}

/**
 * @fn     int usys_thread_cancel()
 * @brief  send a cancellation request to a thread
 *
 * @return returns 0; on error, it returns a
           nonzero error number.
 */
static inline int usys_thread_cancel(USysThreadId thread) {
    return pthread_cancel(thread);
}

/**
 * @fn     int usys_thread_detach(USysThreadId)
 * @brief  function marks the thread identified by thread as detached
 *
 * @param  thread
 * @return On success, pthread_detach() returns 0; on error, it returns an
 *         error number.
 */
static inline int usys_thread_detach(USysThreadId thread){
    return pthread_detach(thread);
}

/**
 * @fn     int usys_thread_join(USysThreadId, void**)
 * @brief  join with a terminated thread
 *
 * @param  thread
 * @param  status
 * @return On success, pthread_join() returns 0; on error, it returns an
 *         error number.
 */
int usys_thread_join(USysThreadId thread, void **status){
    return pthread_join(thread, status);
}

/**
 * @fn     USysThreadId usys_thread_id(void)
 * @brief  returns the ID of the calling thread.
 *
 * @return Thread ID
 */
static inline USysThreadId usys_thread_id(void){
    return pthread_self();
}

/**
 * @fn     int usys_thread_create(USysThreadId*, const USysThreadAttr*, void*(*)(void*), void*)
 * @brief  Creates a new thread. This function starts a new thread in the calling
 *         process.  The new thread starts execution by invoking
 *         start_routine(); arg is passed as the sole argument of
 *         start_routine()
 *
 * @param  thread
 * @param  attr
 * @param  start_routine
 * @param  arg
 * @return On success, pthread_create() returns 0; on error, it returns an
 *         error number, and the contents of *thread are undefined
 */
static inline int usys_thread_create(USysThreadId *thread, const USysThreadAttr *attr,
            void *(*start_routine)(void *), void *arg){
    return pthread_create(thread, attr, start_routine, arg);
}

/**
 * @fn     int usys_thread_setschedparam(pthread_t, int, const struct sched_param*)
 * @brief  sets the scheduling policy and parameters of the thread thread
 *
 * @param  thread
 * @param  policy
 * @param  param
 * @return On success, these functions return 0; on error, they return a
 *         nonzero error number.  If pthread_setschedparam() fails, the
 *         scheduling policy and parameters of thread are not changed.
 */
static inline int usys_thread_setschedparam(USysThreadId thread, int policy,
               const struct sched_param *param) {
    return  pthread_setschedparam(thread, policy,
                   param);
}

/**
 * @fn     int usys_thread_attr_setstacksize(USysThreadId*, size_t)
 * @brief  sets the stack size
 *         attribute of the thread attributes object referred to by attr to
 *         the value specified in stacksize.
 *
 * @param  attr
 * @param  stacksize
 * @return On success, these functions return 0; on error, they return a
 *         nonzero error number.
 */
static inline int usys_thread_attr_setstacksize(USysThreadAttr *attr, size_t stacksize) {
    return pthread_attr_setstacksize(attr, stacksize);
}

#endif /* USYS_SYS_THREAD_H */
