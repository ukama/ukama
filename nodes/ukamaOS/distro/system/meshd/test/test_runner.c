#include "unity.h"

extern void run_all_tests_config(void);
extern void run_all_tests_callback(void);
extern void run_all_tests_jserdes(void);
extern void run_all_tests_main(void);
extern void run_all_tests_map(void);
extern void run_all_tests_websocket(void);
extern void run_all_tests_work(void);

void setUp(void) {
    // Set up code (if needed)
}

void tearDown(void) {
    // Tear down code (if needed)
}

int main(void) {
    UNITY_BEGIN();
    run_all_tests_config();
    run_all_tests_callback();
    run_all_tests_jserdes();
    run_all_tests_main();
    run_all_tests_websocket();
    run_all_tests_work();
    return UNITY_END();
}
