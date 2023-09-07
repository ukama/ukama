/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_UNIT_TEST_CASE_H_
#define USYS_UNIT_TEST_CASE_H_

#ifdef __cplusplus
extern "C" {
#endif

void setUp(void);
void tearDown(void);
void test_usys_errors(void);
void test_usys_fopen_should_return_file_pointer(void);
void test_usys_fopen_should_return_null(void);
void test_usys_fopen_create_new_file(void);
void test_usys_fread_fwrite_fseek(void);
void test_usys_fgets(void);
void test_usys_threads_with_mutex(void);
void test_usys_threads_with_semaphore(void);

void test_usys_memcmp_memset_memcpy(void);
void test_usys_strcat(void);
void test_usys_strcpy(void);
void test_usys_strncmp(void);
void test_usys_strlen(void);
void test_usys_strstr(void);
void test_usys_strtok(void);
void test_usys_shm(void);
void test_usys_timer(void);
void test_usys_read_write_arrays_to_file(void);
void test_usys_read_write_numbers_to_file(void);
void test_usys_read_write_strings_to_file(void);
void test_usys_write_failure_no_file_exist(void);
void test_usys_open_read_close_dir(void);
void test_usys_opendir_fail(void);
void test_usys_getcwd_ch_dir(void);
void test_usys_seek_tell_dir(void);
void test_usys_fork_wait_pid_ppid_prgp(void);
void test_usys_get_set_rlimit(void);

#endif /* USYS_UNIT_TEST_CASE_H_ */
