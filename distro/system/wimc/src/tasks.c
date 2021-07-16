/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string.h>

#include "wimc.h"
#include "agent.h"
#include "jserdes.h"

/* Functions related to tasks. */

void clear_tasks(WTasks **tasks) {

  WTasks *curr=NULL, *tmp;

  curr = *tasks;

  while (curr) {
    
    WContent *content = curr->content;
    
    if (content) {
      free(content->name);
      free(content->tag);
      free(content->method);
      free(content->providerURL);
    }

    free(curr->content);
    free(curr->update);
    if (curr->localPath)
      free(curr->localPath);

    tmp = curr->next;
    free(curr);
    curr = tmp; /*next entry */
  }
  
}

/*
 * copy_contents --
 *
 */

static void copy_contents(WContent *src, WContent *dest) {

  /* sanity check. */
  if (!src && !dest) {
    return;
  }

  dest->name        = strndup(src->name, strlen(src->name));
  dest->tag         = strndup(src->tag, strlen(src->tag));
  dest->method      = strndup(src->method, strlen(src->method));
  dest->providerURL = strndup(src->providerURL, strlen(src->providerURL));
}

/*
 * add_task_entry --
 *
 */

static void add_task_entry(WTasks **task, WimcReq *req) {

  WTasks *ptr;
  Update *update;

  /* sanity check. */
  if (!req && !req->fetch && !task)
    return;

  ptr = *task;

  ptr->content = (WContent *)calloc(1, sizeof(WContent));
  ptr->update  = (Update *)calloc(1, sizeof(Update));
  if (!ptr->content && !ptr->update) {
    return;
  }

  uuid_copy(ptr->uuid, req->fetch->uuid);
  copy_contents(req->fetch->content, ptr->content);
  ptr->state = (TransferState)REQUEST;

  update                = ptr->update;
  update->totalKB       = 0;
  update->transferKB    = 0;
  update->transferState = (int)ptr->state;
  update->voidStr       = NULL;
  uuid_copy(update->uuid, ptr->uuid);

  ptr->localPath = NULL;
}

/*
 * add_to_tasks --
 *
 */

void add_to_tasks(WTasks **tasks, WimcReq *req) {

  WTasks *curr;

  /* Base case - first entry. */
  if (*tasks == NULL) {
    *tasks = (WTasks *)calloc(1, sizeof(WTasks));
    if (!*tasks) {
      return;
    }

    add_task_entry(tasks, req);
    (*tasks)->next = NULL;
    return;
  }

  curr = *tasks;
  
  while (curr->next != NULL) {
    curr = curr->next;
  }

  curr->next = (WTasks *)calloc(1, sizeof(WTasks));
  if (!curr->next) {
    return;
  }

  add_task_entry(&(curr->next), req);
  curr->next->next = NULL;
}

/*
 * free_task_element --
 */

static void free_task_element(WTasks *task) {

  WContent *content = task->content;

  if (content) {
    free(content->name);
    free(content->tag);
    free(content->method);
    free(content->providerURL);
  }

  free(task->content);
  free(task->update);
  if (task->localPath)
    free(task->localPath);

  free(task);
  task = NULL;

  return;
}

/*
 * delete_from_tasks --
 *
 */

void delete_from_tasks(WTasks **tasks, WTasks *target) {

  WTasks *curr;

  /* Sanity check */
  if (target == NULL && *tasks == NULL)
    return;

  curr =*tasks;

  /* Cases: 1. first 2. middle and 3. last */
  if (curr == target) { /* First. */
    free_task_element(target);
    *tasks=NULL;
    return;
  }

  /* Middle and last element in the list. */
  while (curr->next !=NULL && curr->next != target) {
    curr = curr->next;
  }

  if (curr->next == target) {
    curr->next = target->next;
    free_task_element(target);
  }

  return;
}

/*
 * process_task_request --
 *
 */

char *process_task_request(WTasks *task) {

  int ret;
  json_t *json=NULL;
  char *retStr=NULL;

  ret = serialize_task(task, &json);
  
  if (ret) {

    retStr = json_dumps(json, 0);
    if (retStr) {
      log_debug("Task json-str: %s", retStr);
    }
  }

  json_decref(json);
  return retStr;
}

/*
 * process_cli response -- Handle response back to the client.
 *                     There are three type of responses:
 *                     1. Type: Result, Str: local path, Task: void
 *                     2. Type: Processing, Str: ID (for CB), Task: void
 *                     3. Type: Update, Str: void, Task: update_status
 *
 */

char *process_cli_response(WRespType type, char *path, char *idStr,
			   WTasks *task, char *errStr) {

  int ret=FALSE;
  char *retStr=NULL;
  json_t *json=NULL;

  /* Sanity check. */
  if (type == WRESP_RESULT && path == NULL) {
    return NULL;
  }

  if (type == WRESP_UPDATE && task == NULL) {
    return NULL;
  }

  if (type == WRESP_PROCESSING && idStr == NULL){
    return NULL;
  }

  if (type == WRESP_UPDATE) {
    ret = serialize_task(task, &json);
    if (ret) {
      retStr = json_dumps(json, 0);
    }
  } else if (type == (WRespType)WRESP_PROCESSING) {
    ret = serialize_result(type, idStr, &json);
  } else if (type == (WRespType)WRESP_RESULT) {
    ret = serialize_result(type, path, &json);
  } else if (type == (WRespType)WRESP_ERROR) {
    ret = serialize_result(type, errStr, &json);
  }

  if (ret && json) {
    retStr = json_dumps(json, 0);
    log_debug("Sending JSON response: %s", retStr);
    json_decref(json);
  }

  return retStr;
}
