/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <stdlib.h>
#include <stdio.h>
#include <pthread.h>

#include "gpsd.h"
#include "config.h"

#include "usys_log.h"
#include "usys_mem.h"

#include "static.h"

/* main.c */
extern GPSData *gData;

STATIC bool read_last_lat_long(char **lat, char **lon) {

    FILE *file = NULL;
    char line[MAX_LINE_LENGTH]      = {0};
    char lastLat[MAX_LAT_LONG_SIZE] = {0};
    char lastLon[MAX_LAT_LONG_SIZE] = {0};

    file = fopen(GPS_LOC_FILE, "r");
    if (!file) {
        usys_log_error("Unable to open GPS coordinate file: %s Error: %s",
                       GPS_LOC_FILE, strerror(errno));
        return USYS_FALSE;
    }

    /* Reading file line by line and only storing the last entry */
    while (fgets(line, sizeof(line), file)) {
        /* Parse the latitude and longitude from the line */
        if (sscanf(line, "%[^,],%s", lastLat, lastLon) == 2) {
            /* Continue updating the lat, lon with each line
             * (only the last one will remain) 
             */
        } else {
            usys_log_error("Error parsing line: %s", line);
        }
    }

    fclose(file);
    
    if (strlen(lastLat) == 0 && strlen(lastLon) == 0) {

        *lat = NULL;
        *lon = NULL;
        return USYS_FALSE;
    }

    *lat = strdup(lastLat);
    *lon = strdup(lastLon);

    return USYS_TRUE;
}

STATIC bool gps_data_collection_and_processing_thread(Config *config) {

    int ret;
    char runMe[MAX_BUFFER] = {0};
    char *lat = NULL, *lon = NULL;

    pthread_setcancelstate(PTHREAD_CANCEL_ENABLE,  NULL);
    pthread_setcanceltype(PTHREAD_CANCEL_DEFERRED, NULL);

    while (USYS_TRUE) {

        sleep(GPS_WAIT_TIME);

        /* get gps data from the trx board */
        snprintf(runMe, MAX_BUFFER, "%s get_gps_data %s",
                GPS_SCRIPT,
                config->gpsHost);

        ret = system(runMe);
        if (WIFEXITED(ret) && WEXITSTATUS(ret) != 0) {
            continue;
        }

        /* see if gps is locked */
        snprintf(runMe, MAX_BUFFER, "%s gps_fix", GPS_SCRIPT);
        ret = system(runMe);
        if (WIFEXITED(ret) && WEXITSTATUS(ret) == 0) {
            /* gps is locked, get coordinates */
            snprintf(runMe, MAX_BUFFER, "%s get_coordinates", GPS_SCRIPT);
            ret = system(runMe);

            if (WIFEXITED(ret) && WEXITSTATUS(ret) == 0) {
                /* read /tmp/gps_loc.log file: lat,lon */
                if (read_last_lat_long(&lat, &lon)) {

                    if (gData == NULL) continue;

                    /* update record */
                    pthread_mutex_lock(&gData->mutex);

                    gData->gpsLock  = USYS_TRUE;
                    usys_free(gData->latitude);
                    usys_free(gData->longitude);
                    gData->latitude  = strdup(lat);
                    gData->longitude = strdup(lon);

                    pthread_mutex_unlock(&gData->mutex);

                    usys_free(lat);
                    usys_free(lon);
                }
            } else {
                continue;
            }
        } else {
            continue;
        }
    }

    usys_log_debug("GPS thread existing.");
    return USYS_TRUE;
}

STATIC void *gps_thread_wrapper(void* arg) {

    Config* config = (Config*) arg;
    gps_data_collection_and_processing_thread(config);

    return NULL;
}

bool start_gps_data_collection_and_processing(Config *config, pthread_t *tid) {

    int ret = 0;

    pthread_create(tid, NULL, gps_thread_wrapper, (void*) config);
    if (ret != 0) {
        usys_log_error("Failed to create GPS thread");
        return USYS_FALSE;
    }

    return USYS_FALSE;
}

void stop_gps_data_collection_and_processing(pthread_t tid) {

    if (tid == 0) return;

    pthread_cancel(tid);
    pthread_join(tid, NULL);

    if (remove(GPS_LOC_FILE) != 0) {
        usys_log_error("Error deleting %s Error: %s",
                       GPS_LOC_FILE,
                       strerror(errno));
    }

    if (remove(GPS_RAW_FILE) != 0) {
        usys_log_error("Error deleting %s Error: %s",
                       GPS_RAW_FILE,
                       strerror(errno));
    }

    pthread_mutex_lock(&gData->mutex);
    gData->gpsLock = USYS_FALSE;
    pthread_mutex_unlock(&gData->mutex);
}
