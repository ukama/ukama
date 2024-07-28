// tests/test_config.c
#include "unity.h"
#include "config.h"
#include "toml.h"

void test_process_config_file_should_return_false_for_null_config(void) {
    int result = process_config_file(NULL, NULL);
    TEST_ASSERT_FALSE(result);
}

void test_process_config_file_should_return_false_for_null_configData(void) {
    Config config;
    int result = process_config_file("config.toml", NULL);
    TEST_ASSERT_FALSE(result);
}

void test_process_config_file_should_return_true_for_valid_config(void) {
    Config config;
    int result = process_config_file("../../../../configs/capps/bootstrap/config.toml", &config);
    TEST_ASSERT_TRUE(result);
}

int run_all_tests_config(void) {
    RUN_TEST(test_process_config_file_should_return_false_for_null_config);
    RUN_TEST(test_process_config_file_should_return_false_for_null_configData);
    RUN_TEST(test_process_config_file_should_return_true_for_valid_config);
}
