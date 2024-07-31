// tests/test_noded.c
#include "unity.h"
#include "nodeInfo.h"

void test_get_nodeID_from_noded_should_return_false_for_invalid_host(void) {
    char *nodeID;
    int result = get_nodeID_from_noded(&nodeID, NULL, 8080);
    TEST_ASSERT_FALSE(result);
}

void test_get_nodeID_from_noded_should_return_false_for_invalid_port(void) {
    char *nodeID;
    int result = get_nodeID_from_noded(&nodeID, "localhost", 0);
    TEST_ASSERT_FALSE(result);
}

int run_all_tests_noded(void) {
    RUN_TEST(test_get_nodeID_from_noded_should_return_false_for_invalid_host);
    RUN_TEST(test_get_nodeID_from_noded_should_return_false_for_invalid_port);
}
