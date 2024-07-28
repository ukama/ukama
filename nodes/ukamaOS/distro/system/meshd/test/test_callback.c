#include <ulfius.h>
#include <pthread.h>

#include "unity.h"
#include "callback.h"
#include "mesh.h"
#include "work.h"
#include "jserdes.h"
#include "httpStatus.h"
#include "version.h"
#include "map.h"

WorkList *Transmit = NULL;
MapTable *ClientTable = NULL;
pthread_mutex_t mutex;
pthread_cond_t hasData;
char *queue = NULL;

// Test for callback_websocket function
void test_callback_websocket(void) {
    URequest request;
    UResponse response;
    Config config;
    int ret;

    memset(&request, 0, sizeof(request));
    memset(&response, 0, sizeof(response));
    memset(&config, 0, sizeof(config));

    ret = callback_websocket(&request, &response, &config);
    TEST_ASSERT_EQUAL_INT(U_CALLBACK_ERROR, ret);
}

// Test for callback_not_allowed function
void test_callback_not_allowed(void) {
    URequest request;
    UResponse response;
    memset(&request, 0, sizeof(request));
    memset(&response, 0, sizeof(response));

    int ret = callback_not_allowed(&request, &response, NULL);
    TEST_ASSERT_EQUAL_INT(U_CALLBACK_CONTINUE, ret);
    TEST_ASSERT_EQUAL_STRING(HttpStatusStr(HttpStatus_Forbidden), response.binary_body);
}

// Test for callback_forward_service function
void test_callback_forward_service(void) {
    URequest request;
    UResponse response;
    Config config;
    memset(&request, 0, sizeof(request));
    memset(&response, 0, sizeof(response));
    memset(&config, 0, sizeof(config));

    char url[] = "http://localhost:8080";
    char service[] = "test_service";
    u_map_put(request.map_header, "Host", url);
    u_map_put(request.map_header, "User-Agent", service);

    int ret = callback_forward_service(&request, &response, &config);
    TEST_ASSERT_EQUAL_INT(U_CALLBACK_CONTINUE, ret);
}

// Test for web_service_cb_ping function
void test_web_service_cb_ping(void) {
    URequest request;
    UResponse response;
    memset(&request, 0, sizeof(request));
    memset(&response, 0, sizeof(response));

    int ret = web_service_cb_ping(&request, &response, NULL);
    TEST_ASSERT_EQUAL_INT(U_CALLBACK_CONTINUE, ret);
    TEST_ASSERT_EQUAL_STRING(HttpStatusStr(HttpStatus_OK), response.binary_body);
}

// Test for web_service_cb_version function
void test_web_service_cb_version(void) {
    URequest request;
    UResponse response;
    memset(&request, 0, sizeof(request));
    memset(&response, 0, sizeof(response));

    int ret = web_service_cb_version(&request, &response, NULL);
    TEST_ASSERT_EQUAL_INT(U_CALLBACK_CONTINUE, ret);
    TEST_ASSERT_EQUAL_STRING(VERSION, response.binary_body);
}

// Test for web_service_cb_default function
void test_web_service_cb_default(void) {
    URequest request;
    UResponse response;
    memset(&request, 0, sizeof(request));
    memset(&response, 0, sizeof(response));

    int ret = web_service_cb_default(&request, &response, NULL);
    TEST_ASSERT_EQUAL_INT(U_CALLBACK_CONTINUE, ret);
    TEST_ASSERT_EQUAL_STRING(HttpStatusStr(HttpStatus_NotFound), response.binary_body);
}

void run_all_tests_callback(void) {
    RUN_TEST(test_callback_websocket);
    RUN_TEST(test_callback_not_allowed);
    RUN_TEST(test_web_service_cb_ping);
    RUN_TEST(test_web_service_cb_version);
    RUN_TEST(test_web_service_cb_default);
}
