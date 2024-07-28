#include "unity.h"
#include "data.h"
#include "mesh.h"
#include "work.h"
#include "jserdes.h"
#include "httpStatus.h"
#include <curl/curl.h>
#include <ulfius.h>

/* data.c */
extern size_t response_callback(void *contents, size_t size, size_t nmemb, void *userp);
extern void find_service_name_and_ep(char *input, char **name, char **ep);
extern int send_data_to_local_service(URequest *data, char *hostname, int *httpStatus, char **retStr);

typedef struct _response {
        char *buffer;
        size_t size;
} Response;

// Test for response_callback function
void test_response_callback(void) {
    Response response = {NULL, 0};
    char data[] = "response data";
    size_t ret;

    ret = response_callback(data, sizeof(char), strlen(data), &response);
    TEST_ASSERT_EQUAL(strlen(data), ret);
    TEST_ASSERT_EQUAL_STRING(data, response.buffer);

    free(response.buffer);
}

// Test for clear_request function
void test_clear_request(void) {
    MRequest *request = (MRequest *)malloc(sizeof(MRequest));
    request->reqType = strdup("type");
    request->deviceInfo = strdup("device");
    request->serviceInfo = strdup("service");
    request->requestInfo = (URequest *)malloc(sizeof(URequest));

    clear_request(&request);

    TEST_ASSERT_NULL(request);
}

// Test for find_service_name_and_ep function
void test_find_service_name_and_ep(void) {
    char input[] = "/service/endpoint";
    char *name = NULL;
    char *ep = NULL;

    find_service_name_and_ep(input, &name, &ep);

    TEST_ASSERT_EQUAL_STRING("service", name);
    TEST_ASSERT_EQUAL_STRING("/endpoint", ep);

    free(name);
    free(ep);
}

// Test for send_data_to_local_service function
void test_send_data_to_local_service(void) {
    URequest request;
    char *hostname = "localhost";
    int httpStatus;
    char *retStr = NULL;

    memset(&request, 0, sizeof(request));

    int ret = send_data_to_local_service(&request, hostname, &httpStatus, &retStr);

    TEST_ASSERT_EQUAL_INT(TRUE, ret);
    TEST_ASSERT_NOT_NULL(retStr);

    free(retStr);
}

void run_all_tests_data(void) {
    RUN_TEST(test_response_callback);
    RUN_TEST(test_clear_request);
    RUN_TEST(test_find_service_name_and_ep);
    RUN_TEST(test_send_data_to_local_service);
}

