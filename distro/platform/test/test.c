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
#include "usys_string.h"

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

void test_usys_foperations(void) {
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

