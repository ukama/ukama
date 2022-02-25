/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "test.h"

#include "unity.h"
#include "usys_api.h"
#include "usys_error.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_sync.h"
#include "usys_thread.h"

void setUp(void) {
    // set stuff up here
}

void tearDown(void) {
    // clean stuff up here
}

/* File Operations test */
void test_usys_fopen_should_return_file_pointer(void) {
    TEST_ASSERT_NOT_EQUAL(NULL, usys_fopen("./test/data/sample.txt", "r"));
}

void test_usys_fopen_should_return_null(void) {
    TEST_ASSERT_EQUAL(NULL, usys_fopen("./test/data/nosuchfile.txt", "r"));
}

void test_usys_fopen_create_new_file(void) {
    TEST_ASSERT_NOT_EQUAL(NULL, usys_fopen("./test/data/newfile.txt", "w+"));
}

void test_usys_fread_fwrite_fseek(void) {
    char data[] = "abcdef";
    char rdata[10];

    FILE* fp = usys_fopen("./test/data/newfile.txt", "w+");
    TEST_ASSERT_MESSAGE(fp!= NULL, "failed to open file.");
    if (fp != NULL) {
        int wbytes = usys_fwrite(data, sizeof(char), sizeof(data), fp);
        TEST_ASSERT_MESSAGE(wbytes == sizeof(data), "Bytes Written not matching");

        int seek = usys_fseek(fp, 0, SEEK_SET);
        TEST_ASSERT_MESSAGE(seek == 0, "Setting file pointer to beginning of file failed.");

        int rbytes = usys_fread(&rdata, sizeof(char), sizeof(data), fp);
        TEST_ASSERT_MESSAGE(rbytes == sizeof(data), "Bytes read not matching");

        TEST_ASSERT_MESSAGE(usys_memcmp(data, rdata, wbytes) == 0, "Read data not matching to written data");
    }

    TEST_ASSERT_MESSAGE(0 == usys_fclose(fp), "Failed to close file.");
}

void test_usys_fgets(void) {
    char data[] = "abcdef";
    char rdata[10];

    FILE* fp = usys_fopen("./test/data/newfile.txt", "r");
    TEST_ASSERT_MESSAGE(fp!= NULL, "failed to open file.");
    if (fp != NULL) {

        usys_fgets(rdata, sizeof(data), fp);

        TEST_ASSERT_MESSAGE(usys_memcmp(data, rdata, sizeof(data)) == 0, "Read data not matching to written data");

    }

    TEST_ASSERT_MESSAGE(0 == usys_fclose(fp), "Failed to close file.");
}

/* Test cases related to thread */

static int counter = 0;
USysThreadId thread[2];

/* Thread sync with mutex */
USysMutex mutex;
void *thread_function_with_mutex(void* arg) {
    USysError err;
    int *tstat = (int*)arg;

    err = usys_mutex_lock(&mutex);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

    counter += 1;
    usys_log_trace("Job %d has started", counter);

    TEST_ASSERT_MESSAGE(0 == usys_sleep(2), "Sleep failed.");


    usys_log_trace("Job %d has finished", counter);

    err = usys_mutex_unlock(&mutex);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

    *tstat = counter;
    pthread_exit(tstat);
}

void test_usys_threads_with_mutex() {
    int idx = 0;
    counter = 0;
    USysError err;
    int ret = 0;
    int arg[2] = {0};
    int *status[2] = {0};
    err = usys_mutex_init(&mutex);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

    while (idx < 2) {
        err = usys_thread_create( &thread[idx],
                NULL,
                &thread_function_with_mutex,
                (void*)&arg[idx]);
        TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

        idx++;
    }

    ret = usys_thread_join(thread[0], (void **)(&status[0]));
    TEST_ASSERT_MESSAGE(ret == 0, "Failed to join thread 1.");
    TEST_ASSERT_EQUAL_INT(*status[0], 1 );

    ret = usys_thread_join(thread[1], (void **)(&status[1]));
    TEST_ASSERT_MESSAGE(ret == 0, "Failed to join thread 2.");
    TEST_ASSERT_EQUAL_INT(*status[1], 2 );

    /* Mutex destroy */
    err = usys_mutex_destroy(&mutex);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

}

/* Thread sync with semaphore */

USysSem sem;
void *thread_function_with_semaphore(void* arg) {
    USysError err;
    int *tstat = (int*)arg;

    err = usys_sem_wait(&sem);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

    counter += 1;
    usys_log_trace("Job %d has started", counter);

    TEST_ASSERT_MESSAGE(0 == usys_sleep(2), "Sleep failed.");


    usys_log_trace("Job %d has finished", counter);

    err = usys_sem_post(&sem);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

    *tstat = counter;
    pthread_exit(tstat);
}

void test_usys_threads_with_semaphore() {
    int idx = 0;
    counter = 0;
    USysError err;
    int ret = 0;
    int arg[2] = {0};
    int *status[2] = {0};
    err = usys_sem_init(&sem, 1);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

    while (idx < 2) {
        err = usys_thread_create( &thread[idx],
                NULL,
                &thread_function_with_semaphore,
                (void*)&arg[idx]);
        TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

        idx++;
    }

    ret = usys_thread_join(thread[0], (void **)(&status[0]));
    TEST_ASSERT_MESSAGE(ret == 0, "Failed to join thread 1.");
    TEST_ASSERT_EQUAL_INT(*status[0], 1 );

    ret = usys_thread_join(thread[1], (void **)(&status[1]));
    TEST_ASSERT_MESSAGE(ret == 0, "Failed to join thread 2.");
    TEST_ASSERT_EQUAL_INT(*status[1], 2 );

    /* Mutex destroy */
    err = usys_sem_destroy(&sem);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

}
