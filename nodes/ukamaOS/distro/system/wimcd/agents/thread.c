/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <dirent.h>
#include <errno.h>
#include <jansson.h>
#include <pthread.h>
#include <signal.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#include "agent.h"
#include "agent/jserdes.h"
#include "http_status.h"
#include "package_cache.h"
#include "wimc.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_types.h"

#define MAX_ARGS 10
#define CASYNC_RETRIES 3
#define CASYNC_TIMEOUT_SEC 900

static int mkdir_p(const char *path, mode_t mode) {

    char tmp[WIMC_MAX_PATH_LEN];
    char *p;

    if (path == NULL || *path == '\0') {
        return -1;
    }

    if (strlen(path) >= sizeof(tmp)) {
        return -1;
    }

    snprintf(tmp, sizeof(tmp), "%s", path);

    for (p = tmp + 1; *p != '\0'; p++) {
        if (*p == '/') {
            *p = '\0';
            if (mkdir(tmp, mode) != 0 && errno != EEXIST) {
                return -1;
            }
            *p = '/';
        }
    }

    if (mkdir(tmp, mode) != 0 && errno != EEXIST) {
        return -1;
    }

    return 0;
}

static int rm_rf(const char *path) {

    DIR *dir;
    struct dirent *ent;
    struct stat st;
    char child[WIMC_MAX_PATH_LEN];

    if (path == NULL || lstat(path, &st) != 0) {
        return 0;
    }

    if (!S_ISDIR(st.st_mode)) {
        return unlink(path);
    }

    dir = opendir(path);
    if (dir == NULL) {
        return -1;
    }

    while ((ent = readdir(dir)) != NULL) {
        if (strcmp(ent->d_name, ".") == 0 ||
            strcmp(ent->d_name, "..") == 0) {
            continue;
        }

        if (snprintf(child, sizeof(child), "%s/%s", path,
                     ent->d_name) >= (int)sizeof(child)) {
            closedir(dir);
            return -1;
        }

        if (rm_rf(child) != 0) {
            closedir(dir);
            return -1;
        }
    }

    closedir(dir);
    return rmdir(path);
}

static void copy_fetch_request(WFetch **dest, WFetch *src) {

    WFetch *df;
    WContent *content;

    if (dest == NULL || src == NULL || src->content == NULL) {
        return;
    }

    *dest = (WFetch *)calloc(1, sizeof(WFetch));
    if (*dest == NULL) {
        return;
    }

    df = *dest;
    uuid_copy(df->uuid, src->uuid);
    df->cbURL = src->cbURL ? strdup(src->cbURL) : NULL;
    df->interval = src->interval;

    df->content = (WContent *)calloc(1, sizeof(WContent));
    if (df->content == NULL) {
        goto fail;
    }

    content = df->content;
    content->name = src->content->name ? strdup(src->content->name) : NULL;
    content->tag = src->content->tag ? strdup(src->content->tag) : NULL;
    content->method = src->content->method ? strdup(src->content->method) :
                      NULL;
    content->indexURL = src->content->indexURL ?
                        strdup(src->content->indexURL) : NULL;
    content->storeURL = src->content->storeURL ?
                        strdup(src->content->storeURL) : NULL;

    return;

fail:
    usys_free(df->cbURL);
    usys_free(df);
    *dest = NULL;
}

void free_fetch_request(WFetch *ptr) {

    if (ptr == NULL) {
        return;
    }

    if (ptr->content != NULL) {
        usys_free(ptr->content->name);
        usys_free(ptr->content->tag);
        usys_free(ptr->content->method);
        usys_free(ptr->content->indexURL);
        usys_free(ptr->content->storeURL);
        usys_free(ptr->content);
    }

    usys_free(ptr->cbURL);
    usys_free(ptr);
}

static size_t response_callback(void *contents, size_t size, size_t nmemb,
                                void *userp) {

    (void)contents;
    (void)userp;
    return size * nmemb;
}

static long send_update_to_wimc(WFetch *fetch, TransferState state,
                                const char *voidStr) {

    AgentReq req;
    Update update;
    json_t *json = NULL;
    char *jsonStr = NULL;
    CURL *curl = NULL;
    CURLcode res;
    struct curl_slist *headers = NULL;
    long code = 0;

    if (fetch == NULL || fetch->cbURL == NULL) {
        return 0;
    }

    memset(&req, 0, sizeof(req));
    memset(&update, 0, sizeof(update));

    uuid_copy(update.uuid, fetch->uuid);
    update.totalKB = 0;
    update.transferKB = 0;
    update.transferState = state;
    update.voidStr = (char *)(voidStr ? voidStr : "");
    req.update = &update;

    if (!serialize_agent_request(&req, &json)) {
        return 0;
    }

    jsonStr = json_dumps(json, 0);
    if (jsonStr == NULL) {
        json_decref(json);
        return 0;
    }

    curl = curl_easy_init();
    if (curl == NULL) {
        usys_free(jsonStr);
        json_decref(json);
        return 0;
    }

    headers = curl_slist_append(headers, "Accept: application/json");
    headers = curl_slist_append(headers, "Content-Type: application/json");
    headers = curl_slist_append(headers, "charset: utf-8");

    curl_easy_setopt(curl, CURLOPT_URL, fetch->cbURL);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "PUT");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, jsonStr);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "wimc-agent/0.1");
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT, 5L);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 30L);

    res = curl_easy_perform(curl);
    if (res == CURLE_OK) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    } else {
        usys_log_error("Failed to update WIMC: %s", curl_easy_strerror(res));
    }

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    usys_free(jsonStr);
    json_decref(json);

    return code;
}

