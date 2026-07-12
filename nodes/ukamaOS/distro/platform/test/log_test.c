/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "test.h"

#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "unity.h"
#include "usys_log.h"

#define TEST_LOG_BUFFER 8192

static int capture_begin(int pipeFd[2]) {
    int savedStderr;

    if (pipe(pipeFd) != 0) {
        return -1;
    }

    savedStderr = dup(STDERR_FILENO);
    if (savedStderr < 0) {
        close(pipeFd[0]);
        close(pipeFd[1]);
        return -1;
    }

    if (dup2(pipeFd[1], STDERR_FILENO) < 0) {
        close(savedStderr);
        close(pipeFd[0]);
        close(pipeFd[1]);
        return -1;
    }

    close(pipeFd[1]);
    return savedStderr;
}

static ssize_t capture_end(int readFd, int savedStderr,
                           char *buffer, size_t bufferSize) {
    ssize_t bytes;

    if (dup2(savedStderr, STDERR_FILENO) < 0) {
        close(savedStderr);
        close(readFd);
        return -1;
    }
    close(savedStderr);

    bytes = read(readFd, buffer, bufferSize - 1);
    close(readFd);
    if (bytes >= 0) {
        buffer[bytes] = '\0';
    }

    return bytes;
}

void test_usys_log_emits_json_line(void) {
    char output[TEST_LOG_BUFFER] = {0};
    int pipeFd[2];
    int savedStderr;
    ssize_t bytes;

    usys_log_set_service("log-test");
    usys_log_set_level(USYS_LOG_TRACE);
    usys_log_set_quiet(false);

    savedStderr = capture_begin(pipeFd);
    TEST_ASSERT_GREATER_OR_EQUAL_INT(0, savedStderr);

    usys_log_info("quoted value=\"%s\"", "test");

    bytes = capture_end(pipeFd[0], savedStderr,
                        output, sizeof(output));
    TEST_ASSERT_GREATER_THAN_INT(0, bytes);
    TEST_ASSERT_EQUAL_CHAR('{', output[0]);
    TEST_ASSERT_EQUAL_CHAR('\n', output[bytes - 1]);
    TEST_ASSERT_NOT_NULL(strstr(output,
                                "\"schema\":\"ukama.log.v1\""));
    TEST_ASSERT_NOT_NULL(strstr(output, "\"level\":\"info\""));
    TEST_ASSERT_NOT_NULL(strstr(output, "\"app\":\"log-test\""));
    TEST_ASSERT_NOT_NULL(strstr(output,
                                "\"event\":\"log_message\""));
    TEST_ASSERT_NOT_NULL(strstr(output, "quoted value=\\\"test\\\""));
}

void test_usys_log_emits_structured_fields(void) {
    char output[TEST_LOG_BUFFER] = {0};
    USysLogField fields[] = {
        USYS_LOG_STR("dependency", "node_gateway"),
        USYS_LOG_STR("operation", "connect"),
        USYS_LOG_BOOL("retryable", true),
        USYS_LOG_INT("attempt", 2)
    };
    int pipeFd[2];
    int savedStderr;
    ssize_t bytes;

    usys_log_set_service("meshd");
    usys_log_set_level(USYS_LOG_TRACE);
    usys_log_set_quiet(false);

    savedStderr = capture_begin(pipeFd);
    TEST_ASSERT_GREATER_OR_EQUAL_INT(0, savedStderr);

    usys_log_event_error("backend", "dependency_connect_failed",
                         "Unable to connect to backend", fields,
                         USYS_LOG_FIELD_COUNT(fields));

    bytes = capture_end(pipeFd[0], savedStderr,
                        output, sizeof(output));
    TEST_ASSERT_GREATER_THAN_INT(0, bytes);
    TEST_ASSERT_NOT_NULL(strstr(output,
                                "\"component\":\"backend\""));
    TEST_ASSERT_NOT_NULL(strstr(
        output, "\"event\":\"dependency_connect_failed\""));
    TEST_ASSERT_NOT_NULL(strstr(
        output, "\"dependency\":\"node_gateway\""));
    TEST_ASSERT_NOT_NULL(strstr(output, "\"retryable\":true"));
    TEST_ASSERT_NOT_NULL(strstr(output, "\"attempt\":2"));
}
