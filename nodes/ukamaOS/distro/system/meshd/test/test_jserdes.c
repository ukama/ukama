#include "unity.h"
#include "jserdes.h"
#include "mesh.h"
#include <jansson.h>

/* jserdes.c */
extern void add_map_to_request(json_t **json, UMap *map, int mapType);
extern void serialize_message_data(URequest *request, char **data);
extern int deserialize_service_info(ServiceInfo **service, json_t *json);
extern void deserialize_map_array(UMap **map, json_t *json);
extern void deserialize_map(URequest **request, json_t *json);

// Test for serialize_device_info function
void test_serialize_device_info(void) {
    NodeInfo device = { "node_id" };
    json_t *json = NULL;

    int ret = serialize_device_info(&json, &device);
    TEST_ASSERT_EQUAL_INT(TRUE, ret);
    TEST_ASSERT_NOT_NULL(json);

    json_decref(json);
}

// Test for serialize_local_service_response function
void test_serialize_local_service_response(void) {
    Message message = { "seqNo" };
    char *response = NULL;

    int ret = serialize_local_service_response(&response, &message, 200, strlen("data"), "data");
    TEST_ASSERT_EQUAL_INT(TRUE, ret);
    TEST_ASSERT_NOT_NULL(response);

    free(response);
}

// Test for serialize_websocket_message function
void test_serialize_websocket_message(void) {
    URequest request;
    char *nodeID = "nodeID";
    char *port = "port";
    char *agent = "agent";
    char *sourcePort = "sourcePort";
    char *str = NULL;

    memset(&request, 0, sizeof(request));

    int ret = serialize_websocket_message(&str, &request, nodeID, port, agent, sourcePort);
    TEST_ASSERT_EQUAL_INT(TRUE, ret);
    TEST_ASSERT_NOT_NULL(str);

    free(str);
}

// Test for deserialize_node_info function
void test_deserialize_node_info(void) {
    json_t *json = json_pack("{s:s, s:s}", "node_id", "node_id_value", "port", "port_value");
    NodeInfo *node = NULL;

    int ret = deserialize_node_info(&node, json);
    TEST_ASSERT_EQUAL_INT(TRUE, ret);
    TEST_ASSERT_NOT_NULL(node);
    TEST_ASSERT_EQUAL_STRING("node_id_value", node->nodeID);
    TEST_ASSERT_EQUAL_STRING("port_value", node->port);

    free(node->nodeID);
    free(node->port);
    free(node);
    json_decref(json);
}

// Test for deserialize_request_info function
void test_deserialize_request_info(void) {
    char *str = "{\"protocol\": \"http\", \"method\": \"GET\", \"url\": \"/test\"}";
    URequest *request = NULL;

    int ret = deserialize_request_info(&request, str);
    TEST_ASSERT_EQUAL_INT(TRUE, ret);
    TEST_ASSERT_NOT_NULL(request);

    free(request->http_protocol);
    free(request->http_verb);
    free(request->http_url);
    free(request);
}

void run_all_tests_jserdes(void) {
    RUN_TEST(test_serialize_device_info);
    RUN_TEST(test_serialize_local_service_response);
    RUN_TEST(test_serialize_websocket_message);
    RUN_TEST(test_deserialize_node_info);
    RUN_TEST(test_deserialize_request_info);
}