static int wait_for_child(pid_t pid, int timeoutSec) {

    int status;
    int elapsed;
    pid_t ret;

    elapsed = 0;
    while (1) {
        ret = waitpid(pid, &status, WNOHANG);
        if (ret == pid) {
            if (WIFEXITED(status) && WEXITSTATUS(status) == 0) {
                return 0;
            }
            return -1;
        }

        if (ret < 0) {
            return -1;
        }

        if (elapsed >= timeoutSec) {
            kill(pid, SIGTERM);
            sleep(2);
            kill(pid, SIGKILL);
            waitpid(pid, &status, 0);
            return -1;
        }

        sleep(1);
        elapsed++;
    }
}

static int run_casync(WFetch *fetch, const char *extractPath) {

    pid_t pid;
    char *args[MAX_ARGS];
    int i;

    memset(args, 0, sizeof(args));

    args[0] = strdup(AGENT_EXEC);
    args[1] = strdup("extract");
    args[2] = strdup(fetch->content->indexURL);
    args[3] = strdup("--store");
    args[4] = strdup(fetch->content->storeURL);
    args[5] = strdup(extractPath);
    args[6] = NULL;

    for (i = 0; i < 6; i++) {
        if (args[i] == NULL) {
            goto fail;
        }
    }

    pid = fork();
    if (pid < 0) {
        goto fail;
    }

    if (pid == 0) {
        execv(AGENT_EXEC, args);
        _exit(127);
    }

    for (i = 0; args[i] != NULL; i++) {
        free(args[i]);
    }

    return wait_for_child(pid, CASYNC_TIMEOUT_SEC);

fail:
    for (i = 0; args[i] != NULL; i++) {
        free(args[i]);
    }
    return -1;
}

static void *execute_agent(void *data) {

    WFetch *fetch;
    WContent *content;
    char uuidStr[36 + 1];
    char extractPath[WIMC_MAX_PATH_LEN];
    char uuidTmpDir[WIMC_MAX_PATH_LEN];
    char publishedPath[WIMC_MAX_PATH_LEN];
    char actualVersion[WIMC_MAX_NAME_LEN];
    int attempt;
    int ok;

    fetch = (WFetch *)data;
    if (fetch == NULL || fetch->content == NULL) {
        pthread_exit(NULL);
    }

    content = fetch->content;
    uuid_unparse(fetch->uuid, uuidStr);

    if (pkg_ensure_cache_dirs() != 0) {
        send_update_to_wimc(fetch, ERR, "failed to create package cache dirs");
        goto done;
    }

    if (snprintf(uuidTmpDir, sizeof(uuidTmpDir), "%s/.tmp/%s",
                 DEFAULT_APPS_PKGS_PATH, uuidStr) >=
        (int)sizeof(uuidTmpDir)) {
        send_update_to_wimc(fetch, ERR, "tmp path too long");
        goto done;
    }

    if (pkg_extract_path(uuidStr, content->name, content->tag,
                         extractPath, sizeof(extractPath)) != 0) {
        send_update_to_wimc(fetch, ERR, "extract path failure");
        goto done;
    }

    rm_rf(uuidTmpDir);
    if (mkdir_p(uuidTmpDir, 0700) != 0) {
        send_update_to_wimc(fetch, ERR, "failed to create tmp dir");
        goto done;
    }

    send_update_to_wimc(fetch, FETCH, "");

    ok = 0;
    for (attempt = 1; attempt <= CASYNC_RETRIES; attempt++) {
        rm_rf(extractPath);
        usys_log_debug("casync attempt %d for %s:%s",
                       attempt, content->name, content->tag);

        if (run_casync(fetch, extractPath) == 0) {
            ok = 1;
            break;
        }

        sleep(attempt * 2);
    }

    if (!ok) {
        send_update_to_wimc(fetch, ERR, "casync failed");
        goto done;
    }

    memset(publishedPath, 0, sizeof(publishedPath));
    memset(actualVersion, 0, sizeof(actualVersion));

    if (pkg_publish_from_dir(content->name, content->tag, uuidStr,
                             extractPath, publishedPath,
                             sizeof(publishedPath), actualVersion,
                             sizeof(actualVersion)) != 0) {
        send_update_to_wimc(fetch, ERR, "package publish failed");
        goto done;
    }

    usys_log_debug("Published package %s:%s actual=%s path=%s",
                   content->name, content->tag, actualVersion,
                   publishedPath);
    send_update_to_wimc(fetch, DONE, publishedPath);

done:
    rm_rf(uuidTmpDir);
    free_fetch_request(fetch);
    pthread_exit(NULL);
}

void process_capp_fetch_request(WFetch *fetch) {

    int ret;
    pthread_t tid;
    WFetch *threadFetch = NULL;

    copy_fetch_request(&threadFetch, fetch);
    if (threadFetch == NULL) {
        usys_log_error("Unable to copy fetch request");
        return;
    }

    ret = pthread_create(&tid, NULL, execute_agent, threadFetch);
    if (ret) {
        usys_log_error("Error creating agent thread. Return code: %d", ret);
        free_fetch_request(threadFetch);
        return;
    }

    pthread_detach(tid);
    usys_log_debug("Agent fetch thread started");
}
