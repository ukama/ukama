// tests/test_mesh_config.c
#include "unity.h"
#include "mesh_config.h"

void test_read_mesh_config_file_should_return_false_for_invalid_file(void) {
    MeshConfig meshConfig;
    int result = read_mesh_config_file("invalid_file.toml", &meshConfig);
    TEST_ASSERT_FALSE(result);
}

void test_read_mesh_config_file_should_return_true_for_valid_file(void) {
    MeshConfig meshConfig;
    int result = read_mesh_config_file("../../../../configs/capps/mesh/config.toml",
                                       &meshConfig);
    TEST_ASSERT_TRUE(result);
}

int run_all_tests_mesh_config(void) {
    RUN_TEST(test_read_mesh_config_file_should_return_false_for_invalid_file);
    RUN_TEST(test_read_mesh_config_file_should_return_true_for_valid_file);
}
