#include "unity.h"
#include "json_types.h"
#include <jansson.h>
#include "usys_types.h"

void test_json_deserialize_node_id_success(void) {
    json_t *json = json_object();
    json_object_set_new(json, "UUID", json_string("mock-node-id"));

    char *nodeID = NULL;
    bool result = json_deserialize_node_id(&nodeID, json);

    TEST_ASSERT_TRUE(result);
    TEST_ASSERT_EQUAL_STRING("mock-node-id", nodeID);

    free(nodeID);
    json_decref(json);
}

void test_json_deserialize_node_id_failure(void) {
    json_t *json = json_object(); // Missing UUID key

    char *nodeID = NULL;
    bool result = json_deserialize_node_id(&nodeID, json);

    TEST_ASSERT_FALSE(result);
    TEST_ASSERT_NULL(nodeID);

    json_decref(json);
}

void run_all_tests_jserdes(void) {

    RUN_TEST(test_json_deserialize_node_id_success);
    RUN_TEST(test_json_deserialize_node_id_failure);
}

