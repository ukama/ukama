// tests/test_runner.c
#include "unity.h"
#include "gpsd.h"

void run_all_tests_config(void);
void run_all_tests_jserdes(void);
void run_all_tests_web_service(void);
void run_all_tests_web_client(void);
void run_all_tests_gpsd(void);

/* Mock external variables */
GPSData *gData;

/* Setup function */
void setUp(void) {
    gData = (GPSData *)calloc(1, sizeof(GPSData));
    pthread_mutex_init(&gData->mutex, NULL);
}

void tearDown(void) {
    pthread_mutex_destroy(&gData->mutex);
    free(gData->latitude);
    free(gData->longitude);
    free(gData);
}

int main(void) {

    UNITY_BEGIN();
    run_all_tests_config();
    run_all_tests_jserdes();
    //    run_all_tests_web_service();
    // run_all_tests_web_client();
    run_all_tests_gpsd();

    return UNITY_END();
}

