#include "unity.h"
#include "web_client.h"
#include "ulfius.h"

void test_wc_create_http_request(void) {
    URequest *req = wc_create_http_request("http://localhost", "test", "GET", NULL);

    TEST_ASSERT_NOT_NULL(req);
    TEST_ASSERT_EQUAL_STRING("http://localhost/test", req->http_url);
    TEST_ASSERT_EQUAL_STRING("GET", req->http_verb);

    ulfius_clean_request(req);
    free(req);
}

int run_all_tests_web_client(void) {

    RUN_TEST(test_wc_create_http_request);
}

