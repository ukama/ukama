#include "unity.h"
#include "gpsd.h"
#include "config.h"
#include <pthread.h>
#include <stdlib.h>
#include <string.h>

/* gpsd.c */
extern bool read_last_gps_data(char **lat, char **lon, char **gpsTime);

void test_read_last_gps_data_success(void) {
    char *lat = NULL;
    char *lon = NULL;
    char *gpsTime = NULL;

    FILE *file = fopen(GPS_LOC_FILE, "w");
    fprintf(file, "12.3456,78.9101,2026-08-23T22:15:30Z\n");
    fclose(file);

    bool result = read_last_gps_data(&lat, &lon, &gpsTime);

    TEST_ASSERT_TRUE(result);
    TEST_ASSERT_EQUAL_STRING("12.3456", lat);
    TEST_ASSERT_EQUAL_STRING("78.9101", lon);
    TEST_ASSERT_EQUAL_STRING("2026-08-23T22:15:30Z", gpsTime);

    free(lat);
    free(lon);
    free(gpsTime);
}

void test_read_last_gps_data_empty_file(void) {
    char *lat = NULL;
    char *lon = NULL;
    char *gpsTime = NULL;

    FILE *file = fopen(GPS_LOC_FILE, "w");
    fclose(file);

    bool result = read_last_gps_data(&lat, &lon, &gpsTime);
    TEST_ASSERT_FALSE(result);
    TEST_ASSERT_NULL(lat);
    TEST_ASSERT_NULL(lon);
    TEST_ASSERT_NULL(gpsTime);
}

int run_all_tests_gpsd(void) {

    RUN_TEST(test_read_last_gps_data_success);
    RUN_TEST(test_read_last_gps_data_empty_file);
}
