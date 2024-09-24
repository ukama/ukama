#include "unity.h"
#include "gpsd.h"
#include "config.h"
#include <pthread.h>
#include <stdlib.h>
#include <string.h>

/* gpsd.c */
extern bool read_last_lat_long(char **lat, char **lon);

void test_read_last_lat_long_success(void) {
    char *lat = NULL;
    char *lon = NULL;

    FILE *file = fopen("test_gps_loc.log", "w");
    fprintf(file, "12.3456,78.9101\n");
    fclose(file);

    bool result = read_last_lat_long(&lat, &lon);

    TEST_ASSERT_TRUE(result);
    TEST_ASSERT_EQUAL_STRING("12.3456", lat);
    TEST_ASSERT_EQUAL_STRING("78.9101", lon);

    free(lat);
    free(lon);
}

void test_read_last_lat_long_empty_file(void) {
    char *lat = NULL;
    char *lon = NULL;

    FILE *file = fopen("test_gps_loc.log", "w");
    fclose(file);

    bool result = read_last_lat_long(&lat, &lon);
    TEST_ASSERT_FALSE(result);
    TEST_ASSERT_NULL(lat);
    TEST_ASSERT_NULL(lon);
}

int run_all_tests_gpsd(void) {

    RUN_TEST(test_read_last_lat_long_success);
    RUN_TEST(test_read_last_lat_long_empty_file);
}

