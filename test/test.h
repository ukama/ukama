/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
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

void test_shm(void);

#endif /* USYS_UNIT_TEST_CASE_H_ */
