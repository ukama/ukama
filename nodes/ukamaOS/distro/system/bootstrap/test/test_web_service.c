// tests/test_web_service.c
#include "unity.h"

#include <stdio.h>
#include <stdlib.h>
#include <ulfius.h>

#include "httpStatus.h"

#include "usys_error.h"
#include "usys_types.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_file.h"
#include "usys_string.h"
#include "usys_services.h"

#include "version.h"

typedef struct _u_instance UInst;

void test_start_web_services_should_return_true_for_valid_instance(void) {
    UInst instance;
    int result = start_web_services(&instance);
    TEST_ASSERT_TRUE(result);
}

int run_all_tests_web_service(void) {
    RUN_TEST(test_start_web_services_should_return_true_for_valid_instance);
}

