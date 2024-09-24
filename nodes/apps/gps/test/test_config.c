#include "unity.h"
#include "config.h"

void test_config_init(void) {
    Config config;
    config.serviceName = "GPS_Service";
    
    TEST_ASSERT_EQUAL_STRING("GPS_Service", config.serviceName);
}

void run_all_tests_config(void) {
    RUN_TEST(test_config_init);
}
