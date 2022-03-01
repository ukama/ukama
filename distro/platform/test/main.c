/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_api.h"
#include "usys_error.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_sync.h"
#include "usys_thread.h"
#include "usys_types.h"
#include "usys_mem.h"
#include "test.h"

#include "unity.h"


int main() {
    usys_log_set_level(LOG_TRACE);

    usys_log_info("Starting test app for platform.");

    UNITY_BEGIN();
    RUN_TEST(test_usys_errors);
    RUN_TEST(test_usys_fopen_should_return_file_pointer);
    RUN_TEST(test_usys_fopen_should_return_null);
    RUN_TEST(test_usys_fopen_create_new_file);
    RUN_TEST(test_usys_fread_fwrite_fseek);
    RUN_TEST(test_usys_fgets);
    RUN_TEST(test_usys_threads_with_mutex);
    RUN_TEST(test_usys_threads_with_semaphore);
    RUN_TEST(test_usys_memcmp_memset_memcpy);
    RUN_TEST(test_usys_strcat);
    RUN_TEST(test_usys_strcpy);
    RUN_TEST(test_usys_strncmp);
    RUN_TEST(test_usys_strlen);
    RUN_TEST(test_usys_strstr);
    RUN_TEST(test_usys_strtok);
    RUN_TEST(test_shm);
    return UNITY_END();
}





