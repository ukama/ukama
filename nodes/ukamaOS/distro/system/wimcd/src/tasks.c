/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>

#include "agent.h"
#include "jserdes.h"
#include "wimc.h"

#include "usys_log.h"
#include "usys_mem.h"

static void free_task_element(WTasks *task) {

    WContent *content;

    if (task == NULL) {
        return;
    }

    content = task->content;
    if (content != NULL) {
        free(content->name);
        free(content->tag);
        free(content->method);
        free(content->indexURL);
        free(content->storeURL);
    }

    free(task->content);

    if (task->update != NULL) {
        free(task->update->voidStr);
        free(task->update);
    }

    free(task->localPath);
    free(task);
}

void clear_tasks(WTasks **tasks) {

    WTasks *curr;
    WTasks *tmp;

    if (tasks == NULL) {
        return;
    }

    curr = *tasks;
    while (curr != NULL) {
        tmp = curr->next;
        free_task_element(curr);
        curr = tmp;
    }

    *tasks = NULL;
}

static void copy_contents(WContent *src, WContent *dest) {

    if (src == NULL || dest == NULL) {
        return;
    }

    dest->name = src->name ? strdup(src->name) : NULL;
    dest->tag = src->tag ? strdup(src->tag) : NULL;
    dest->method = src->method ? strdup(src->method) : NULL;
    dest->indexURL = src->indexURL ? strdup(src->indexURL) : NULL;
    dest->storeURL = src->storeURL ? strdup(src->storeURL) : NULL;
    dest->expectedSizeBytes = src->expectedSizeBytes;
}

static void add_task_entry(WTasks *task, WimcReq *req) {

    Update *update;

    if (task == NULL || req == NULL || req->fetch == NULL) {
        return;
    }

    task->content = (WContent *)calloc(1, sizeof(WContent));
    task->update = (Update *)calloc(1, sizeof(Update));
    if (task->content == NULL || task->update == NULL) {
        return;
    }

    uuid_copy(task->uuid, req->fetch->uuid);
    copy_contents(req->fetch->content, task->content);
    task->state = REQUEST;

    update = task->update;
    update->totalKB = 0;
    update->transferKB = 0;
    update->transferState = task->state;
    update->voidStr = NULL;
    uuid_copy(update->uuid, task->uuid);

    task->localPath = NULL;
    task->next = NULL;
}

void add_to_tasks(WTasks **tasks, WimcReq *req) {

    WTasks *curr;
    WTasks *newTask;

    if (tasks == NULL || req == NULL || req->fetch == NULL) {
        return;
    }

    newTask = (WTasks *)calloc(1, sizeof(WTasks));
    if (newTask == NULL) {
        return;
    }

    add_task_entry(newTask, req);
    if (newTask->content == NULL || newTask->update == NULL) {
        free_task_element(newTask);
        return;
    }

    if (*tasks == NULL) {
        *tasks = newTask;
        return;
    }

    curr = *tasks;
    while (curr->next != NULL) {
        curr = curr->next;
    }
    curr->next = newTask;
}

void delete_from_tasks(WTasks **tasks, WTasks *target) {

    WTasks *curr;

    if (tasks == NULL || *tasks == NULL || target == NULL) {
        return;
    }

    curr = *tasks;
    if (curr == target) {
        *tasks = target->next;
        free_task_element(target);
        return;
    }

    while (curr->next != NULL && curr->next != target) {
        curr = curr->next;
    }

    if (curr->next == target) {
        curr->next = target->next;
        free_task_element(target);
    }
}

char *process_task_request(WTasks *task) {

    int ret;
    json_t *json;
    char *retStr;

    json = NULL;
    retStr = NULL;

    ret = serialize_task(task, &json);
    if (ret) {
        retStr = json_dumps(json, 0);
        if (retStr != NULL) {
            usys_log_debug("Task json-str: %s", retStr);
        }
    }

    if (json != NULL) {
        json_decref(json);
    }

    return retStr;
}

char *process_cli_response(WRespType type, char *path, char *idStr,
                           WTasks *task, char *errStr) {

    int ret;
    char *retStr;
    json_t *json;

    (void)path;

    ret = FALSE;
    retStr = NULL;
    json = NULL;

    if (type == WRESP_RESULT &&
        (task == NULL || task->localPath == NULL)) {
        return NULL;
    }

    if (type == WRESP_UPDATE && task == NULL) {
        return NULL;
    }

    if (type == WRESP_PROCESSING && idStr == NULL) {
        return NULL;
    }

    if (type == WRESP_UPDATE) {
        ret = serialize_task(task, &json);
    } else if (type == WRESP_PROCESSING) {
        ret = serialize_result(type, idStr, &json);
    } else if (type == WRESP_RESULT) {
        ret = serialize_result(type, task->localPath, &json);
    } else if (type == WRESP_ERROR) {
        ret = serialize_result(type, errStr, &json);
    }

    if (ret && json != NULL) {
        retStr = json_dumps(json, 0);
        if (retStr != NULL) {
            usys_log_debug("Sending JSON response: %s", retStr);
        }
    }

    if (json != NULL) {
        json_decref(json);
    }

    return retStr;
}

WTasks *find_task_by_uuid(WTasks *tasks, uuid_t uuid) {

    WTasks *curr;

    curr = tasks;
    while (curr != NULL) {
        if (uuid_compare(curr->uuid, uuid) == 0) {
            return curr;
        }
        curr = curr->next;
    }

    return NULL;
}
