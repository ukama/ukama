#include "unity.h"
#include "web_service.h"
#include "http_status.h"
#include "ulfius.h"

void test_web_service_cb_ping(void) {
    UResponse response;
    ulfius_init_response(&response);
    
    int result = web_service_cb_ping(NULL, &response, NULL);
    
    TEST_ASSERT_EQUAL(U_CALLBACK_CONTINUE, result);
    TEST_ASSERT_EQUAL(HttpStatus_OK, response.status);

    ulfius_clean_response(&response);
}

void test_web_service_cb_version(void) {
    UResponse response;
    ulfius_init_response(&response);

    int result = web_service_cb_version(NULL, &response, NULL);

    TEST_ASSERT_EQUAL(U_CALLBACK_CONTINUE, result);
    TEST_ASSERT_EQUAL(HttpStatus_OK, response.status);

    ulfius_clean_response(&response);
}

int run_all_tests_web_service(void) {

    RUN_TEST(test_web_service_cb_ping);
    RUN_TEST(test_web_service_cb_version);
}

