/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */
#include <ulfius.h>
#include <pthread.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <unistd.h>
#include <errno.h>
#include <curl/curl.h>
#include <curl/easy.h>

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_file.h"
#include "usys_services.h"

#include "http_status.h"
#include "rlogd.h"

extern ThreadData *gData;

struct Response {
    char *buffer;
    size_t size;
};

/* Mutexs to ensure thread-safe writes for various dest */
static pthread_mutex_t logFileMutex = PTHREAD_MUTEX_INITIALIZER;
static pthread_mutex_t stdoutMutex  = PTHREAD_MUTEX_INITIALIZER;
static pthread_mutex_t stderrMutex  = PTHREAD_MUTEX_INITIALIZER;

static void write_to_log_file(const char *buffer);
static void write_to_stdout(const char *buffer);
static void write_to_stderr(const char *buffer);

static int find_json_buffer_size(json_t *json) {

    char *jStr = NULL;
    int  len = 0;
    
    jStr = json_dumps(json, JSON_ENCODE_ANY);
    if (jStr) {
        len = strlen(jStr);
        usys_free(jStr);
    }

    return len;
}

static int log_level(char *slevel) {

    if (!strcmp(slevel, "DEBUG")) {
        return USYS_LOG_DEBUG;
    } else if (!strcmp(slevel, "INFO")) {
        return USYS_LOG_INFO;
    } else if (!strcmp(slevel, "ERROR")) {
        return USYS_LOG_ERROR;
    } else {
        return USYS_LOG_TRACE;
    }
}


static size_t response_callback(void *contents,
                                size_t size,
                                size_t nmemb,
                                void *userp) {

    size_t realsize = size * nmemb;
    struct Response *response = (struct Response *)userp;

    response->buffer = realloc(response->buffer, response->size + realsize + 1);

    if(response->buffer == NULL) {
        usys_log_error("Not enough memory to realloc of size: %d",
                       response->size + realsize + 1);
        return 0;
    }

    memcpy(&(response->buffer[response->size]), contents, realsize);
    response->size += realsize;
    response->buffer[response->size] = 0;

    return realsize;
}

static long send_request_to_server(char *url,
                                   const char *data,
                                   struct Response *response) {

    long resCode=0;
    CURL *curl=NULL;
    CURLcode res;
    struct curl_slist *headers=NULL;

    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (curl == NULL) return resCode;

    response->buffer = malloc(1);
    response->size   = 0;

    headers = curl_slist_append(headers, "Accept: application/json");
    headers = curl_slist_append(headers, "charset: utf-8");

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "POST");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)response);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "rlog/0.1");
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, data);

    res = curl_easy_perform(curl);

    if (res != CURLE_OK) {
        usys_log_error("Error sending request to server at URL %s: %s", url,
                       curl_easy_strerror(res));
        usys_free(response->buffer);
    } else {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &resCode);
    }

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    curl_global_cleanup();

    return resCode;
}

