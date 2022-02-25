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

void setUp(void) {
    // set stuff up here
}

void tearDown(void) {
    // clean stuff up here
}

void test_usys_fopen_should_return_file_pointer(void) {
    TEST_ASSERT_NOT_EQUAL(NULL, usys_fopen("./test/data/sample.txt", "r"));
}

void test_usys_fopen_should_return_null(void) {
    TEST_ASSERT_EQUAL(NULL, usys_fopen("./test/data/nosuchfile.txt", "r"));
}

void test_usys_fopen_create_new_file(void) {
    TEST_ASSERT_NOT_EQUAL(NULL, usys_fopen("./test/data/newfile.txt", "w+"));
}

