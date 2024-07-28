// tests/test_runner.c
#include "unity.h"

void run_all_tests_config(void);
void run_all_tests_jserdes(void);
void run_all_tests_main(void);
void run_all_tests_mesh_config(void);
void run_all_tests_noded(void);
void run_all_tests_server(void);
void run_all_tests_web_service(void);

void setUp(void) {
    // Set up code (if needed)
}

void tearDown(void) {
    // Tear down code (if needed)
}

int main(void) {
    UNITY_BEGIN();
    run_all_tests_config();
    run_all_tests_jserdes();
    run_all_tests_main();
    run_all_tests_mesh_config();
    run_all_tests_noded();
    run_all_tests_web_service();
    return UNITY_END();
}