static void write_to_ukama_service(char *nodeID, const char *log) {

    time_t now;
    int port = 0;
    char url[MAX_URL_LEN] = {0};
    char appName[MAX_SIZE] = {0};
    char logTime[9] = {0};
    char logLevel[16] = {0};
    char sourceFile[MAX_SIZE] = {0};
    char message[MAX_MSG_LEN] = {0};
    int sourceLine = 0;
    json_t *jElem  = NULL;
    json_t *jArray = NULL;
    char *buffer = NULL;

    struct Response response;

    time(&now);

    /* if buffer overflow or its been FLUSH_TIME. */
    if ((gData->bufferSize + strlen(log) + 1 > MAX_LOG_BUFFER) ||
        (difftime(now, gData->lastWriteTime) >= gData->flushTime &&
         gData->bufferSize)) {

        /* convert json to str and write to the desired output */
        buffer = json_dumps(gData->jOutputBuffer, JSON_ENCODE_ANY);
        if (buffer == NULL) return;

        port = usys_find_service_port(SERVICE_REMOTE);
        if (port == 0) {
            usys_log_error("Unable to get mesh.d port via services db");
            return;
        }

        sprintf(url, "http://localhost:%d/%s/logger/node/%s",
                port, API_VERSION, nodeID);

        if (send_request_to_server(&url[0], buffer, &response) != HttpStatus_OK) {
            /* fall back to log-file */
            write_to_log_file(log);
        }

        json_decref(gData->jOutputBuffer);
        gData->bufferSize    = 0;
        gData->jOutputBuffer = json_pack("{s:[]}", JTAG_LOGS);

        usys_free(buffer);
        usys_free(response.buffer);
        time(&gData->lastWriteTime);
    }

    /* <app_name> <time> <level> <file_name:line_number> <message> */
    if (sscanf(log, LOG_FORMAT, appName, logTime, logLevel,
               sourceFile, &sourceLine, message) != LOG_ELEMENTS) {
        usys_log_debug("Invalid log message: %s", log);
        return;
    }

    /* Add the new log if only it's log-level matches with set log-level */
    if (log_level(logLevel) >= gData->level) {
        jArray = json_object_get(gData->jOutputBuffer, JTAG_LOGS);
        if (jArray) {
            jElem = json_object();
            json_object_set_new(jElem, JTAG_APP_NAME, json_string(appName));
            json_object_set_new(jElem, JTAG_TIME,     json_string(logTime));
            json_object_set_new(jElem, JTAG_LEVEL,    json_string(logLevel));
            json_object_set_new(jElem, JTAG_MESSAGE,  json_string(message));

            /* add the element to array and update size. */
            json_array_append_new(jArray, jElem);
            gData->bufferSize = find_json_buffer_size(gData->jOutputBuffer);
        }
    }
}

static void write_to_log_file(const char *buffer) {

    FILE *fPtr = NULL;

    if (buffer == NULL) return;

    pthread_mutex_lock(&logFileMutex);

    fPtr = fopen(DEF_LOG_FILE, "a+");
    if (fPtr == NULL) {
        usys_log_error("Unable to open file: %s error: %s",
                       DEF_LOG_FILE,
                       strerror(errno));
        return;
    } else {
        fputs(buffer, fPtr);
        fclose(fPtr);
    }

    pthread_mutex_unlock(&logFileMutex);

    return;
}

static void write_to_stdout(const char *buffer) {

    if (buffer == NULL) return;

    pthread_mutex_lock(&stdoutMutex);
    printf("%s\n", buffer);
    pthread_mutex_unlock(&stdoutMutex);
}

static void write_to_stderr(const char *buffer) {

    if (buffer == NULL) return;

    pthread_mutex_lock(&stderrMutex);
    printf("%s\n", buffer);
    pthread_mutex_unlock(&stderrMutex);
}

void process_logs(void *nodeID, const char *log) {

    char appName[MAX_SIZE] = {0};
    char logTime[9] = {0};
    char logLevel[16] = {0};
    char sourceFile[MAX_SIZE] = {0};
    char message[MAX_MSG_LEN] = {0};
    int sourceLine = 0;

    if (strlen(log) > MAX_LOG_LEN) return;

    /* <app_name> <time> <level> <file_name:line_number> <message> */
    if (sscanf(log, LOG_FORMAT, appName, logTime, logLevel,
               sourceFile, &sourceLine, message) != LOG_ELEMENTS) {
        usys_log_debug("Invalid log message: %s", log);
        return;
    }

    if (log_level(logLevel) >= gData->level) {

        switch (gData->output) {
        case STDOUT:
            write_to_stdout(log);
            break;
        case STDERR:
            write_to_stderr(log);
            break;
        case LOG_FILE:
            write_to_log_file(log);
            break;
        case UKAMA_SERVICE:
            write_to_ukama_service((char *)nodeID, log);
            break;
        default:
            break;
        }
    }

    return;
}
