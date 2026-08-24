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
#include <string.h>
#include <errno.h>
#include <unistd.h>
#include <sys/wait.h>

#include "gpsd.h"
#include "config.h"

#include "usys_log.h"
#include "usys_mem.h"

#include "static.h"

/* main.c */
extern GPSData *gData;

static void gps_mark_unlocked(void) {

    if (gData == NULL) return;

    pthread_mutex_lock(&gData->mutex);
    if (gData->gpsLock || gData->lockLostAt == 0) {
        gData->lockLostAt = time(NULL);
    }
    gData->gpsLock = USYS_FALSE;
    pthread_mutex_unlock(&gData->mutex);
}

STATIC bool read_last_gps_data(char **lat, char **lon, char **gpsTime) {

    FILE *file = NULL;
    char line[MAX_LINE_LENGTH]      = {0};
    char lastLat[MAX_LAT_LONG_SIZE] = {0};
    char lastLon[MAX_LAT_LONG_SIZE] = {0};
    char lastTime[MAX_GPS_TIME_SIZE] = {0};

    if (lat == NULL || lon == NULL || gpsTime == NULL) {
        return USYS_FALSE;
    }

    *lat = NULL;
    *lon = NULL;
    *gpsTime = NULL;

    file = fopen(GPS_LOC_FILE, "r");
    if (!file) {
        usys_log_error("Unable to open GPS coordinate file: %s Error: %s",
                       GPS_LOC_FILE, strerror(errno));
        return USYS_FALSE;
    }

    /* Reading file line by line and only storing the last entry */
    while (fgets(line, sizeof(line), file)) {
        /* Parse latitude, longitude and GPS time from the line. */
        if (sscanf(line,
                   "%31[^,],%31[^,],%31s",
                   lastLat,
                   lastLon,
                   lastTime) == 3) {
            /* Continue updating the values with each line
             * (only the last one will remain) 
             */
        } else {
            usys_log_error("Error parsing line: %s", line);
        }
    }

    fclose(file);
    
    if (lastLat[0] == '\0' ||
        lastLon[0] == '\0' ||
        lastTime[0] == '\0') {
        return USYS_FALSE;
    }

    *lat = strdup(lastLat);
    *lon = strdup(lastLon);
    *gpsTime = strdup(lastTime);

    if (*lat == NULL || *lon == NULL || *gpsTime == NULL) {
        usys_free(*lat);
        usys_free(*lon);
        usys_free(*gpsTime);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

STATIC bool gps_data_collection_and_processing_thread(Config *config) {

    int ret;
    char runMe[MAX_BUFFER] = {0};
    char *lat = NULL, *lon = NULL, *gpsTime = NULL;

    pthread_setcancelstate(PTHREAD_CANCEL_ENABLE,  NULL);
    pthread_setcanceltype(PTHREAD_CANCEL_DEFERRED, NULL);

    while (USYS_TRUE) {

        sleep(GPS_WAIT_TIME);

        /* get gps data from the trx board */
        snprintf(runMe, MAX_BUFFER, "%s get_gps_data %s",
                GPS_SCRIPT,
                config->gpsHost);

        ret = system(runMe);
        if (ret == -1 ||
            !WIFEXITED(ret) ||
            WEXITSTATUS(ret) != 0) {
            gps_mark_unlocked();
            continue;
        }

        /* see if gps is locked */
        snprintf(runMe, MAX_BUFFER, "%s gps_fix", GPS_SCRIPT);
        ret = system(runMe);
        if (ret != -1 &&
            WIFEXITED(ret) &&
            WEXITSTATUS(ret) == 0) {
            /* gps is locked, get coordinates */
            snprintf(runMe, MAX_BUFFER, "%s get_coordinates", GPS_SCRIPT);
            ret = system(runMe);

            if (ret != -1 &&
                WIFEXITED(ret) &&
                WEXITSTATUS(ret) == 0) {
                /* read /tmp/gps_loc.log file: lat,lon,time */
                if (read_last_gps_data(&lat, &lon, &gpsTime)) {

                    if (gData == NULL) continue;

                    /* Update the complete GPS snapshot atomically. */
                    pthread_mutex_lock(&gData->mutex);

                    gData->gpsLock  = USYS_TRUE;
                    usys_free(gData->latitude);
                    usys_free(gData->longitude);
                    usys_free(gData->time);
                    gData->latitude  = lat;
                    gData->longitude = lon;
                    gData->time      = gpsTime;
                    gData->lastLockAt = time(NULL);
                    gData->lockLostAt = 0;

                    lat = NULL;
                    lon = NULL;
                    gpsTime = NULL;

                    pthread_mutex_unlock(&gData->mutex);
                } else {
                    gps_mark_unlocked();
                }
            } else {
                gps_mark_unlocked();
                continue;
            }
        } else {
            gps_mark_unlocked();
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

    if (config == NULL || tid == NULL) return USYS_FALSE;

    ret = pthread_create(tid, NULL, gps_thread_wrapper, (void*) config);
    if (ret != 0) {
        usys_log_error("Failed to create GPS thread");
        return USYS_FALSE;
    }

    return USYS_TRUE;
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
