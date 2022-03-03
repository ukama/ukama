/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "test.h"
#include "usys_log.h"

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
    RUN_TEST(test_usys_shm);
    RUN_TEST(test_usys_timer);
    RUN_TEST(test_usys_read_write_arrays_to_file);
    RUN_TEST(test_usys_read_write_numbers_to_file);
    RUN_TEST(test_usys_read_write_strings_to_file);
    RUN_TEST(test_usys_open_read_close_dir);
    RUN_TEST(test_usys_opendir_fail);
    RUN_TEST(test_usys_getcwd_ch_dir);
    RUN_TEST(test_usys_seek_tell_dir);
    RUN_TEST(test_usys_fork_wait_pid_ppid_prgp);
    RUN_TEST(test_usys_get_set_rlimit);
    return UNITY_END();
}





