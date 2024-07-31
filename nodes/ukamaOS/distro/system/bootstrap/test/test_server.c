// tests/test_server.c
#include "unity.h"
#include "server.h"

void test_send_request_to_init_should_return_false_for_invalid_params(void) {
    ServerInfo serverInfo;
    int result = send_request_to_init_with_exponential_backoff(NULL, 8080, NULL, &serverInfo);
    TEST_ASSERT_FALSE(result);
}

int run_all_tests_server(void) {
    RUN_TEST(test_send_request_to_init_should_return_false_for_invalid_params);
}
