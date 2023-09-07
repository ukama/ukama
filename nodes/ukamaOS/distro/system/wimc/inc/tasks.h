/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WIMC_TASKS_H
#define WIMC_TASKS_H

void clear_tasks(WTasks **tasks);
void add_to_tasks(WTasks **tasks, WimcReq *req);
void delete_from_tasks(WTasks **tasks, WTasks *target);
char *process_task_request(WTasks *task);
char *process_cli_response(WRespType type, char *path, char *idStr,
                           WTasks *task, char *errStr); 

#endif /* WIMC_TASKS_H */
