// tests/test_jserdes.c
#include "unity.h"
#include "jserdes.h"
#include "nodeInfo.h"

void test_deserialize_node_info_should_return_false_for_null_json(void) {
    NodeInfo *nodeInfo;
    int result = deserialize_node_info(&nodeInfo, NULL);
    TEST_ASSERT_FALSE(result);
}

void test_deserialize_node_info_should_return_false_for_invalid_json(void) {
    NodeInfo *nodeInfo;
    json_t *json = json_object();
    int result = deserialize_node_info(&nodeInfo, json);
    TEST_ASSERT_FALSE(result);
}

void test_deserialize_node_info_should_return_true_for_valid_json(void) {
    NodeInfo *nodeInfo;
    const char *json_str = "{\"nodeInfo\": {\"UUID\": \"ukma-7001-tnode-sa03-1100\"}}";
    json_t *json = json_loads(json_str, 0, NULL);

    int result = deserialize_node_info(&nodeInfo, json);
    TEST_ASSERT_TRUE(result);
    TEST_ASSERT_EQUAL_STRING("ukma-7001-tnode-sa03-1100", nodeInfo->uuid);

    json_decref(json);
}

int run_all_tests_jserdes(void) {
    RUN_TEST(test_deserialize_node_info_should_return_false_for_null_json);
    RUN_TEST(test_deserialize_node_info_should_return_false_for_invalid_json);
    RUN_TEST(test_deserialize_node_info_should_return_true_for_valid_json);
}
