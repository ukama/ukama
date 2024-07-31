// tests/test_main.c
#include "unity.h"
#include "config.h"

void test_process_config_file_should_return_false_for_invalid_file(void) {
    Config config;
    int result = process_config_file("invalid_file.toml", &config);
    TEST_ASSERT_FALSE(result);
}

void test_process_config_file_should_return_true_for_valid_file(void) {
    Config config;
    int result = process_config_file("../../../../configs/capps/bootstrap/config.toml",
                                     &config);
    TEST_ASSERT_TRUE(result);
}

int run_all_tests_main(void) {
    RUN_TEST(test_process_config_file_should_return_false_for_invalid_file);
    RUN_TEST(test_process_config_file_should_return_true_for_valid_file);
}
